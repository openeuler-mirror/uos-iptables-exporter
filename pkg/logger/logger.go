// SPDX-FileCopyrightText: 2025 UnionTech Software Technology Co., Ltd.
// SPDX-License-Identifier: MIT

package logger

import (
	"strings"
	"time"

	formatter "gitee.com/weidongkl/logrus-formatter"
	"github.com/sirupsen/logrus"
)

type Config struct {
	Level   string        `yaml:"level"`
	LogPath string        `yaml:"log_path"`
	MaxSize string        `yaml:"max_size"`
	MaxAge  time.Duration `yaml:"max_age"`
}

type fileLogConfig struct {
	FileRotator *FileRotator
	level       string
}

func NewConfig(level, logPath string, maxSize int64, maxAge time.Duration) fileLogConfig {
	return fileLogConfig{
		level:       level,
		FileRotator: NewFileRotator(logPath, maxSize, maxAge),
	}
}

func Init(config fileLogConfig) {
	if config.FileRotator == nil {
		logrus.SetOutput(logrus.StandardLogger().Out)
	} else {
		logrus.SetReportCaller(true)
		logrus.SetFormatter(&formatter.Formatter{})
		logrus.SetOutput(config.FileRotator)
	}
	level, ok := parseLevel(config.level)
	if !ok {
		logrus.SetLevel(logrus.WarnLevel)
		logrus.Warnf("unknown log level: %s, use default level: warn", strings.ToLower(config.level))
		logrus.Warnf("support level is [debug,info,warn,warning,error]")
		return
	}
	logrus.SetLevel(level)
}

func InitDefaultLog() {
	Init(fileLogConfig{
		level: "info",
	})
	logrus.SetFormatter(&logrus.TextFormatter{
		DisableTimestamp: true,
	})
}

func parseLevel(level string) (logrus.Level, bool) {
	switch strings.ToLower(level) {
	case "debug":
		return logrus.DebugLevel, true
	case "info":
		return logrus.InfoLevel, true
	case "warn", "warning":
		return logrus.WarnLevel, true
	case "error":
		return logrus.ErrorLevel, true
	default:
		return logrus.WarnLevel, false
	}
}
