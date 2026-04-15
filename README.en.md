// SPDX-FileCopyrightText: 2025 UnionTech Software Technology Co., Ltd.
// SPDX-License-Identifier: MIT

# uos-iptables-exporter

## Description

`uos-iptables-exporter` is a Prometheus exporter that collects packet/byte
statistics from the local `iptables` rules and exposes them as metrics over HTTP.

It calls `iptables-save -c`, parses the tables/chains/rules statistics, and
lets you monitor firewall traffic in Prometheus / Grafana.

This document focuses on how to **build**, **deploy** and **use** the exporter.

## Build

### Requirements

- Go 1.22 or later
- Linux environment (must be able to run `iptables-save`)

### Build from source

From the project root directory `uos-iptables-exporter`:

1. Fetch dependencies and build:
   - `make build` (recommended), or
   - `go build ./...`.

When using `make build`, the binary is usually placed at:

- `build/bin/iptables-exporter`

## Deploy

### Install to system (optional)

The `Makefile` provides an `install` target which will:

1. Install the binary to `/usr/local/bin/iptables-exporter` (or `${PREFIX}/bin`).
2. Install the default config to `/etc/uos-exporter/iptables-exporter.yaml`.
3. Install the systemd unit to `/usr/lib/systemd/system/uos-iptables-exporter.service`.

After installation, you can manage the service via systemd, for example:

```bash
systemctl enable --now uos-iptables-exporter.service
```

### Configuration file

The default config file path is: `/etc/uos-exporter/iptables-exporter.yaml`.

A minimal working example:

```yaml
address: "0.0.0.0"
port: 9077
metricsPath: "/metrics"
log:
  level: "info"
```

Common fields:

- `address`: exporter listen address, e.g. `0.0.0.0`.
- `port`: listen port, e.g. `9077`.
- `metricsPath`: metrics endpoint path for Prometheus, default `/metrics`.
- `log.level`: log level, such as `debug` / `info` / `warn`.
- `scrapeTimeoutSeconds`: timeout for a single `iptables-save -c` execution, default `20` seconds.
- `minScrapeIntervalSeconds`: minimum interval between real background scrapes, default `10` seconds, to avoid running `iptables-save` too frequently.

> Note: if the config file is missing or invalid, the exporter will fail fast
> during startup instead of running with implicit defaults.

## Usage

### Run from command line

After the binary and config file are ready, you can start the exporter with:

```bash
iptables-exporter -c /etc/uos-exporter/iptables-exporter.yaml
```

Common options:

- `-c, --config`: path to the YAML config file.
- `-h`: show command line help.

### Access metrics and health check

Assuming `address: 0.0.0.0`, `port: 9077`, `metricsPath: /metrics`:

- Metrics: `http://<host>:9077/metrics`
- Health check: `http://<host>:9077/healthz`

You can open these URLs in a browser or with `curl` to verify the exporter is running.

### Prometheus scrape example

In your Prometheus configuration, add something like:

```yaml
scrape_configs:
  - job_name: 'iptables'
    static_configs:
      - targets: ['127.0.0.1:9077']
```

After reloading Prometheus, you should see iptables-related metrics (tables,
chains, rules, packets, bytes, etc.) in Prometheus / Grafana.
