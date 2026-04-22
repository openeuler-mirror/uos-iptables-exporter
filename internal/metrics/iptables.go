// SPDX-FileCopyrightText: 2025 UnionTech Software Technology Co., Ltd.
// SPDX-License-Identifier: MIT

package metrics

import (
	"sync"

	"iptables_exporter/internal/exporter"

	"github.com/prometheus/client_golang/prometheus"
)

func init() {
	// Collector registration is done by the server after YAML config is loaded.
}

var (
	lastScrapeStateMu      sync.RWMutex
	lastScrapeSuccessState = true
	lastScrapeErrorMessage = ""
)

// LastScrapeOK returns whether the last iptables scrape succeeded.
func LastScrapeOK() bool {
	lastScrapeStateMu.RLock()
	defer lastScrapeStateMu.RUnlock()
	return lastScrapeSuccessState
}

// LastScrapeError returns the error message of the last failed scrape (empty if none).
func LastScrapeError() string {
	lastScrapeStateMu.RLock()
	defer lastScrapeStateMu.RUnlock()
	return lastScrapeErrorMessage
}

func setLastScrapeState(ok bool, errMessage string) {
	lastScrapeStateMu.Lock()
	defer lastScrapeStateMu.Unlock()
	lastScrapeSuccessState = ok
	lastScrapeErrorMessage = errMessage
}

type IptablesCollector struct {
	mu                  sync.Mutex
	scraper             *IptablesScraper
	scrapeDurantion     *prometheus.Desc
	scrapeSuccess       *prometheus.Desc
	scrapeTotal         *prometheus.Desc
	scrapeErrorsTotal   *prometheus.Desc
	lastScrapeSuccess   *prometheus.Desc
	lastScrapeTimestamp *prometheus.Desc
	tableChainsTotal    *prometheus.Desc
	chainRulesTotal     *prometheus.Desc
	chainEmpty          *prometheus.Desc
	defaultBytesTotal   *prometheus.Desc
	defaultPacketsTotal *prometheus.Desc
	ruleBytesTotal      *prometheus.Desc
	rulePacketsTotal    *prometheus.Desc
}

// NewIptablesCollector is a backwards-compatible constructor.
//
// It does not start background scraping. Production code should prefer
// NewIptablesCollectorWithScraper and start the scraper in the server lifecycle.
func NewIptablesCollector(cfg exporter.Config) *IptablesCollector {
	return NewIptablesCollectorWithScraper(NewIptablesScraper(cfg))
}

func NewIptablesCollectorWithScraper(scraper *IptablesScraper) *IptablesCollector {
	return &IptablesCollector{
		scraper: scraper,
		scrapeDurantion: prometheus.NewDesc(
			"iptables_exporter_scrape_duration_seconds",
			"Duration of scraping iptables.",
			nil,
			nil,
		),
		scrapeSuccess: prometheus.NewDesc(
			"iptables_exporter_scrape_success",
			"Whether scraping iptables succeeded.",
			nil,
			nil,
		),
		scrapeTotal: prometheus.NewDesc(
			"iptables_exporter_scrape_total",
			"Total number of iptables scrape attempts.",
			nil,
			nil,
		),
		scrapeErrorsTotal: prometheus.NewDesc(
			"iptables_exporter_scrape_errors_total",
			"Total number of iptables scrape errors.",
			nil,
			nil,
		),
		lastScrapeSuccess: prometheus.NewDesc(
			"iptables_exporter_last_scrape_success",
			"Whether the last iptables scrape succeeded.",
			nil,
			nil,
		),
		lastScrapeTimestamp: prometheus.NewDesc(
			"iptables_exporter_last_scrape_timestamp_seconds",
			"Unix timestamp of the last successful iptables scrape.",
			nil,
			nil,
		),
		tableChainsTotal: prometheus.NewDesc(
			"iptables_table_chains_total",
			"Total number of chains in an iptables table.",
			[]string{"table"},
			nil,
		),
		chainRulesTotal: prometheus.NewDesc(
			"iptables_chain_rules_total",
			"Total number of rules in an iptables chain.",
			[]string{"table", "chain"},
			nil,
		),
		chainEmpty: prometheus.NewDesc(
			"iptables_chain_empty",
			"Whether an iptables chain has no rules (1 for empty, 0 otherwise).",
			[]string{"table", "chain"},
			nil,
		),
		defaultBytesTotal: prometheus.NewDesc(
			"iptables_default_bytes_total",
			"Total bytes matching a chain's default policy.",
			[]string{"table", "chain", "policy"},
			nil,
		),
		defaultPacketsTotal: prometheus.NewDesc(
			"iptables_default_packets_total",
			"Total packets matching a chain's default policy.",
			[]string{"table", "chain", "policy"},
			nil,
		),
		ruleBytesTotal: prometheus.NewDesc(
			"iptables_rule_bytes_total",
			"Total bytes matching a rule.",
			[]string{"table", "chain", "in_iface", "out_iface", "source", "destination", "target"},
			nil,
		),
		rulePacketsTotal: prometheus.NewDesc(
			"iptables_rule_packets_total",
			"Total packets matching a rule.",
			[]string{"table", "chain", "in_iface", "out_iface", "source", "destination", "target"},
			nil,
		),
	}
}

