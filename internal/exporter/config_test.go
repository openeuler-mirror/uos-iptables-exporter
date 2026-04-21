// SPDX-FileCopyrightText: 2025 UnionTech Software Technology Co., Ltd.
// SPDX-License-Identifier: MIT

package exporter

import "testing"

func TestValidateConfig(t *testing.T) {
	tests := []struct {
		name    string
		mutate  func(*Config)
		wantErr bool
	}{
		{
			name:    "default config",
			mutate:  func(*Config) {},
			wantErr: false,
		},
		{
			name: "invalid port",
			mutate: func(cfg *Config) {
				cfg.Port = 0
			},
			wantErr: true,
		},
		{
			name: "invalid metrics path",
			mutate: func(cfg *Config) {
				cfg.MetricsPath = "metrics"
			},
			wantErr: true,
		},
		{
			name: "invalid scrape timeout",
			mutate: func(cfg *Config) {
				cfg.ScrapeTimeoutSeconds = -1
			},
			wantErr: true,
		},
		{
			name: "invalid scrape interval",
			mutate: func(cfg *Config) {
				cfg.MinScrapeIntervalSeconds = 0
			},
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cfg := DefaultConfig
			test.mutate(&cfg)

			err := validateConfig(cfg)
			if (err != nil) != test.wantErr {
				t.Fatalf("validateConfig() error = %v, wantErr %v", err, test.wantErr)
			}
		})
	}
}