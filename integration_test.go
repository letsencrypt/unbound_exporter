//go:build integration

package main

import "testing"

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
	// TODO
}
