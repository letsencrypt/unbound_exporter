// Copyright 2017 Kumina, https://kumina.nl/
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"bytes"
	"flag"
	"html/template"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors/version"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/promslog"

	"github.com/letsencrypt/unbound_exporter/exporter"
)

const homePageTemplate string = `
<!DOCTYPE html>
<html>
	<head><title>Unbound Exporter</title></head>
	<body>
		<h1>Unbound Exporter</h1>
		<p><a href='{{ .MetricsPath }}'>Metrics</a></p>
		<p><a href='{{ .HealthPath }}'>Health</a></p>
	</body>
</html>
`

// homePageText renders the html template for the homepage with the user-configured links
func homePageText(metricsPath, healthPath string) []byte {
	tmpl := template.Must(template.New("homePage").Parse(homePageTemplate))

	var out bytes.Buffer
	err := tmpl.Execute(&out, struct {
		MetricsPath string
		HealthPath  string
	}{
		MetricsPath: metricsPath,
		HealthPath:  healthPath,
	})
	if err != nil {
		panic(err) // Unreachable: Template is static and known to be well-formed
	}
	return out.Bytes()
}

// newMetricServer starts the http server on listenAddress
func newMetricsServer(listenAddress, metricsPath, healthPath string, exp *exporter.UnboundExporter) error {
	http.Handle(metricsPath, promhttp.Handler())

	http.HandleFunc(healthPath, func(w http.ResponseWriter, req *http.Request) {
		if exp.UnboundUp() {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("ok"))
		} else {
			w.WriteHeader(http.StatusServiceUnavailable)
			_, _ = w.Write([]byte("sad"))
		}
	})

	renderedHomePage := homePageText(metricsPath, healthPath)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write(renderedHomePage)
	})

	return http.ListenAndServe(listenAddress, nil)
}

func main() {
	log := promslog.New(&promslog.Config{})

	var (
		listenAddress = flag.String("web.listen-address", ":9167", "Address to listen on for web interface and telemetry.")
		metricsPath   = flag.String("web.telemetry-path", "/metrics", "Path under which to expose metrics.")
		healthPath    = flag.String("web.health-path", "/_healthz", "Path under which to expose healthcheck.")
		unboundHost   = flag.String("unbound.host", "tcp://localhost:8953", "Unix or TCP address of Unbound control socket.")
		unboundCa     = flag.String("unbound.ca", "/etc/unbound/unbound_server.pem", "Unbound server certificate.")
		unboundCert   = flag.String("unbound.cert", "/etc/unbound/unbound_control.pem", "Unbound client certificate.")
		unboundKey    = flag.String("unbound.key", "/etc/unbound/unbound_control.key", "Unbound client key.")
	)
	flag.Parse()

	log.Info("Starting unbound_exporter")
	exp, err := exporter.NewUnboundExporter(*unboundHost, *unboundCa, *unboundCert, *unboundKey, log)
	if err != nil {
		panic(err)
	}

	prometheus.MustRegister(exp)
	prometheus.MustRegister(version.NewCollector("unbound_exporter"))

	log.Info("Starting server", "address", *listenAddress)
	err = newMetricsServer(*listenAddress, *metricsPath, *healthPath, exp)
	if err != nil {
		log.Error("Listen failed", "err", err.Error())
		os.Exit(1)
	}
}
