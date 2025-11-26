package exporter

import (
	"os"
	"regexp"
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

	err = collectFromReader(compileMetrics(), testData, ch)
	if err != nil {
		t.Fatal(err)
	}

	close(ch)
	<-done

	if len(metrics) != 109 {
		t.Fatal("expected 109 metrics, got ", len(metrics))
	}
}

func TestLabels(t *testing.T) {
	for _, metric := range unboundMetrics {
		r := regexp.MustCompile(metric.pattern)
		if r.NumSubexp() != len(metric.labels) {
			t.Errorf("Expected %d patterns in regex, got %d on %s", len(metric.labels), r.NumSubexp(), metric.name)
		}
	}
}
