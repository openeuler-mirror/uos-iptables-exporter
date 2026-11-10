// SPDX-FileCopyrightText: 2025 UnionTech Software Technology Co., Ltd.
// SPDX-License-Identifier: MIT


# uos-iptables-exporter

## 简介

`uos-iptables-exporter` 是一个 Prometheus Exporter，用于从本机 `iptables`
中采集包/字节统计信息，并通过 HTTP 暴露为监控指标。

Exporter 会调用 `iptables-save -c`，解析输出中的表 / 链 / 规则级别统计数据，
方便你在 Prometheus / Grafana 中对防火墙流量做监控和告警。

下面主要说明如何 **编译**、**部署** 和 **使用** 这个程序。

## 编译

### 环境要求

- Go 1.22 及以上版本；
- Linux 环境（需要能调用 `iptables-save`）。

### 从源码构建

在项目根目录 `uos-iptables-exporter` 下执行：

- 使用 `make fmt` 格式化 Go 代码；
- 使用 `make test` 运行测试；
- 使用 `make vet` 做静态检查；
- 使用 `make build`（推荐）；或使用 `go build ./...`。

使用 `make build` 时，编译好的二进制通常位于：

- `build/bin/iptables-exporter`

## 部署

### 安装到系统（可选）

`Makefile` 提供了 `install` 目标，执行后会：

1. 安装二进制到：`/usr/local/bin/iptables-exporter`（或 `${PREFIX}/bin`）；
2. 安装默认配置到：`/etc/uos-exporter/iptables-exporter.yaml`；
3. 安装 systemd unit 到：`/usr/lib/systemd/system/uos-iptables-exporter.service`。

安装完成后可以通过 systemd 管理服务，例如：

```bash
systemctl enable --now uos-iptables-exporter.service
```

### 配置文件

默认配置文件路径为：`/etc/uos-exporter/iptables-exporter.yaml`。

一个最小可用的配置示例：

```yaml
address: "0.0.0.0"
port: 9077
metricsPath: "/metrics"
log:
  level: "info"
```

常用字段：

- `address`：Exporter 监听地址，例如 `0.0.0.0`；
- `port`：监听端口，例如 `9077`；
- `metricsPath`：Prometheus 抓取指标的路径，默认 `/metrics`；
- `log.level`：日志级别，`debug` / `info` / `warn` 等；
- `scrapeTimeoutSeconds`：单次执行 `iptables-save -c` 的超时时间，默认 `20` 秒；
- `minScrapeIntervalSeconds`：后台真实抓取的最小间隔，默认 `10` 秒，用于避免高频重复执行 `iptables-save`。

> 注意：如果配置文件不存在或格式错误，程序会在启动时直接报错退出。

## 使用

### 命令行启动

在安装好二进制并准备好配置文件后，可以通过命令行直接启动：

```bash
iptables-exporter -c /etc/uos-exporter/iptables-exporter.yaml
```

常用参数：

- `-c, --config`：指定 YAML 配置文件路径；
- `-h`：显示命令行帮助。

### 访问指标和健康检查

假设配置为：`address: 0.0.0.0`、`port: 9077`、`metricsPath: /metrics`，则：

- 指标地址：`http://<host>:9077/metrics`
- 健康检查：`http://<host>:9077/healthz`

你可以在浏览器或 `curl` 中访问这些地址，确认 Exporter 是否正常运行。

### Prometheus 抓取配置示例

在 Prometheus 的配置文件中添加类似配置：

```yaml
scrape_configs:
  - job_name: 'iptables'
    static_configs:
      - targets: ['127.0.0.1:9077']
```

重新加载 Prometheus 配置后，就可以在 Prometheus / Grafana 中查看规则和链的
包数、字节数等相关指标了。

## 指标列表

Exporter 暴露的主要指标如下（完整定义见 `internal/metrics/iptables.go`）：

| 指标名 | 类型 | 标签 |
| --- | --- | --- |
| `iptables_exporter_scrape_duration_seconds` | Gauge | — |
| `iptables_exporter_scrape_success` | Gauge | — |
| `iptables_exporter_scrape_total` | Counter | — |
| `iptables_exporter_scrape_errors_total` | Counter | — |
| `iptables_exporter_last_scrape_success` | Gauge | — |
| `iptables_exporter_last_scrape_timestamp_seconds` | Gauge | — |
| `iptables_table_chains_total` | Gauge | `table` |
| `iptables_chain_rules_total` | Gauge | `table`, `chain` |
| `iptables_chain_empty` | Gauge | `table`, `chain` |
| `iptables_default_bytes_total` | Counter | `table`, `chain`, `policy` |
| `iptables_default_packets_total` | Counter | `table`, `chain`, `policy` |
| `iptables_rule_bytes_total` | Counter | `table`, `chain`, `in_iface`, `out_iface`, `source`, `destination`, `target`, `protocol` |
| `iptables_rule_packets_total` | Counter | `table`, `chain`, `in_iface`, `out_iface`, `source`, `destination`, `target`, `protocol` |
