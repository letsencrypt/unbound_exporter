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
	"net"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors/version"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/promslog"

	"github.com/letsencrypt/unbound_exporter/exporter"
)

func main() {
	log := promslog.New(&promslog.Config{})

	var (
		listenAddress = flag.String("web.listen-address", ":9167", "Address to listen on for web interface and telemetry.")
		metricsPath   = flag.String("web.telemetry-path", "/metrics", "Path under which to expose metrics.")
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

	http.Handle(*metricsPath, promhttp.Handler())
	http.HandleFunc("/_healthz", func(w http.ResponseWriter, req *http.Request) {
		if exp.UnboundUp() {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("ok"))
		} else {
			w.WriteHeader(http.StatusServiceUnavailable)
			_, _ = w.Write([]byte("sad"))
		}
	})
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
	{
		address, port, err := net.SplitHostPort(*listenAddress)
		if err != nil {
			log.Error("Cannot parse web.listen-address", "err", err.Error())
			os.Exit(1)
		}
		log.Info("Listening", "address", address, "port", port)
	}
	err = http.ListenAndServe(*listenAddress, nil)
	log.Error("Listen failed", "err", err.Error())
	os.Exit(1)
}
