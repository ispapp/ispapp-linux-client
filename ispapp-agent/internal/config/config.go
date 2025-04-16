package config

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"

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
		return defaultConfig(), nil // no realy usefull
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
		// Check firmware environment for DeviceID
		deviceID, err := getDeviceIDFromFirmware()
		if err != nil || deviceID == "" {
			log.Warn("DeviceID not found in firmware, attempting to fetch from API")
			deviceID, err = fetchDeviceIDFromAPI("https://prv.cloud.ispapp.co/auth/uuid")
			if err != nil || deviceID == "" {
				log.Warn("Failed to fetch DeviceID from API, generating default UUID")
				deviceID = "00000000-0000-0000-0000-000000000000"
			}
			// Persist the DeviceID
			if err := persistDeviceIDToFirmware(deviceID); err != nil {
				log.Errorf("Failed to persist DeviceID: %v", err)
			}
		}
		config.DeviceID = deviceID
		log.Infof("Using DeviceID: %s", deviceID)
	}

	if config.CheckInterval <= 0 {
		config.CheckInterval = 60
		log.Info("Invalid check interval, using default of 60 seconds")
	}

	return config, nil
}

// getDeviceIDFromFirmware retrieves the DeviceID from the firmware environment
func getDeviceIDFromFirmware() (string, error) {
	cmd := exec.Command("fw_printenv", "login")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	// Parse the output to extract the value
	parts := strings.SplitN(string(output), "=", 2)
	if len(parts) != 2 {
		return "", fmt.Errorf("unexpected output from fw_printenv: %s", output)
	}

	return strings.TrimSpace(parts[1]), nil
}

// fetchDeviceIDFromAPI fetches the DeviceID from the online API
func fetchDeviceIDFromAPI(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API returned non-200 status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(body)), nil
}

// persistDeviceIDToFirmware persists the DeviceID to the firmware environment
func persistDeviceIDToFirmware(deviceID string) error {
	cmd := exec.Command("fw_setenv", "login", deviceID)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to set device_id: %v", err)
	}
	return nil
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
