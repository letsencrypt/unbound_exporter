package exporter

import (
	"bufio"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/url"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync/atomic"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	unboundUpDesc = prometheus.NewDesc(
		prometheus.BuildFQName("unbound", "", "up"),
		"Whether scraping Unbound's metrics was successful.",
		nil, nil)

	unboundHistogram = prometheus.NewDesc(
		prometheus.BuildFQName("unbound", "", "response_time_seconds"),
		"Query response time in seconds.",
		nil, nil)

	unboundMetrics = []*unboundMetric{
		newUnboundMetric(
			"answer_rcodes_total",
			"Total number of answers to queries, from cache or from recursion, by response code.",
			prometheus.CounterValue,
			[]string{"rcode"},
			"^num\\.answer\\.rcode\\.(\\w+)$"),
		newUnboundMetric(
			"answers_bogus",
			"Total number of answers that were bogus.",
			prometheus.CounterValue,
			nil,
			"^num\\.answer\\.bogus$"),
		newUnboundMetric(
			"answers_secure_total",
			"Total number of answers that were secure.",
			prometheus.CounterValue,
			nil,
			"^num\\.answer\\.secure$"),
		newUnboundMetric(
			"cache_hits_total",
			"Total number of queries that were successfully answered using a cache lookup.",
			prometheus.CounterValue,
			[]string{"thread"},
			"^thread(\\d+)\\.num\\.cachehits$"),
		newUnboundMetric(
			"cache_misses_total",
			"Total number of cache queries that needed recursive processing.",
			prometheus.CounterValue,
			[]string{"thread"},
			"^thread(\\d+)\\.num\\.cachemiss$"),
		newUnboundMetric(
			"query_subnet_total",
			"Total number of queries that got an answer that contained EDNS client subnet data.",
			prometheus.CounterValue,
			nil,
			"^num\\.query\\.subnet$"),
		newUnboundMetric(
			"query_subnet_cache_total",
			"Total number of queries answered from the edns client subnet cache.",
			prometheus.CounterValue,
			nil,
			"^num\\.query\\.subnet_cache$"),
		newUnboundMetric(
			"queries_cookie_client_total",
			"Total number of queries with a client cookie.",
			prometheus.CounterValue,
			[]string{"thread"},
			"^thread(\\d+)\\.num\\.queries_cookie_client$"),
		newUnboundMetric(
			"queries_cookie_invalid_total",
			"Total number of queries with a invalid cookie.",
			prometheus.CounterValue,
			[]string{"thread"},
			"^thread(\\d+)\\.num\\.queries_invalid_client$"),
		newUnboundMetric(
			"queries_cookie_valid_total",
			"Total number of queries with a valid cookie.",
			prometheus.CounterValue,
			[]string{"thread"},
			"^thread(\\d+)\\.num\\.queries_cookie_valid$"),
		newUnboundMetric(
			"memory_caches_bytes",
			"Memory in bytes in use by caches.",
			prometheus.GaugeValue,
			[]string{"cache"},
			"^mem\\.cache\\.(\\w+)$"),
		newUnboundMetric(
			"memory_modules_bytes",
			"Memory in bytes in use by modules.",
			prometheus.GaugeValue,
			[]string{"module"},
			"^mem\\.mod\\.(\\w+)$"),
		newUnboundMetric(
			"memory_sbrk_bytes",
			"Memory in bytes allocated through sbrk.",
			prometheus.GaugeValue,
			nil,
			"^mem\\.total\\.sbrk$"),
		newUnboundMetric(
			"prefetches_total",
			"Total number of cache prefetches performed.",
			prometheus.CounterValue,
			[]string{"thread"},
			"^thread(\\d+)\\.num\\.prefetch$"),
		newUnboundMetric(
			"queries_total",
			"Total number of queries received.",
			prometheus.CounterValue,
			[]string{"thread"},
			"^thread(\\d+)\\.num\\.queries$"),
		newUnboundMetric(
			"expired_total",
			"Total number of expired entries served.",
			prometheus.CounterValue,
			[]string{"thread"},
			"^thread(\\d+)\\.num\\.expired$"),
		newUnboundMetric(
			"query_classes_total",
			"Total number of queries with a given query class.",
			prometheus.CounterValue,
			[]string{"class"},
			"^num\\.query\\.class\\.([\\w]+)$"),
		newUnboundMetric(
			"query_flags_total",
			"Total number of queries that had a given flag set in the header.",
			prometheus.CounterValue,
			[]string{"flag"},
			"^num\\.query\\.flags\\.([\\w]+)$"),
		newUnboundMetric(
			"query_ipv6_total",
			"Total number of queries that were made using IPv6 towards the Unbound server.",
			prometheus.CounterValue,
			nil,
			"^num\\.query\\.ipv6$"),
		newUnboundMetric(
			"query_opcodes_total",
			"Total number of queries with a given query opcode.",
			prometheus.CounterValue,
			[]string{"opcode"},
			"^num\\.query\\.opcode\\.([\\w]+)$"),
		newUnboundMetric(
			"query_edns_DO_total",
			"Total number of queries that had an EDNS OPT record with the DO (DNSSEC OK) bit set present.",
			prometheus.CounterValue,
			nil,
			"^num\\.query\\.edns\\.DO$"),
		newUnboundMetric(
			"query_edns_present_total",
			"Total number of queries that had an EDNS OPT record present.",
			prometheus.CounterValue,
			nil,
			"^num\\.query\\.edns\\.present$"),
		newUnboundMetric(
			"query_tcp_total",
			"Total number of queries that were made using TCP towards the Unbound server, including DoT and DoH queries.",
			prometheus.CounterValue,
			nil,
			"^num\\.query\\.tcp$"),
		newUnboundMetric(
			"query_tcpout_total",
			"Total number of queries that the Unbound server made using TCP outgoing towards other servers.",
			prometheus.CounterValue,
			nil,
			"^num\\.query\\.tcpout$"),
		newUnboundMetric(
			"query_tls_total",
			"Total number of queries that were made using TCP TLS towards the Unbound server, including DoT and DoH queries.",
			prometheus.CounterValue,
			nil,
			"^num\\.query\\.tls$"),
		newUnboundMetric(
			"query_tls_resume_total",
			"Total number of queries that were made using TCP TLS Resume towards the Unbound server.",
			prometheus.CounterValue,
			nil,
			"^num\\.query\\.tls\\.resume$"),
		newUnboundMetric(
			"query_https_total",
			"Total number of DoH queries that were made towards the Unbound server.",
			prometheus.CounterValue,
			nil,
			"^num\\.query\\.https$"),
		newUnboundMetric(
			"query_types_total",
			"Total number of queries with a given query type.",
			prometheus.CounterValue,
			[]string{"type"},
			"^num\\.query\\.type\\.([\\w]+)$"),
		newUnboundMetric(
			"query_udpout_total",
			"Total number of queries that the Unbound server made using UDP outgoing towards￼other servers.",
			prometheus.CounterValue,
			nil,
			"^num\\.query\\.udpout$"),
		newUnboundMetric(
			"query_aggressive_nsec",
			"Total number of queries that the Unbound server generated response using Aggressive NSEC.",
			prometheus.CounterValue,
			[]string{"rcode"},
			"^num\\.query\\.aggressive\\.(\\w+)$"),
		newUnboundMetric(
			"request_list_current_all",
			"Current size of the request list, including internally generated queries.",
			prometheus.GaugeValue,
			[]string{"thread"},
			"^thread([0-9]+)\\.requestlist\\.current\\.all$"),
		newUnboundMetric(
			"request_list_current_user",
			"Current size of the request list, only counting the requests from client queries.",
			prometheus.GaugeValue,
			[]string{"thread"},
			"^thread([0-9]+)\\.requestlist\\.current\\.user$"),
		newUnboundMetric(
			"request_list_exceeded_total",
			"Number of queries that were dropped because the request list was full.",
			prometheus.CounterValue,
			[]string{"thread"},
			"^thread([0-9]+)\\.requestlist\\.exceeded$"),
		newUnboundMetric(
			"request_list_overwritten_total",
			"Total number of requests in the request list that were overwritten by newer entries.",
			prometheus.CounterValue,
			[]string{"thread"},
			"^thread([0-9]+)\\.requestlist\\.overwritten$"),
		newUnboundMetric(
			"recursive_replies_total",
			"Total number of replies sent to queries that needed recursive processing.",
			prometheus.CounterValue,
			[]string{"thread"},
			"^thread(\\d+)\\.num\\.recursivereplies$"),
		newUnboundMetric(
			"rrset_bogus_total",
			"Total number of rrsets marked bogus by the validator.",
			prometheus.CounterValue,
			nil,
			"^num\\.rrset\\.bogus$"),
		newUnboundMetric(
			"rrset_cache_max_collisions_total",
			"Total number of rrset cache hashtable collisions.",
			prometheus.CounterValue,
			nil,
			"^rrset\\.cache\\.max_collisions$"),
		newUnboundMetric(
			"time_elapsed_seconds",
			"Time since last statistics printout in seconds.",
			prometheus.CounterValue,
			nil,
			"^time\\.elapsed$"),
		newUnboundMetric(
			"time_now_seconds",
			"Current time in seconds since 1970.",
			prometheus.GaugeValue,
			nil,
			"^time\\.now$"),
		newUnboundMetric(
			"time_up_seconds_total",
			"Uptime since server boot in seconds.",
			prometheus.CounterValue,
			nil,
			"^time\\.up$"),
		newUnboundMetric(
			"unwanted_queries_total",
			"Total number of queries that were refused or dropped because they failed the access control settings.",
			prometheus.CounterValue,
			nil,
			"^unwanted\\.queries$"),
		newUnboundMetric(
			"unwanted_replies_total",
			"Total number of replies that were unwanted or unsolicited.",
			prometheus.CounterValue,
			nil,
			"^unwanted\\.replies$"),
		newUnboundMetric(
			"recursion_time_seconds_avg",
			"Average time it took to answer queries that needed recursive processing (does not include in-cache requests).",
			prometheus.GaugeValue,
			nil,
			"^total\\.recursion\\.time\\.avg$"),
		newUnboundMetric(
			"recursion_time_seconds_median",
			"The median of the time it took to answer queries that needed recursive processing.",
			prometheus.GaugeValue,
			nil,
			"^total\\.recursion\\.time\\.median$"),
		newUnboundMetric(
			"msg_cache_count",
			"The Number of Messages cached",
			prometheus.GaugeValue,
			nil,
			"^msg\\.cache\\.count$"),
		newUnboundMetric(
			"msg_cache_max_collisions_total",
			"Total number of msg cache hashtable collisions.",
			prometheus.CounterValue,
			nil,
			"^msg\\.cache\\.max_collisions$"),
		newUnboundMetric(
			"rrset_cache_count",
			"The Number of rrset cached",
			prometheus.GaugeValue,
			nil,
			"^rrset\\.cache\\.count$"),
		newUnboundMetric(
			"rpz_action_count",
			"Total number of triggered Response Policy Zone actions, by type.",
			prometheus.CounterValue,
			[]string{"type"},
			"^num\\.rpz\\.action\\.rpz-([\\w-]+)$"),
		newUnboundMetric(
			"memory_doh_bytes",
			"Memory used by DoH buffers, in bytes.",
			prometheus.GaugeValue,
			[]string{"buffer"},
			"^mem\\.http\\.(\\w+)$"),
		newUnboundMetric(
			"infra_cache_count",
			"Total number of infra cache entries",
			prometheus.CounterValue,
			nil,
			"^infra\\.cache\\.count$"),
		newUnboundMetric(
			"memory_doq_bytes",
			"Memory used by DoQ buffers, in bytes.",
			prometheus.GaugeValue,
			[]string{"buffer"},
			"^mem\\.quic$"),
		newUnboundMetric(
			"query_quic_total",
			"Total number of DNS-over-QUIC (DoQ) queries performed towards the Unbound server.",
			prometheus.CounterValue,
			nil,
			"^num\\.query\\.quic$"),
		newUnboundMetric(
			"dns_error_reports",
			"Total number of DNS Error Reports generated",
			prometheus.CounterValue,
			nil,
			"^num\\.dns_error_reports$"),
		newUnboundMetric(
			"queries_discard_timeout",
			"Total number of queries removed due to discard-timeout.",
			prometheus.CounterValue,
			nil,
			"^num\\.queries_discard_timeout$"),
		newUnboundMetric(
			"queries_wait_limit",
			"Total number of queries removed due to wait-limit.",
			prometheus.CounterValue,
			nil,
			"^num\\.queries_wait_limit$"),
		newUnboundMetric(
			"signature_validations",
			"Total umber of signature validation operations performed by the validator module",
			prometheus.CounterValue,
			nil,
			"^num\\.valops$"),
	}
)