func (c *IptablesCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.scrapeDurantion
	ch <- c.scrapeSuccess
	ch <- c.scrapeTotal
	ch <- c.scrapeErrorsTotal
	ch <- c.lastScrapeSuccess
	ch <- c.lastScrapeTimestamp
	ch <- c.tableChainsTotal
	ch <- c.chainRulesTotal
	ch <- c.chainEmpty
	ch <- c.defaultBytesTotal
	ch <- c.defaultPacketsTotal
	ch <- c.ruleBytesTotal
	ch <- c.rulePacketsTotal
}

func (c *IptablesCollector) Collect(ch chan<- prometheus.Metric) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.scraper == nil {
		// No scraper configured; report exporter as unhealthy.
		ch <- prometheus.MustNewConstMetric(c.scrapeDurantion, prometheus.GaugeValue, 0)
		ch <- prometheus.MustNewConstMetric(c.scrapeSuccess, prometheus.GaugeValue, 0)
		ch <- prometheus.MustNewConstMetric(c.scrapeTotal, prometheus.CounterValue, 0)
		ch <- prometheus.MustNewConstMetric(c.scrapeErrorsTotal, prometheus.CounterValue, 0)
		ch <- prometheus.MustNewConstMetric(c.lastScrapeSuccess, prometheus.GaugeValue, 0)
		ch <- prometheus.MustNewConstMetric(c.lastScrapeTimestamp, prometheus.GaugeValue, 0)
		setLastScrapeState(false, "scraper not configured")
		return
	}

	tables, lastOK, lastErr, lastSuccess, scrapeTotal, scrapeErrors, lastDuration := c.scraper.Snapshot()

	success := 0.0
	if lastOK {
		success = 1.0
	}

	// Base metrics
	ch <- prometheus.MustNewConstMetric(c.scrapeDurantion, prometheus.GaugeValue, lastDuration.Seconds())
	ch <- prometheus.MustNewConstMetric(c.scrapeSuccess, prometheus.GaugeValue, success)
	ch <- prometheus.MustNewConstMetric(c.scrapeTotal, prometheus.CounterValue, float64(scrapeTotal))
	ch <- prometheus.MustNewConstMetric(c.scrapeErrorsTotal, prometheus.CounterValue, float64(scrapeErrors))
	ch <- prometheus.MustNewConstMetric(c.lastScrapeSuccess, prometheus.GaugeValue, success)
	lastOKUnix := 0.0
	if !lastSuccess.IsZero() {
		lastOKUnix = float64(lastSuccess.Unix())
	}
	ch <- prometheus.MustNewConstMetric(c.lastScrapeTimestamp, prometheus.GaugeValue, lastOKUnix)
	setLastScrapeState(lastOK, lastErr)

	// If the last scrape failed, keep behavior consistent: do not emit detailed metrics.
	if !lastOK {
		return
	}

	// Defensive: if we have no tables yet, nothing more to export.
	if tables == nil {
		return
	}

	// 发送详细指标
	for _, table := range tables {
		// 表级聚合：表中链的数量
		ch <- prometheus.MustNewConstMetric(
			c.tableChainsTotal,
			prometheus.GaugeValue,
			float64(len(table.Chains)),
			table.Name,
		)

		for _, chain := range table.Chains {
			// 链级聚合：规则数量和是否为空
			ruleCount := len(chain.Rules)
			isEmpty := 0.0
			if ruleCount == 0 {
				isEmpty = 1.0
			}

			ch <- prometheus.MustNewConstMetric(
				c.chainRulesTotal,
				prometheus.GaugeValue,
				float64(ruleCount),
				table.Name,
				chain.Name,
			)
			ch <- prometheus.MustNewConstMetric(
				c.chainEmpty,
				prometheus.GaugeValue,
				isEmpty,
				table.Name,
				chain.Name,
			)

			// 默认策略指标
			ch <- prometheus.MustNewConstMetric(
				c.defaultPacketsTotal,
				prometheus.CounterValue,
				float64(chain.Packets),
				table.Name,
				chain.Name,
				chain.Policy,
			)
			ch <- prometheus.MustNewConstMetric(
				c.defaultBytesTotal,
				prometheus.CounterValue,
				float64(chain.Bytes),
				table.Name,
				chain.Name,
				chain.Policy,
			)

			// 规则级别指标
			// log.Print(table)
			for _, rule := range chain.Rules {
				ch <- prometheus.MustNewConstMetric(
					c.rulePacketsTotal,
					prometheus.CounterValue,
					float64(rule.Packets),
					table.Name,
					chain.Name,
					rule.InIface,
					rule.OutIface,
					rule.Source,
					rule.Destination,
					rule.Target,
				)
				ch <- prometheus.MustNewConstMetric(
					c.ruleBytesTotal,
					prometheus.CounterValue,
					float64(rule.Bytes),
					table.Name,
					chain.Name,
					rule.InIface,
					rule.OutIface,
					rule.Source,
					rule.Destination,
					rule.Target,
				)
			}
		}
	}
}
