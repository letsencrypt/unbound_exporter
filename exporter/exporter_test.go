package exporter

import (
	"os"
	"regexp"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
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

func TestCollectInfo(t *testing.T) {
	testData, err := os.Open("testdata/status.txt")
	if err != nil {
		t.Fatal(err)
	}

	ch := make(chan prometheus.Metric, 1)
	if err := collectInfoFromReader(testData, ch); err != nil {
		t.Fatal(err)
	}
	close(ch)

	metric, ok := <-ch
	if !ok {
		t.Fatal("expected one metric, got none")
	}

	var m dto.Metric
	if err := metric.Write(&m); err != nil {
		t.Fatal(err)
	}

	want := map[string]string{
		"version": "1.25.0",
		"threads": "16",
		"modules": "validator iterator",
	}
	for _, label := range m.GetLabel() {
		if want[label.GetName()] != label.GetValue() {
			t.Errorf("label %s: expected %q, got %q", label.GetName(), want[label.GetName()], label.GetValue())
		}
		delete(want, label.GetName())
	}
	for name := range want {
		t.Errorf("missing label %s", name)
	}
	if m.GetGauge().GetValue() != 1 {
		t.Errorf("expected value 1, got %f", m.GetGauge().GetValue())
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