type unboundMetric struct {
	desc      *prometheus.Desc
	valueType prometheus.ValueType
	pattern   *regexp.Regexp
}

func newUnboundMetric(name string, description string, valueType prometheus.ValueType, labels []string, pattern string) *unboundMetric {
	return &unboundMetric{
		desc: prometheus.NewDesc(
			prometheus.BuildFQName("unbound", "", name),
			description,
			labels,
			nil),
		valueType: valueType,
		pattern:   regexp.MustCompile(pattern),
	}
}

func CollectFromReader(file io.Reader, ch chan<- prometheus.Metric) error {
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	histogramPattern := regexp.MustCompile(`^histogram\.\d+\.\d+\.to\.(\d+\.\d+)$`)

	histogramCount := uint64(0)
	histogramAvg := float64(0)
	histogramBuckets := make(map[float64]uint64)

	for scanner.Scan() {
		fields := strings.Split(scanner.Text(), "=")
		if len(fields) != 2 {
			return fmt.Errorf(
				"%q is not a valid key-value pair",
				scanner.Text())
		}

		for _, metric := range unboundMetrics {
			if matches := metric.pattern.FindStringSubmatch(fields[0]); matches != nil {
				value, err := strconv.ParseFloat(fields[1], 64)
				if err != nil {
					return err
				}
				ch <- prometheus.MustNewConstMetric(
					metric.desc,
					metric.valueType,
					value,
					matches[1:]...)

				break
			}
		}

		if matches := histogramPattern.FindStringSubmatch(fields[0]); matches != nil {
			end, err := strconv.ParseFloat(matches[1], 64)
			if err != nil {
				return err
			}
			value, err := strconv.ParseUint(fields[1], 10, 64)
			if err != nil {
				return err
			}
			histogramBuckets[end] = value
			histogramCount += value
		} else if fields[0] == "total.recursion.time.avg" {
			value, err := strconv.ParseFloat(fields[1], 64)
			if err != nil {
				return err
			}
			histogramAvg = value
		}
	}

	// Convert the metrics to a cumulative Prometheus histogram.
	// Reconstruct the sum of all samples from the average value
	// provided by Unbound. Hopefully this does not break
	// monotonicity.
	keys := []float64{}
	for k := range histogramBuckets {
		keys = append(keys, k)
	}
	sort.Float64s(keys)
	prev := uint64(0)
	for _, i := range keys {
		histogramBuckets[i] += prev
		prev = histogramBuckets[i]
	}
	ch <- prometheus.MustNewConstHistogram(
		unboundHistogram,
		histogramCount,
		histogramAvg*float64(histogramCount),
		histogramBuckets)

	return scanner.Err()
}

