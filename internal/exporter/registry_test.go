// SPDX-FileCopyrightText: 2025 UnionTech Software Technology Co., Ltd.
// SPDX-License-Identifier: MIT

package exporter

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus"
)

type testMetric struct {
	desc *prometheus.Desc
}

func (m testMetric) Describe(ch chan<- *prometheus.Desc) {
	ch <- m.desc
}

func (m testMetric) Collect(ch chan<- prometheus.Metric) {
	ch <- prometheus.MustNewConstMetric(m.desc, prometheus.GaugeValue, 1)
}

func TestRegistryDescribeForwardsCollectorDescriptions(t *testing.T) {
	registry := NewRegistry()
	metric := testMetric{
		desc: prometheus.NewDesc("test_metric_total", "test metric", nil, nil),
	}
	registry.Register(metric)

	descCh := make(chan *prometheus.Desc, 1)
	registry.Describe(descCh)
	close(descCh)

	count := 0
	for desc := range descCh {
		if desc == nil {
			t.Fatal("received nil desc")
		}
		count++
	}

	if count != 1 {
		t.Fatalf("unexpected number of descs: got %d, want 1", count)
	}
}

func TestRegistryGetMetricsReturnsCopy(t *testing.T) {
	registry := NewRegistry()
	original := testMetric{
		desc: prometheus.NewDesc("test_metric_total", "test metric", nil, nil),
	}
	replacement := testMetric{
		desc: prometheus.NewDesc("replacement_metric_total", "replacement metric", nil, nil),
	}
	registry.Register(original)

	metrics := registry.GetMetrics()
	metrics[0] = replacement

	got := registry.GetMetrics()
	if got[0] != original {
		t.Fatalf("expected registry metrics to remain unchanged")
	}
}