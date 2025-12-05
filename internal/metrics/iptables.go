// SPDX-FileCopyrightText: 2025 UnionTech Software Technology Co., Ltd.
// SPDX-License-Identifier: MIT

package metrics

import (
	"context"
	"iptables_exporter/internal/exporter"
	"iptables_exporter/internal/parser"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

func init() {
	exporter.Register(NewIptablesCollector())
}

type IptablesCollector struct {
	scrapeDurantion     *prometheus.Desc
	scrapeSuccess       *prometheus.Desc
	scrapeTotal         *prometheus.Desc
	scrapeErrorsTotal   *prometheus.Desc
	lastScrapeSuccess   *prometheus.Desc
	lastScrapeTimestamp *prometheus.Desc
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
			[]string{"table", "chain", "rule"},
			nil,
		),
		rulePacketsTotal: prometheus.NewDesc(
			"iptables_rule_packets_total",
			"Total packets matching a rule.",
			[]string{"table", "chain", "rule"},
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
	ch <- c.defaultBytesTotal
	ch <- c.defaultPacketsTotal
	ch <- c.ruleBytesTotal
	ch <- c.rulePacketsTotal
}

func (c *IptablesCollector) Collect(ch chan<- prometheus.Metric) {
	start := time.Now()
	success := 1.0
	scrapeErrors := 0.0
	// 每次 Collect 都增加一次抓取总数
	ch <- prometheus.MustNewConstMetric(c.scrapeTotal, prometheus.CounterValue, 1)

	// 带超时的数据采集
	tables, err := parser.GetTables(context.Background())
	if err != nil {
		success = 0
		scrapeErrors = 1
		// 记录抓取错误总数以及最近一次抓取失败
		ch <- prometheus.MustNewConstMetric(c.scrapeErrorsTotal, prometheus.CounterValue, scrapeErrors)
		ch <- prometheus.MustNewConstMetric(c.scrapeSuccess, prometheus.GaugeValue, success)
		ch <- prometheus.MustNewConstMetric(c.lastScrapeSuccess, prometheus.GaugeValue, 0)
		// 失败时不更新时间戳，只保留最近一次成功的时间戳
		return
	}

	// 发送基础指标
	ch <- prometheus.MustNewConstMetric(c.scrapeDurantion, prometheus.GaugeValue, time.Since(start).Seconds())
	ch <- prometheus.MustNewConstMetric(c.scrapeSuccess, prometheus.GaugeValue, success)
	ch <- prometheus.MustNewConstMetric(c.lastScrapeSuccess, prometheus.GaugeValue, 1)
	ch <- prometheus.MustNewConstMetric(c.lastScrapeTimestamp, prometheus.GaugeValue, float64(time.Now().Unix()))

	// 发送详细指标
	for _, table := range tables {
		for _, chain := range table.Chains {
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
					rule.RuleSpec,
				)
				ch <- prometheus.MustNewConstMetric(
					c.ruleBytesTotal,
					prometheus.CounterValue,
					float64(rule.Bytes),
					table.Name,
					chain.Name,
					rule.RuleSpec,
				)
			}
		}
	}
}