func CollectFromSocket(socketFamily string, host string, tlsConfig *tls.Config, ch chan<- prometheus.Metric) error {
	var (
		conn net.Conn
		err  error
	)

	if socketFamily == "unix" || tlsConfig == nil {
		conn, err = net.Dial(socketFamily, host)
	} else {
		conn, err = tls.Dial(socketFamily, host, tlsConfig)
	}
	if err != nil {
		return err
	}
	defer conn.Close()
	_, err = conn.Write([]byte("UBCT1 stats_noreset\n"))
	if err != nil {
		return err
	}
	return CollectFromReader(conn, ch)
}

type UnboundExporter struct {
	log *slog.Logger

	socketFamily string
	host         string
	tlsConfig    *tls.Config

	// unboundUp is true if the last scrape was healthy. Used for /_healthz
	// False initially, so this will return unhealthy until the first metric scrape has succeeded.
	unboundUp atomic.Bool
}

func tlsConfig(ca string, cert string, key string) (*tls.Config, error) {
	// Server authentication
	caData, err := os.ReadFile(ca)
	if err != nil {
		return nil, err
	}
	roots := x509.NewCertPool()
	if !roots.AppendCertsFromPEM(caData) {
		return nil, errors.New("failed to parse CA")
	}

	// Client authentication
	certData, err := os.ReadFile(cert)
	if err != nil {
		return nil, err
	}

	keyData, err := os.ReadFile(key)
	if err != nil {
		return nil, err
	}

	keyPair, err := tls.X509KeyPair(certData, keyData)
	if err != nil {
		return nil, err
	}

	return &tls.Config{
		Certificates: []tls.Certificate{keyPair},
		RootCAs:      roots,
		ServerName:   "unbound",
	}, nil
}

