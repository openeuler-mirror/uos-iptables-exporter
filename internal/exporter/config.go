// SPDX-FileCopyrightText: 2025 UnionTech Software Technology Co., Ltd.
// SPDX-License-Identifier: MIT

package exporter

import (
	"fmt"
	"os"
	"strings"
	"time"

	"iptables_exporter/pkg/logger"
	"iptables_exporter/pkg/utils"

	"github.com/alecthomas/kingpin"
	"gopkg.in/yaml.v2"
)

var (
	Configfile = kingpin.Flag("config", "Configuration file").
			Short('c').
			Default("/etc/uos-exporter/iptables-exporter.yaml").
			String()

	DefaultConfig = Config{
		Logging: logger.Config{
			Level:   "debug",
			LogPath: "/var/log/exporter.log",
			MaxSize: "10MB",
			MaxAge:  time.Hour * 24 * 7,
		},
		Address:                  "127.0.0.1",
		Port:                     9077,
		MetricsPath:              "/metrics",
		ScrapeTimeoutSeconds:     20,
		MinScrapeIntervalSeconds: 10,
	}
)

type Config struct {
	Logging     logger.Config `yaml:"log"`
	Address     string        `yaml:"address"`
	Port        int           `yaml:"port"`
	MetricsPath string        `yaml:"metricsPath"`
	// ScrapeTimeoutSeconds defines the maximum duration (in seconds) allowed for
	// a single iptables scrape before it is treated as failed.
	ScrapeTimeoutSeconds int `yaml:"scrapeTimeoutSeconds"`
	// MinScrapeIntervalSeconds defines the minimum interval (in seconds) between
	// two real iptables scrapes. If Prometheus scrapes more frequently than this
	// value, cached results from the last successful scrape may be reused.
	MinScrapeIntervalSeconds int `yaml:"minScrapeIntervalSeconds"`
}

func Unpack(config interface{}) error {
	if !utils.FileExists(*Configfile) {
		return fmt.Errorf("config file %s not found", *Configfile)
	}
	file, err := os.Open(*Configfile)
	if err != nil {
		return fmt.Errorf("failed to open config file %s: %w", *Configfile, err)
	}
	defer file.Close()

	if err := yaml.NewDecoder(file).Decode(config); err != nil {
		return fmt.Errorf("failed to decode config file %s: %w", *Configfile, err)
	}
	return nil
}

// LoadConfig 解析命令行参数并加载配置文件，返回最终配置。
//
// 使用方式：
//
//	kingpin.HelpFlag.Short('h')
//	kingpin.Parse()
//	cfg, err := exporter.LoadConfig()
func LoadConfig() (Config, error) {
	// 先以默认值初始化
	cfg := DefaultConfig

	if err := Unpack(&cfg); err != nil {
		return Config{}, fmt.Errorf("failed to unpack config: %w", err)
	}
	if err := validateConfig(cfg); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func validateConfig(cfg Config) error {
	if cfg.Port <= 0 || cfg.Port > 65535 {
		return fmt.Errorf("invalid port %d", cfg.Port)
	}
	if cfg.MetricsPath == "" || !strings.HasPrefix(cfg.MetricsPath, "/") {
		return fmt.Errorf("invalid metricsPath %q", cfg.MetricsPath)
	}
	if cfg.ScrapeTimeoutSeconds <= 0 {
		return fmt.Errorf("invalid scrapeTimeoutSeconds %d", cfg.ScrapeTimeoutSeconds)
	}
	if cfg.MinScrapeIntervalSeconds <= 0 {
		return fmt.Errorf("invalid minScrapeIntervalSeconds %d", cfg.MinScrapeIntervalSeconds)
	}

	return nil
}
