// SPDX-FileCopyrightText: 2025 UnionTech Software Technology Co., Ltd.
// SPDX-License-Identifier: MIT

package config

import "iptables_exporter/pkg/utils"

// Settings 保留为未来扩展预留，目前未在主流程中使用。
// 如需启用基于 HTTP 的 scrape，可在此基础上补充逻辑，并在 main 中显式调用。
type Settings struct {
	ScrapeUri string `yaml:"scrape_uri"`
	Insecure  bool   `yaml:"insecure"`
}

// ValidateSettings 用于在显式解析配置后进行合法性校验。
// 当前仅示例性地校验 ScrapeUri，如不需要可在调用层忽略。
func ValidateSettings(s *Settings) {
	if s == nil {
		return
	}
	if s.ScrapeUri == "" {
		return
	}
	if err := utils.ValidateURI(s.ScrapeUri); err != nil {
		// 调用方可以决定是否输出日志，这里只修正为默认值的占位逻辑
		// s.ScrapeUri = DefaultSettings.ScrapeUri
	}
}
