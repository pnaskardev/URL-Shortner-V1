package config

import (
	"fmt"
	"log/slog"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type Config struct {
	Port string `mapstructure:"app_port"`
}

var (
	config     *Config
	configOnce sync.Once
)

func LoadConfig() error {
	var loadErr error

	configOnce.Do(func() {
		v := viper.New()

		v.SetConfigType("env")
		v.AddConfigPath(".")
		// Read environment variables prefixed with APP_.
		v.SetEnvPrefix("APP")
		v.AutomaticEnv()

		if err := v.ReadInConfig(); err != nil {
			loadErr = fmt.Errorf("read config file: %w", err)
			return
		}
		fmt.Println("Using config:", v.ConfigFileUsed())

		config = &Config{}
		if err := v.Unmarshal(config); err != nil {
			loadErr = fmt.Errorf("unmarshal config: %w", err)
			return
		}

		v.WatchConfig()
		v.OnConfigChange(func(e fsnotify.Event) {
			slog.Info("config file changed", "file", e.Name)
			if err := v.Unmarshal(config); err != nil {
				slog.Error("failed to reload config", "error", err)
			}
		})
	})

	return loadErr
}

func GetConfig() *Config {
	err := LoadConfig()
	if err != nil {
		panic(err)
	}
	return config
}
