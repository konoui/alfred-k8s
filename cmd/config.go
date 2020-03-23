package cmd

import (
	"github.com/spf13/viper"
)

// Config configuration
type Config struct {
	Kubectl Kubectl `mapstructure:"kubectl"`
}

// Kubectl configuration kubectl and plugin path
type Kubectl struct {
	Bin         string   `mapstructure:"kubectl_absolute_path"`
	PluginPaths []string `mapstructure:"plugin_paths"`
}

// NewConfig return alfred k8s configuration
func newConfig() (*Config, error) {
	var c Config
	viper.SetConfigType("yaml")
	viper.SetConfigName(".alfred-k8s")
	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME/")

	// Set Default Value overwritten with config file
	viper.SetDefault("kubectl.bin", "/usr/local/bin/kubectl")
	viper.SetDefault("kubectl.pluginPath", "/usr/local/bin/")
	if err := viper.ReadInConfig(); err != nil {
		// ignore not found error. try to exec default options
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return &Config{}, err
		}
	}

	if err := viper.Unmarshal(&c); err != nil {
		return &Config{}, err
	}

	return &c, nil
}
