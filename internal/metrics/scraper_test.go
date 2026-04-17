// SPDX-FileCopyrightText: 2025 UnionTech Software Technology Co., Ltd.
// SPDX-License-Identifier: MIT

package metrics

import (
	"context"
	"testing"

	"iptables_exporter/internal/exporter"
	"iptables_exporter/internal/parser"
)

func TestScrapeNowUsesUpdatedGetter(t *testing.T) {
	cfg := exporter.DefaultConfig
	cfg.ScrapeTimeoutSeconds = 1
	scraper := NewIptablesScraper(cfg)

	firstTables := parser.Tables{
		"filter": &parser.Table{Name: "filter", Chains: map[string]parser.Chain{}},
	}
	secondTables := parser.Tables{
		"nat": &parser.Table{Name: "nat", Chains: map[string]parser.Chain{}},
	}

	scraper.SetGetTables(func(ctx context.Context) (parser.Tables, error) {
		return firstTables, nil
	})
	scraper.ScrapeNow(context.Background())

	tables, lastOK, lastErr, _, _, _, _ := scraper.Snapshot()
	if !lastOK {
		t.Fatalf("expected initial scrape to succeed, got error %q", lastErr)
	}
	if _, ok := tables["filter"]; !ok {
		t.Fatalf("expected initial scrape to use first getter, got tables %v", tables)
	}

	scraper.SetGetTables(func(ctx context.Context) (parser.Tables, error) {
		return secondTables, nil
	})
	scraper.ScrapeNow(context.Background())

	tables, lastOK, lastErr, _, _, _, _ = scraper.Snapshot()
	if !lastOK {
		t.Fatalf("expected updated scrape to succeed, got error %q", lastErr)
	}
	if _, ok := tables["nat"]; !ok {
		t.Fatalf("expected updated scrape to use second getter, got tables %v", tables)
	}
	if _, ok := tables["filter"]; ok {
		t.Fatalf("expected latest snapshot to replace old tables, got tables %v", tables)
	}
}
