package config

import (
	"time"

	"github.com/spf13/viper"
)

const defaulCacheValue = 70

// Config configuration
type Config struct {
	Kubectl         Kubectl `mapstructure:"kubectl"`
	CacheTimeSecond int     `mapstructure:"cache_time_second"`
}

// Kubectl configuration kubectl and plugin path
type Kubectl struct {
	Bin         string   `mapstructure:"kubectl_absolute_path"`
	PluginPaths []string `mapstructure:"plugin_paths"`
}

// New return alfred k8s configuration
func New() (*Config, error) {
	v := viper.New()
	v.SetConfigType("yaml")
	v.SetConfigName(".alfred-k8s")
	v.AddConfigPath(".")
	v.AddConfigPath("$HOME/")

	if err := v.ReadInConfig(); err != nil {
		// ignore not found error. try to exec default options
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return &Config{}, err
		}
		// return default value
		// TODO v.SetDefault
		return &Config{
			Kubectl: Kubectl{
				Bin:         "/usr/local/bin/kubectl",
				PluginPaths: []string{"/usr/local/bin/"},
			},
			CacheTimeSecond: defaulCacheValue,
		}, nil
	}

	var c Config
	if err := v.Unmarshal(&c); err != nil {
		return &Config{}, err
	}

	return &c, nil
}

func (cfg *Config) TTL() time.Duration {
	var cacheTime time.Duration
	maxAge := cfg.CacheTimeSecond
	switch {
	case maxAge == 0:
		cacheTime = defaulCacheValue * time.Second
	case maxAge < 0:
		cacheTime = 0 * time.Second
	default:
		cacheTime = time.Duration(maxAge) * time.Second
	}
	return cacheTime
}
