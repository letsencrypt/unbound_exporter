package exporter

import (
	"os"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
)

// TestCollect is a basic unit test for parsing the output format
func TestCollect(t *testing.T) {
	testData, err := os.Open("testdata/metrics.txt")
	if err != nil {
		t.Fatal(err)
	}

	ch := make(chan prometheus.Metric)
	done := make(chan struct{})

	var metrics []prometheus.Metric
	go func() {
		for m := range ch {
			metrics = append(metrics, m)
		}
		done <- struct{}{}
	}()

	err = CollectFromReader(testData, ch)
	if err != nil {
		t.Fatal(err)
	}

	close(ch)
	<-done

	if len(metrics) != 93 {
		t.Fatal("expected 93 metrics, got ", len(metrics))
	}
}
