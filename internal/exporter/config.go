package exporter


import (
	"fmt"
	"os"
	"time"

	"iptables_exporter/pkg/logger"
	"iptables_exporter/pkg/utils"

	"github.com/alecthomas/kingpin"
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
