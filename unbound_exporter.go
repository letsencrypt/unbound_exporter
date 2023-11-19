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
	"flag"
	"net/http"
	"os"

	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/promlog"

	"github.com/letsencrypt/unbound_exporter/internal/exporter"
)

var (
	log = promlog.New(&promlog.Config{})
)

func main() {
	var (
		listenAddress = flag.String("web.listen-address", ":9167", "Address to listen on for web interface and telemetry.")
		metricsPath   = flag.String("web.telemetry-path", "/metrics", "Path under which to expose metrics.")
		unboundHost   = flag.String("unbound.host", "tcp://localhost:8953", "Unix or TCP address of Unbound control socket.")
		unboundCa     = flag.String("unbound.ca", "/etc/unbound/unbound_server.pem", "Unbound server certificate.")
		unboundCert   = flag.String("unbound.cert", "/etc/unbound/unbound_control.pem", "Unbound client certificate.")
		unboundKey    = flag.String("unbound.key", "/etc/unbound/unbound_control.key", "Unbound client key.")
	)
	flag.Parse()

	_ = level.Info(log).Log("Starting unbound_exporter")
	exporter, err := exporter.NewUnboundExporter(log, *unboundHost, *unboundCa, *unboundCert, *unboundKey)
	if err != nil {
		panic(err)
	}
	prometheus.MustRegister(exporter)

	http.Handle(*metricsPath, promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`
			<html>
			<head><title>Unbound Exporter</title></head>
			<body>
			<h1>Unbound Exporter</h1>
			<p><a href='` + *metricsPath + `'>Metrics</a></p>
			</body>
			</html>`))
	})
	_ = level.Info(log).Log("Listening on address:port => ", *listenAddress)
	_ = level.Error(log).Log(http.ListenAndServe(*listenAddress, nil))
	os.Exit(1)
}
