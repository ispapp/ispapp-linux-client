package utils

import (
	"github.com/spf13/viper"
)

// FileSyncConfig represents a single sync configuration
type FileSyncConfig struct {
	Remote string `mapstructure:"remote"`
	Local  string `mapstructure:"local"`
}

// Config holds all file sync configurations
type Config struct {
	SyncPaths []FileSyncConfig `mapstructure:"sync_paths"`
}

// readConfig reads the configuration from a file using Viper
func ReadConfig(filename string) (*Config, error) {
	viper.SetConfigFile(filename)
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
