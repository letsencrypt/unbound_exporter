//go:build integration

package main

import (
	"net/http"
	"testing"

	"github.com/prometheus/common/expfmt"
	"github.com/prometheus/common/model"
)

// TestIntegration checks that unbound_exporter is running, successfully
// scraping and exporting metrics.
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
	resp, err := http.Get("http://localhost:9167/metrics")
	if err != nil {
		t.Fatalf("Failed to fetch metrics from unbound_exporter: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected a 200 OK from unbound_exporter, got: %v", resp.StatusCode)
	}

	parser := expfmt.NewTextParser(model.UTF8Validation)
	metrics, err := parser.TextToMetricFamilies(resp.Body)
	if err != nil {
		t.Fatalf("Failed to parse metrics from unbound_exporter: %v", err)
	}

	// unbound_up is 1 if we've successfully scraped metrics from it
	unbound_up := metrics["unbound_up"].Metric[0].Gauge.GetValue()
	if unbound_up != 1 {
		t.Errorf("Expected unbound_up to be 1, not: %v", unbound_up)
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
			t.Errorf("Expected metric is missing: %s", metric)
		}
	}

	resp, err = http.Get("http://localhost:9167/_healthz")
	if err != nil {
		t.Fatalf("Failed to fetch healthz from unbound_exporter: %v", err)
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("unbound_exporter reported unhealthy, status code: %d", resp.StatusCode)
	}

}
