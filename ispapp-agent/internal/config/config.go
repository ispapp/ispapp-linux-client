package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// Config holds all configuration parameters
type Config struct {
	// Basic settings
	DeviceID      string
	ServerURL     string
	ListenerPort  string
	LogLevel      string
	
	// Authentication
	APIKey        string
	AccessToken   string
	RefreshToken  string
	
	// Operational settings
	UpdateInterval int
	Enabled        bool
	
	// Network testing
	IperfServer   string
	PingTargets   []string
	
	// UI settings
	HighRefresh    bool
	RefreshInterval int
}

// DefaultConfig returns a configuration with reasonable defaults
func DefaultConfig() *Config {
	return &Config{
		DeviceID:       uuid.New().String(),
		ServerURL:      "prv.cloud.ispapp.co",
		ListenerPort:   "443",
		LogLevel:       "info",
		UpdateInterval: 1,
		Enabled:        true,
		PingTargets:    []string{"cloud.ispapp.co"},
	}
}

// Load reads configuration from UCI or file
func Load(log *logrus.Logger) (*Config, error) {
	// Start with defaults
	cfg := DefaultConfig()
	
	// Try to load from UCI first (OpenWrt)
	if uciCfg, err := loadFromUCI(); err == nil {
		log.Info("Loaded configuration from UCI")
		return uciCfg, nil
	} else {
		log.Debug("Could not load from UCI, trying file-based config")
	}
	
	// Fall back to file-based config
	if fileCfg, err := loadFromFile(); err == nil {
		log.Info("Loaded configuration from file")
		return fileCfg, nil
	} else {
		log.Warn("Could not load from file, using defaults")
	}
	
	// If we couldn't load config, save defaults
	if err := saveConfig(cfg); err != nil {
		log.Warnf("Failed to save default config: %v", err)
	}
	
	return cfg, nil
}

func loadFromUCI() (*Config, error) {
	// This would use a UCI library or exec uci commands
	// For now, placeholder implementation
	return nil, fmt.Errorf("UCI not implemented")
}

func loadFromFile() (*Config, error) {
	configPaths := []string{
		"/etc/ispapp/config.json",
		"./config.json",
	}
	
	for _, path := range configPaths {
		if _, err := os.Stat(path); err == nil {
			// Found a config file, load it
			file, err := os.ReadFile(path)
			if err != nil {
				continue
			}
			
			var config Config
			if err := json.Unmarshal(file, &config); err != nil {
				continue
			}
			
			return &config, nil
		}
	}
	
	return nil, fmt.Errorf("no config file found")
}

func saveConfig(cfg *Config) error {
	// Try to save to the appropriate location
	configDirs := []string{
		"/etc/ispapp",
		".",
	}
	
	for _, dir := range configDirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			continue
		}
		
		configPath := filepath.Join(dir, "config.json")
		
		data, err := json.MarshalIndent(cfg, "", "  ")
		if err != nil {
			continue
		}
		
		if err := os.WriteFile(configPath, data, 0644); err != nil {
			continue
		}
		
		return nil
	}
	
	return fmt.Errorf("could not save config to any location")
}

// Save persists the current configuration
func (c *Config) Save() error {
	return saveConfig(c)
}

// SetLogLevel changes the log level
func (c *Config) SetLogLevel(level string) {
	c.LogLevel = strings.ToLower(level)
}
