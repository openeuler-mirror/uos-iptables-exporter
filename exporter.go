package main

import (
	"iptables_exporter/internal/exporter"
	"iptables_exporter/internal/server"
	"iptables_exporter/pkg/logger"

	"github.com/sirupsen/logrus"
)

// Run is the main entry point for starting the exporter server.
// It initializes logging, loads configuration once via exporter.LoadConfig,
// then constructs and runs the HTTP server.
func Run(name string, version string) error {
	logger.InitDefaultLog()

	// 统一加载配置：优先从 YAML 读取，失败时直接返回错误。
	commonCfg, err := exporter.LoadConfig()
	if err != nil {
		logrus.Errorf("LoadConfig error: %v", err)
		return err
	}

	s := server.NewServer(name, version)
	s.CommonConfig = commonCfg

	s.PrintVersion()
	err = s.SetUp()
	if err != nil {
		logrus.Errorf("SetUp error: %v", err)
		return err
	}
	go func() {
		err := s.Run()
		if err != nil {
			logrus.Errorf("Run error: %v", err)
			s.Error = err
		}

		s.Exit()
	}()
	select {
	case <-s.ExitSignal:
		s.Stop()
		logrus.Info("Exit exporter server completed")
		return s.Error
	}
}


