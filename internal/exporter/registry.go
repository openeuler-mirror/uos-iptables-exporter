// SPDX-FileCopyrightText: 2025 UnionTech Software Technology Co., Ltd.
// SPDX-License-Identifier: MIT

package exporter

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

var defaultReg *Registry

func init() {
	defaultReg = NewRegistry()
}

type Registry struct {
	metrics []Metric
	mu      sync.RWMutex
}

func Register(metrics ...Metric) {
	for _, metric := range metrics {
		defaultReg.Register(metric)
	}
}

func RegisterPrometheus(reg *prometheus.Registry) {
	reg.MustRegister(defaultReg)
}

func NewRegistry() *Registry {
	return &Registry{
		metrics: []Metric{},
	}
}

func (r *Registry) Register(metrics Metric) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.metrics = append(r.metrics, metrics)
}

func (r *Registry) GetMetrics() []Metric {
	r.mu.RLock()
	defer r.mu.RUnlock()
	metrics := make([]Metric, len(r.metrics))
	copy(metrics, r.metrics)
	return metrics
}

func (r *Registry) Describe(descs chan<- *prometheus.Desc) {
	for _, metric := range r.GetMetrics() {
		if collector, ok := metric.(prometheus.Collector); ok {
			collector.Describe(descs)
		}
	}
}

func (r *Registry) Collect(ch chan<- prometheus.Metric) {
	for _, m := range r.GetMetrics() {
		m.Collect(ch)
	}
}
