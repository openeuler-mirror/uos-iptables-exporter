// SPDX-FileCopyrightText: 2025 UnionTech Software Technology Co., Ltd.
// SPDX-License-Identifier: MIT

package metrics

import (
	"context"
	"sync"
	"time"

	"iptables_exporter/internal/exporter"
	"iptables_exporter/internal/parser"

	"github.com/prometheus/client_golang/prometheus"
)

func init() {
	exporter.Register(NewIptablesCollector())
}

var (
	lastScrapeSuccessState = true
	lastScrapeErrorMessage = ""
)

// LastScrapeOK returns whether the last iptables scrape succeeded.
func LastScrapeOK() bool {
	return lastScrapeSuccessState
}

// LastScrapeError returns the error message of the last failed scrape (empty if none).
func LastScrapeError() string {
	return lastScrapeErrorMessage
}

type IptablesCollector struct {
	mu                  sync.Mutex
	lastTables          parser.Tables
	lastScrapeTime      time.Time
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

func NewIptablesCollector() *IptablesCollector {
	return &IptablesCollector{
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
			[]string{"table", "chain", "in_iface", "out_iface", "rule"},
			nil,
		),
		rulePacketsTotal: prometheus.NewDesc(
			"iptables_rule_packets_total",
			"Total packets matching a rule.",
			[]string{"table", "chain", "in_iface", "out_iface", "rule"},
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

	start := time.Now()
	success := 1.0
	scrapeErrors := 0.0
	// 每次 Collect 都增加一次抓取总数
	ch <- prometheus.MustNewConstMetric(c.scrapeTotal, prometheus.CounterValue, 1)

	// 带超时的数据采集：使用配置中的 ScrapeTimeoutSeconds
	timeout := time.Duration(exporter.DefaultConfig.ScrapeTimeoutSeconds) * time.Second
	if timeout <= 0 {
		timeout = 5 * time.Second
	}

	// 结果缓存：根据 MinScrapeIntervalSeconds 决定是否复用上一次抓取结果
	minInterval := time.Duration(exporter.DefaultConfig.MinScrapeIntervalSeconds) * time.Second
	var tables parser.Tables

	// 如果设置了最小抓取间隔且上次抓取时间在该间隔内，复用缓存
	if minInterval > 0 && !c.lastScrapeTime.IsZero() && time.Since(c.lastScrapeTime) < minInterval && lastScrapeSuccessState {
		tables = c.lastTables
	} else {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		var err error
		tables, err = parser.GetTables(ctx)
		if err != nil {
			success = 0
			scrapeErrors = 1
			lastScrapeSuccessState = false
			lastScrapeErrorMessage = err.Error()
			// 记录抓取错误总数以及最近一次抓取失败
			ch <- prometheus.MustNewConstMetric(c.scrapeErrorsTotal, prometheus.CounterValue, scrapeErrors)
			ch <- prometheus.MustNewConstMetric(c.scrapeSuccess, prometheus.GaugeValue, success)
			ch <- prometheus.MustNewConstMetric(c.lastScrapeSuccess, prometheus.GaugeValue, 0)
			// 失败时不更新时间戳，只保留最近一次成功的时间戳
			return
		}

		// 成功抓取后更新缓存
		c.lastTables = tables
		c.lastScrapeTime = time.Now()
	}

	// 发送基础指标
	ch <- prometheus.MustNewConstMetric(c.scrapeDurantion, prometheus.GaugeValue, time.Since(start).Seconds())
	ch <- prometheus.MustNewConstMetric(c.scrapeSuccess, prometheus.GaugeValue, success)
	ch <- prometheus.MustNewConstMetric(c.lastScrapeSuccess, prometheus.GaugeValue, 1)
	ch <- prometheus.MustNewConstMetric(c.lastScrapeTimestamp, prometheus.GaugeValue, float64(time.Now().Unix()))
	lastScrapeSuccessState = true
	lastScrapeErrorMessage = ""

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
					rule.RuleSpec,
				)
				ch <- prometheus.MustNewConstMetric(
					c.ruleBytesTotal,
					prometheus.CounterValue,
					float64(rule.Bytes),
					table.Name,
					chain.Name,
					rule.InIface,
					rule.OutIface,
					rule.RuleSpec,
				)
			}
		}
	}
}
