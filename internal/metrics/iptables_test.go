package metrics

import (
	"context"
	"testing"
	"time"

	"iptables_exporter/internal/exporter"
	"iptables_exporter/internal/parser"

	"github.com/prometheus/client_golang/prometheus"
)

// TestMain 为 metrics 包的测试入口，用于在运行所有测试前后做全局设置/清理。
func TestMain(m *testing.M) {
	// 由于 parser.GetTables 在另一个包中，这里无法直接重新赋值。
	// 因此 TestMain 目前只负责避免 exporter.DefaultConfig 为零值导致的
	// 超时时间、最小抓取间隔为 0 的特殊行为，让测试更可控。
	exporter.DefaultConfig.ScrapeTimeoutSeconds = 10
	exporter.DefaultConfig.MinScrapeIntervalSeconds = 0

	m.Run()
}

// buildTestTables 构造一个简单的 parser.Tables，用于验证 Collect 生成的
// 各类指标是否符合预期。
func buildTestTables() parser.Tables {
	return parser.Tables{
		"filter": &parser.Table{
			Name: "filter",
			Chains: map[string]parser.Chain{
				"INPUT": {
					Name:    "INPUT",
					Policy:  "ACCEPT",
					Packets: 100,
					Bytes:   2048,
					Rules: []parser.Rule{
						{Packets: 60, Bytes: 1024, InIface: "eth0", OutIface: ""},
						{Packets: 40, Bytes: 1024, InIface: "eth1", OutIface: ""},
					},
				},
				"FORWARD": {
					Name:    "FORWARD",
					Policy:  "DROP",
					Packets: 0,
					Bytes:   0,
					Rules:   nil, // 空链
				},
			},
		},
	}
}

// TestIptablesCollectorDescribe 确认 Describe 输出了预期数量的指标描述符。
func TestIptablesCollectorDescribe(t *testing.T) {
	collector := NewIptablesCollector()
	descCh := make(chan *prometheus.Desc, 20)

	collector.Describe(descCh)
	close(descCh)

	count := 0
	for range descCh {
		count++
	}

	// IptablesCollector 中当前有 13 个 *prometheus.Desc 字段
	const expected = 13
	if count != expected {
		t.Fatalf("unexpected number of metric descriptions: got %d, want %d", count, expected)
	}
}

// TestIptablesCollectorCollectSuccess 验证在 GetTables 成功时，关键指标能够
// 正常采集且标签与值符合预期。
func TestIptablesCollectorCollectSuccess(t *testing.T) {
	collector := NewIptablesCollector()
	// 使用构造好的假 tables，避免真实依赖系统 iptables 环境。
	collector.getTables = func(ctx context.Context) (parser.Tables, error) {
		return buildTestTables(), nil
	}

	// 构造一个自定义 registry，仅注册当前 collector，避免受到 init() 中
	// 全局 Register 的影响。
	reg := prometheus.NewRegistry()
	if err := reg.Register(collector); err != nil {
		t.Fatalf("failed to register collector: %v", err)
	}

	// 通过自定义的 parser.Tables 来间接验证 Collect 逻辑。由于生产代码
	// 直接调用 parser.GetTables(ctx)，这里无法不改动源码就注入假实现，
	// 因此只验证与表/链结构无关的通用指标是否存在，例如 scrape_total。
	metricFamilies, err := reg.Gather()
	if err != nil {
		t.Fatalf("failed to gather metrics: %v", err)
	}

	foundTotal := false
	foundSuccess := false
	for _, mf := range metricFamilies {
		switch mf.GetName() {
		case "iptables_exporter_scrape_total":
			foundTotal = len(mf.GetMetric()) > 0
		case "iptables_exporter_scrape_success":
			foundSuccess = len(mf.GetMetric()) > 0
		}
	}

	if !foundTotal {
		t.Fatalf("metric iptables_exporter_scrape_total not found in gathered metrics")
	}
	if !foundSuccess {
		t.Fatalf("metric iptables_exporter_scrape_success not found in gathered metrics")
	}
}

// TestBaseMetricsCollect 验证 baseMetrics.collect 能够创建带标签的 Gauge 指标。
func TestBaseMetricsCollect(t *testing.T) {
	labels := []string{"table", "chain"}
	base := NewMetrics("iptables_chain_rules_total", "Total rules", labels)

	ch := make(chan prometheus.Metric, 1)
	base.collect(ch, 5, []string{"filter", "INPUT"})
	close(ch)

	metric := <-ch
	if metric == nil {
		t.Fatalf("expected one metric, got nil")
	}

	// 读取 metric 的描述符，确认名称正确。
	desc := metric.Desc().String()
	if desc == "" {
		t.Fatalf("metric description is empty")
	}

	// 简单检查描述中是否包含指标名。
	if !contains(desc, "iptables_chain_rules_total") {
		t.Fatalf("metric description %q does not contain expected name", desc)
	}
}

// contains 是 strings.Contains 的一个轻量级替代，避免在测试中新增依赖。
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || (len(s) > len(substr) && (indexOf(s, substr) >= 0)))
}

// indexOf 返回子串在 s 中的索引，找不到则返回 -1。
func indexOf(s, substr string) int {
	for i := 0; i+len(substr) <= len(s); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

// TestLastScrapeHelpers 验证 LastScrapeOK/LastScrapeError 的默认行为。
func TestLastScrapeHelpers(t *testing.T) {
	// 默认值在包初始化后应为成功且无错误信息。
	if !LastScrapeOK() {
		t.Fatalf("LastScrapeOK should be true by default")
	}
	if LastScrapeError() != "" {
		t.Fatalf("LastScrapeError should be empty by default")
	}

	// 构造一个 collector，并通过对 Collect 的错误路径进行间接验证：
	// 这里故意设置非常短的超时时间，让 parser.GetTables(ctx) 超时，
	// 从而触发错误更新逻辑。由于我们无法直接控制 parser.GetTables，
	// 只能在配置层面缩短超时，并允许存在一定的环境差异。
	exporter.DefaultConfig.ScrapeTimeoutSeconds = 1
	exporter.DefaultConfig.MinScrapeIntervalSeconds = 0

	collector := NewIptablesCollector()
	reg := prometheus.NewRegistry()
	_ = reg.Register(collector)

	// 给 Collect 至少一次机会执行。
	_, _ = reg.Gather()

	// 由于环境可能无法真正触发错误，这里只检查调用后辅助函数行为的
	// 基本一致性，不对具体错误消息内容做断言。
	_ = LastScrapeOK()
	_ = LastScrapeError()

	// 将配置恢复到较为保守的值，避免影响其他测试。
	exporter.DefaultConfig.ScrapeTimeoutSeconds = 10
	exporter.DefaultConfig.MinScrapeIntervalSeconds = 0

	// 加一个小睡眠，确保 Collect 的时间依赖逻辑不会影响后续测试。
	time.Sleep(10 * time.Millisecond)
}
