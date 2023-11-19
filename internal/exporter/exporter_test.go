package exporter

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
)

// metrics is an example of what Unbound returns from a "UBCT1 stats_noreset" command
const metrics string = `thread0.num.queries=369
thread0.num.queries_ip_ratelimited=0
thread0.num.queries_cookie_valid=0
thread0.num.queries_cookie_client=0
thread0.num.queries_cookie_invalid=0
thread0.num.cachehits=333
thread0.num.cachemiss=36
thread0.num.prefetch=0
thread0.num.queries_timed_out=0
thread0.query.queue_time_us.max=0
thread0.num.expired=0
thread0.num.recursivereplies=36
thread0.requestlist.avg=0
thread0.requestlist.max=0
thread0.requestlist.overwritten=0
thread0.requestlist.exceeded=0
thread0.requestlist.current.all=0
thread0.requestlist.current.user=0
thread0.recursion.time.avg=0.028036
thread0.recursion.time.median=0.0232825
thread0.tcpusage=0
thread1.num.queries=365
thread1.num.queries_ip_ratelimited=0
thread1.num.queries_cookie_valid=0
thread1.num.queries_cookie_client=0
thread1.num.queries_cookie_invalid=0
thread1.num.cachehits=340
thread1.num.cachemiss=25
thread1.num.prefetch=0
thread1.num.queries_timed_out=0
thread1.query.queue_time_us.max=0
thread1.num.expired=0
thread1.num.recursivereplies=25
thread1.requestlist.avg=0
thread1.requestlist.max=0
thread1.requestlist.overwritten=0
thread1.requestlist.exceeded=0
thread1.requestlist.current.all=0
thread1.requestlist.current.user=0
thread1.recursion.time.avg=0.043104
thread1.recursion.time.median=0.0251611
thread1.tcpusage=0
thread2.num.queries=373
thread2.num.queries_ip_ratelimited=0
thread2.num.queries_cookie_valid=0
thread2.num.queries_cookie_client=0
thread2.num.queries_cookie_invalid=0
thread2.num.cachehits=330
thread2.num.cachemiss=43
thread2.num.prefetch=0
`

func TestCollectFromReader(t *testing.T) {
	// Channel that CollectFrom Re
	metricsCh := make(chan prometheus.Metric)
	doneCh := make(chan error)

	go func(ch chan prometheus.Metric) {
		doneCh <- CollectFromReader(bytes.NewReader([]byte(metrics)), ch)
	}(metricsCh)

	var metrics []prometheus.Metric

L:
	for {
		select {
		case metric, ok := <-metricsCh:
			if !ok {
				t.Fatal("collector channel unexpectedly closed")
			}
			metrics = append(metrics, metric)
		case err := <-doneCh:
			if err != nil {
				t.Fatalf("Failed to CollectFromReader: %v", err)
			}
			break L
		}
	}

	//

	for _, metric := range metrics {
		fmt.Printf("%v\n", metric.Desc())
	}
}
