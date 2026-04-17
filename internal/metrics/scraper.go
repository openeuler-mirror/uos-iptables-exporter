// SPDX-FileCopyrightText: 2025 UnionTech Software Technology Co., Ltd.
// SPDX-License-Identifier: MIT

package metrics

import (
	"context"
	"sync"
	"time"

	"iptables_exporter/internal/exporter"
	"iptables_exporter/internal/parser"
)

// IptablesScraper runs iptables scrapes in the background and keeps the latest
// snapshot for collectors to read without blocking Prometheus scrapes.
type IptablesScraper struct {
	cfg         exporter.Config
	getTables   func(ctx context.Context) (parser.Tables, error)
	getTablesMu sync.RWMutex

	mu           sync.RWMutex
	lastTables   parser.Tables
	lastOK       bool
	lastErr      string
	lastDuration time.Duration
	lastAttempt  time.Time
	lastSuccess  time.Time
	scrapeTotal  uint64
	scrapeErrors uint64

	inflightMu sync.Mutex
	inflight   bool
}

func NewIptablesScraper(cfg exporter.Config) *IptablesScraper {
	return &IptablesScraper{
		cfg:       cfg,
		getTables: parser.GetTables,
		lastOK:    true,
	}
}

// SetGetTables overrides the table getter (primarily for tests).
func (s *IptablesScraper) SetGetTables(getTables func(ctx context.Context) (parser.Tables, error)) {
	s.getTablesMu.Lock()
	s.getTables = getTables
	s.getTablesMu.Unlock()
}

func (s *IptablesScraper) currentGetTables() func(ctx context.Context) (parser.Tables, error) {
	s.getTablesMu.RLock()
	defer s.getTablesMu.RUnlock()
	return s.getTables
}

func (s *IptablesScraper) LastOK() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.lastOK
}

func (s *IptablesScraper) LastError() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.lastErr
}

// Snapshot returns the latest scrape snapshot. Treat returned tables as read-only.
func (s *IptablesScraper) Snapshot() (tables parser.Tables, lastOK bool, lastErr string, lastSuccess time.Time, scrapeTotal, scrapeErrors uint64, lastDuration time.Duration) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.lastTables, s.lastOK, s.lastErr, s.lastSuccess, s.scrapeTotal, s.scrapeErrors, s.lastDuration
}

// Start begins periodic background scraping. It performs an initial scrape
// asynchronously, then scrapes every MinScrapeIntervalSeconds.
func (s *IptablesScraper) Start(ctx context.Context) {
	interval := time.Duration(s.cfg.MinScrapeIntervalSeconds) * time.Second
	if interval <= 0 {
		interval = 10 * time.Second
	}

	// Kick once immediately.
	s.kick(ctx)

	ticker := time.NewTicker(interval)
	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				s.kick(ctx)
			}
		}
	}()
}

// ScrapeNow performs a synchronous scrape. Intended for tests.
func (s *IptablesScraper) ScrapeNow(ctx context.Context) {
	s.runScrape(ctx)
}

func (s *IptablesScraper) kick(ctx context.Context) {
	s.inflightMu.Lock()
	if s.inflight {
		s.inflightMu.Unlock()
		return
	}
	s.inflight = true
	s.inflightMu.Unlock()

	go func() {
		defer func() {
			s.inflightMu.Lock()
			s.inflight = false
			s.inflightMu.Unlock()
		}()
		s.runScrape(ctx)
	}()
}

func (s *IptablesScraper) runScrape(ctx context.Context) {
	timeout := time.Duration(s.cfg.ScrapeTimeoutSeconds) * time.Second
	if timeout <= 0 {
		timeout = 5 * time.Second
	}

	s.mu.Lock()
	s.scrapeTotal++
	s.lastAttempt = time.Now()
	s.mu.Unlock()

	scrapeCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	start := time.Now()
	get := s.currentGetTables()
	if get == nil {
		get = parser.GetTables
	}
	tables, err := get(scrapeCtx)
	duration := time.Since(start)

	s.mu.Lock()
	defer s.mu.Unlock()
	s.lastDuration = duration

	if err != nil {
		s.scrapeErrors++
		s.lastOK = false
		s.lastErr = err.Error()
		return
	}

	s.lastTables = tables
	s.lastOK = true
	s.lastErr = ""
	s.lastSuccess = time.Now()
}
