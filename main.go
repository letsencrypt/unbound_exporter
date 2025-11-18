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
	"os"

	"github.com/prometheus/common/promslog"

	"github.com/letsencrypt/unbound_exporter/exporter"
	"github.com/letsencrypt/unbound_exporter/metrics"
)

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
		log.Error("Unbound Exporter setup failed", "err", err.Error())
		os.Exit(1)
	}

	log.Info("Starting server", "address", *listenAddress)
	err = metrics.NewMetricServer(*listenAddress, *metricsPath, *healthPath, exp)
	if err != nil {
		log.Error("Listen failed", "err", err.Error())
		os.Exit(1)
	}
}