func NewUnboundExporter(host string, ca string, cert string, key string, log *slog.Logger) (*UnboundExporter, error) {
	u, err := url.Parse(host)
	if err != nil {
		return &UnboundExporter{}, err
	}

	newExporter := UnboundExporter{
		log:          log,
		socketFamily: u.Scheme,
		host:         u.Path,
	}

	if u.Scheme == "unix" || (ca == "" && cert == "" && key == "") {
		return &newExporter, nil
	}

	cfg, err := tlsConfig(ca, cert, key)
	if err != nil {
		return nil, err
	}
	newExporter.tlsConfig = cfg

	return &newExporter, nil
}

func (e *UnboundExporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- unboundUpDesc
	for _, metric := range unboundMetrics {
		ch <- metric.desc
	}
}

func (e *UnboundExporter) Collect(ch chan<- prometheus.Metric) {
	err := CollectFromSocket(e.socketFamily, e.host, e.tlsConfig, ch)
	if err == nil {
		e.unboundUp.Store(true)
		ch <- prometheus.MustNewConstMetric(
			unboundUpDesc,
			prometheus.GaugeValue,
			1.0)
	} else {
		e.log.Error("Failed to scrape socket", "err", err.Error())
		e.unboundUp.Store(false)
		ch <- prometheus.MustNewConstMetric(
			unboundUpDesc,
			prometheus.GaugeValue,
			0.0)
	}
}

// UnboundUp returns true if we have successfully scraped metrics and Unbound is up.
func (e *UnboundExporter) UnboundUp() bool {
	return e.unboundUp.Load()
}
