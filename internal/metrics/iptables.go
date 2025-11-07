package metrics

import (
	"iptables_exporter/internal/exporter"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

func init() {
	exporter.Register(NewIptablesCollector())
}

type IptablesCollector struct {
	scrapeDurantion *prometheus.Desc
	scrapeSuccess   *prometheus.Desc
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
	}
}

func (c *IptablesCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.scrapeDurantion
	ch <- c.scrapeSuccess
}

func (c *IptablesCollector) Collect(ch chan<- prometheus.Metric) {
	start := time.Now()
	success := 1

	// 发送基础指标
	ch <- prometheus.MustNewConstMetric(c.scrapeDurantion, prometheus.GaugeValue, time.Since(start).Seconds())
	ch <- prometheus.MustNewConstMetric(c.scrapeSuccess, prometheus.GaugeValue, float64(success))

}
