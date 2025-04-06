package config

import (
	"fmt"
	"os"
	"runtime"

	"github.com/sirupsen/logrus"
)

// Config holds the agent configuration
type Config struct {
	LogLevel      string
	ServerURL     string
	WSEndpoint    string
	DeviceID      string
	AuthToken     string
	CheckInterval int // in seconds
	Debug         bool
}

// Load loads the configuration
func Load(log *logrus.Logger) (*Config, error) {
	// On non-OpenWrt systems, we need to provide default config
	if !isOpenWrt() {
		log.Warn("Not running on OpenWrt - using default configuration")
		return defaultConfig(), nil
	}

	// Load UCI configuration
	log.Info("Loading configuration from UCI")
	uciCfg := NewUCIConfig(log)

	if err := uciCfg.Load(); err != nil {
		return nil, fmt.Errorf("failed to load UCI config: %v", err)
	}

	// Convert UCI config to standard config
	config := uciCfg.ToStandardConfig()

	// Validate and set defaults for any missing required fields
	if config.DeviceID == "" {
		hostname, err := os.Hostname()
		if err == nil {
			config.DeviceID = hostname
			log.Infof("DeviceID not configured, using hostname: %s", hostname)
		} else {
			return nil, fmt.Errorf("no device ID configured and couldn't get hostname")
		}
	}

	if config.CheckInterval <= 0 {
		config.CheckInterval = 60
		log.Info("Invalid check interval, using default of 60 seconds")
	}

	return config, nil
}

// defaultConfig provides defaults for non-OpenWrt systems (development only)
func defaultConfig() *Config {
	// Set device ID to hostname
	deviceID := "unknown"
	hostname, err := os.Hostname()
	if err == nil {
		deviceID = hostname
	}

	return &Config{
		LogLevel:      "info",
		ServerURL:     "https://prv.cloud.ispapp.co",
		WSEndpoint:    "wss://prv.cloud.ispapp.co/agent/ws",
		DeviceID:      deviceID,
		CheckInterval: 60,
		Debug:         false,
	}
}

// isOpenWrt checks if we're running on OpenWrt
func isOpenWrt() bool {
	// Check if on Linux
	if runtime.GOOS != "linux" {
		return false
	}

	// Check for OpenWrt-specific files
	_, err := os.Stat("/etc/openwrt_release")
	return err == nil
}
