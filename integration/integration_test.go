//go:build integration

package main

import (
	"net/http"
	"testing"

	"github.com/prometheus/common/expfmt"
	"github.com/prometheus/common/model"
)

// TestIntegration checks that unbound_exporter is running, successfully
// scraping and exporting metrics using Unix socket.
//
// It assumes unbound_exporter is available on localhost:9167, and Unbound on
// localhost:1053, as is set up in the docker-compose.yml file.
//
// A typical invocation of this test would look like
//
//	docker compose up --build -d
//	go test --tags=integration
//	docker compose down
func TestIntegration(t *testing.T) {
	testExporter(t, "http://localhost:9167", "Unix socket")
}

// TestIntegrationTCP checks that unbound_exporter is running, successfully
// scraping and exporting metrics using TCP (without TLS).
//
// It assumes unbound_exporter is available on localhost:9168, and Unbound TCP
// control interface on localhost:8953, as is set up in the docker-compose.yml file.
func TestIntegrationTCP(t *testing.T) {
	testExporter(t, "http://localhost:9168", "TCP")
}

// TestIntegrationTLS checks that unbound_exporter is running, successfully
// scraping and exporting metrics using TLS.
//
// It assumes unbound_exporter is available on localhost:9169, and Unbound TLS
// control interface on localhost:8954, as is set up in the docker-compose.yml file.
func TestIntegrationTLS(t *testing.T) {
	testExporter(t, "http://localhost:9169", "TLS")
}

// testExporter is a helper function that tests an unbound_exporter instance.
func testExporter(t *testing.T, exporterURL string, connectionType string) {
	t.Logf("Testing %s connection to unbound_exporter at %s", connectionType, exporterURL)

	resp, err := http.Get(exporterURL + "/metrics")
	if err != nil {
		t.Fatalf("Failed to fetch metrics from unbound_exporter (%s): %v", connectionType, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected a 200 OK from unbound_exporter (%s), got: %v", connectionType, resp.StatusCode)
	}

	parser := expfmt.NewTextParser(model.UTF8Validation)
	metrics, err := parser.TextToMetricFamilies(resp.Body)
	if err != nil {
		t.Fatalf("Failed to parse metrics from unbound_exporter (%s): %v", connectionType, err)
	}

	// unbound_up is 1 if we've successfully scraped metrics from it
	unbound_up := metrics["unbound_up"].Metric[0].Gauge.GetValue()
	if unbound_up != 1 {
		t.Errorf("Expected unbound_up to be 1 for %s connection, not: %v", connectionType, unbound_up)
	}

	// Check some expected metrics are present
	for _, metric := range []string{
		"go_info",
		"unbound_exporter_build_info",
		"unbound_queries_total",
		"unbound_response_time_seconds",
		"unbound_cache_hits_total",
		"unbound_query_https_total",
		"unbound_memory_doh_bytes",
		"unbound_query_subnet_total",
		"unbound_query_subnet_cache_total",
	} {
		if _, ok := metrics[metric]; !ok {
			t.Errorf("Expected metric is missing for %s connection: %s", connectionType, metric)
		}
	}

	resp, err = http.Get(exporterURL + "/_healthz")
	if err != nil {
		t.Fatalf("Failed to fetch healthz from unbound_exporter (%s): %v", connectionType, err)
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("unbound_exporter reported unhealthy for %s connection, status code: %d", connectionType, resp.StatusCode)
	}

	t.Logf("Successfully validated %s connection", connectionType)
}

