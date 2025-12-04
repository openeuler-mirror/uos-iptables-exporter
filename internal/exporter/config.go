package exporter


import (
	"fmt"
	"os"
	"time"

	"iptables_exporter/pkg/logger"
	"iptables_exporter/pkg/utils"

	"github.com/alecthomas/kingpin"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

var (
	Configfile    = kingpin.Flag("config", "Configuration file").
		Short('c').
		Default("/etc/uos-exporter/iptables-exporter.yaml").
		String()

	DefaultConfig = Config{
		Logging: logger.Config{
			Level:   "debug",
			LogPath: "/var/log/exporter.log",
			MaxSize: "10MB",
			MaxAge:  time.Hour * 24 * 7},
		Address:     "127.0.0.1",
		Port:        8080,
		MetricsPath: "/metrics",
}
)

type Config struct {
	Logging     logger.Config `yaml:"log"`
	Address     string        `yaml:"address"`
	Port        int           `yaml:"port"`
	MetricsPath string        `yaml:"metricsPath"`
}

func Unpack(config interface{}) error {
	if !utils.FileExists(*Configfile) {
		logrus.Warnf("%s file not found, using default config", *Configfile)
	} else {
		file, err := os.Open(*Configfile)
		if err != nil {
			logrus.Error("Failed to open config file: ", err)
			return err
		}
		err = yaml.NewDecoder(file).Decode(config)
		if err != nil {
			return err
		}
	}
	return nil
}

// LoadConfig 解析命令行参数并加载配置文件，返回最终配置。
//
// 使用方式：
//   kingpin.HelpFlag.Short('h')
//   kingpin.Parse()
//   cfg, err := exporter.LoadConfig()
func LoadConfig() (Config, error) {
	// 先以默认值初始化
	cfg := DefaultConfig

	if err := Unpack(&cfg); err != nil {
		return Config{}, fmt.Errorf("failed to unpack config: %w", err)
	}

	return cfg, nil
}
