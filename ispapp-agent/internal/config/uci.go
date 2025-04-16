package config

import (
	"fmt"
	"ispapp-agent/internal/config/muci"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"

	"github.com/sirupsen/logrus"
)

const (
	// UCIConfigPackage is the UCI package name for our config
	UCIConfigPackage = "ispapp"

	// UCIConfigSection is the main UCI config section
	UCIConfigSettingsSection = "settings"

	// UCIConfigOverviewSection is the overview section
	UCIConfigOverviewSection = "overview"
)

// UCIConfig represents the UCI configuration
type UCIConfig struct {
	// Settings section
	Enabled        bool     `uci:"enabled,bool"`
	Login          string   `uci:"login"`
	Domain         string   `uci:"Domain"`
	ListenerPort   int      `uci:"ListenerPort,int"`
	Connected      bool     `uci:"connected,bool"`
	Key            string   `uci:"Key"`
	RefreshToken   string   `uci:"refreshToken"`
	AccessToken    string   `uci:"accessToken"`
	UpdateInterval int      `uci:"updateInterval,int"`
	IperfServer    string   `uci:"IperfServer"`
	PingTargets    []string `uci:"pingTargets,list"`

	// Overview section
	OverviewEnabled bool `uci:"enabled,bool,overview"`
	HighRefresh     bool `uci:"highrefresh,bool,overview"`
	RefreshInterval int  `uci:"refreshInterval,int,overview"`

	log     *logrus.Logger
	mutex   sync.RWMutex
	uciRoot string
}

// NewUCIConfig creates a new UCI configuration controller
func NewUCIConfig(log *logrus.Logger) *UCIConfig {
	return &UCIConfig{
		log:     log,
		uciRoot: "/etc/config",
	}
}

// Load loads the UCI configuration
func (c *UCIConfig) Load() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.log.Info("Loading UCI configuration")

	// Check if the UCI config exists, create it if it doesn't
	if err := c.ensureConfigExists(); err != nil {
		return fmt.Errorf("failed to ensure config exists: %v", err)
	}

	// Load UCI configuration
	uciInstance := muci.NewUCI()
	configFile := fmt.Sprintf("%s/%s", c.uciRoot, UCIConfigPackage)
	config, err := muci.LoadConfig(configFile)
	if err != nil {
		return fmt.Errorf("failed to load UCI config: %v", err)
	}
	uciInstance.Configs[UCIConfigPackage] = config

	// Load settings section
	if err := c.loadSettings(uciInstance); err != nil {
		return fmt.Errorf("failed to load settings: %v", err)
	}
	c.log.Infof("Loaded settings: %+v", c)

	// Load overview section
	if err := c.loadOverview(uciInstance); err != nil {
		return fmt.Errorf("failed to load overview: %v", err)
	}

	c.log.Infof("UCI configuration loaded successfully : %+v", c)
	return nil
}

// Save saves the UCI configuration
func (c *UCIConfig) Save() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.log.Info("Saving UCI configuration")

	// Create UCI instance
	uciInstance := muci.NewUCI()
	configFile := fmt.Sprintf("%s/%s", c.uciRoot, UCIConfigPackage)
	config, err := muci.LoadConfig(configFile)
	if err != nil {
		return fmt.Errorf("failed to load UCI config for saving: %v", err)
	}
	uciInstance.Configs[UCIConfigPackage] = config

	// Save settings section
	if err := c.saveSettings(uciInstance); err != nil {
		return fmt.Errorf("failed to save settings: %v", err)
	}

	// Save overview section
	if err := c.saveOverview(uciInstance); err != nil {
		return fmt.Errorf("failed to save overview: %v", err)
	}

	// Save the configuration to disk
	if err := config.Save(); err != nil {
		return fmt.Errorf("failed to save UCI config: %v", err)
	}

	c.log.Info("UCI configuration saved successfully")
	return nil
}

// ToStandardConfig converts the UCI configuration to the standard configuration
func (c *UCIConfig) ToStandardConfig() *Config {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	wsEndpoint := fmt.Sprintf("wss://%s:%d/agent/ws", c.Domain, c.ListenerPort)
	serverURL := fmt.Sprintf("https://%s:%d", c.Domain, c.ListenerPort)

	logLevel := "info"
	if c.HighRefresh {
		logLevel = "debug"
	}

	return &Config{
		LogLevel:      logLevel,
		ServerURL:     serverURL,
		WSEndpoint:    wsEndpoint,
		DeviceID:      c.Login,
		AuthToken:     c.AccessToken,
		CheckInterval: c.UpdateInterval * 60, // Convert to seconds
		Debug:         c.HighRefresh,
	}
}

// SetConnected sets the connected status
func (c *UCIConfig) SetConnected(connected bool) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.Connected = connected

	// Create UCI instance using our custom muci package
	uciInstance := muci.NewUCI()
	configFile := fmt.Sprintf("%s/%s", c.uciRoot, UCIConfigPackage)
	config, err := muci.LoadConfig(configFile)
	if err != nil {
		return fmt.Errorf("failed to load UCI config for setting connected status: %v", err)
	}
	uciInstance.Configs[UCIConfigPackage] = config

	// Set value
	connValue := "0"
	if connected {
		connValue = "1"
	}

	if err := config.Set(UCIConfigSettingsSection, "settings", "connected", connValue, false); err != nil {
		return fmt.Errorf("failed to set connected status: %v", err)
	}

	// Save the configuration to disk
	if err := config.Save(); err != nil {
		return fmt.Errorf("failed to save UCI config: %v", err)
	}

	return nil
}

// UpdateTokens updates the auth tokens
func (c *UCIConfig) UpdateTokens(accessToken, refreshToken string) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.AccessToken = accessToken
	c.RefreshToken = refreshToken

	// Create UCI instance
	uciInstance := muci.NewUCI()
	configFile := fmt.Sprintf("%s/%s", c.uciRoot, UCIConfigPackage)
	config, err := muci.LoadConfig(configFile)
	if err != nil {
		return fmt.Errorf("failed to load UCI config for token update: %v", err)
	}
	uciInstance.Configs[UCIConfigPackage] = config

	if err := config.Set(UCIConfigSettingsSection, "settings", "accessToken", accessToken, false); err != nil {
		return fmt.Errorf("failed to set access token: %v", err)
	}

	if err := config.Set(UCIConfigSettingsSection, "settings", "refreshToken", refreshToken, false); err != nil {
		return fmt.Errorf("failed to set refresh token: %v", err)
	}

	// Save the configuration to disk
	if err := config.Save(); err != nil {
		return fmt.Errorf("failed to save UCI config: %v", err)
	}

	return nil
}

// ensureConfigExists ensures the UCI config exists
func (c *UCIConfig) ensureConfigExists() error {
	// Check if config file exists
	configFile := fmt.Sprintf("%s/%s", c.uciRoot, UCIConfigPackage)
	_, err := os.Stat(configFile)

	// Create config dir if it doesn't exist
	os.MkdirAll(c.uciRoot, 0755)

	// If file doesn't exist, create default config
	if os.IsNotExist(err) {
		c.log.Info("Creating default UCI configuration file")
		return c.createDefaultConfig()
	}

	return nil
}

// createDefaultConfig creates a default UCI configuration
func (c *UCIConfig) createDefaultConfig() error {
	// Create UCI instance
	configFile := fmt.Sprintf("%s/%s", c.uciRoot, UCIConfigPackage)
	config, err := muci.LoadConfig(configFile)
	if err != nil {
		return fmt.Errorf("failed to load UCI config: %v", err)
	}
	// defer ctx.Revert()
	// Create settings section
	if err := config.AddSection("settings", "settings"); err != nil {
		return fmt.Errorf("failed to create settings section: %v", err)
	}

	// Set default settings
	defaults := map[string]string{
		"enabled":        "1",
		"login":          "00000000-000-0000-0000-000000000000",
		"Domain":         "prv.cloud.ispapp.co",
		"ListenerPort":   "443",
		"connected":      "0",
		"Key":            "",
		"refreshToken":   "",
		"accessToken":    "",
		"updateInterval": "1",
		"IperfServer":    "",
	}

	for key, value := range defaults {
		if err := config.Set(UCIConfigSettingsSection, UCIConfigSettingsSection, key, value, len(strings.Split(value, " ")) > 0); err != nil {
			return fmt.Errorf("failed to set default setting %s", key)
		}
	}

	// Add default ping targets
	defaultTargets := []string{
		"cloud.ispapp.co",
		"aws-eu-west-2-ping.ispapp.co",
		"aws-sa-east-1-ping.ispapp.co",
		"aws-us-east-1-ping.ispapp.co",
		"aws-us-west-1-ping.ispapp.co",
	}

	for _, target := range defaultTargets {
		if err := config.Set(UCIConfigSettingsSection, "settings", "pingTargets", target, true); err != nil {
			return fmt.Errorf("failed to add default ping target: %v", err)
		}
	}

	// Create overview section
	if err := config.AddSection(UCIConfigOverviewSection, UCIConfigOverviewSection); err != nil {
		return fmt.Errorf("failed to create overview section: %v", err)
	}

	// Set default overview settings
	overviewDefaults := map[string]string{
		"enabled":         "1",
		"highrefresh":     "0",
		"refreshInterval": "0",
	}

	for key, value := range overviewDefaults {
		if err := config.Set(UCIConfigOverviewSection, UCIConfigOverviewSection, key, value, false); err != nil {
			return fmt.Errorf("failed to set default overview setting %s: %v", key, err)
		}
	}

	// Commit changes
	if err := config.Save(); err != nil {
		return fmt.Errorf("failed to commit UCI config: %v", err)
	}

	return nil
}

// loadSettings loads the settings section
func (c *UCIConfig) loadSettings(uci *muci.UCI) error {
	config := uci.Configs[UCIConfigPackage]

	// Get "enabled"
	enabledValue, err := config.Get(UCIConfigSettingsSection, "settings", "enabled")
	if err != nil {
		c.log.Errorf("failed to get 'enabled' value: %v", err)
		c.Enabled = false // Default to false if not found
	} else {
		if strValue, ok := enabledValue.(string); ok {
			c.Enabled = strValue == "1"
		}
	}

	// Get "login" value
	loginValue, err := config.Get(UCIConfigSettingsSection, "settings", "login")
	if err == nil {
		if strValue, ok := loginValue.(string); ok {
			c.Login = strValue
			c.log.Infof("Loaded Login: %s", c.Login)
		}
	} else {
		c.Login = "00000000-0000-0000-0000-000000000000" // Default value
	}
	if c.Login == "" || c.Login == "00000000-0000-0000-0000-000000000000" {
		c.log.Warn("Login not found in UCI configuration, using default value")
		if err := c.getFromFlash("login"); err != nil {
			c.log.Errorf("Failed to get login from flash memory: %v", err)
		} else {
			c.log.Infof("Loaded Login from flash: %s", c.Login)
		}
	}

	// Get "Domain" value
	domainValue, err := config.Get(UCIConfigSettingsSection, "settings", "Domain")
	if err == nil {
		if strValue, ok := domainValue.(string); ok {
			c.Domain = strValue
			c.log.Infof("Loaded Domain: %s", c.Domain)
		}
	} else {
		c.log.Warn("Domain not found in UCI configuration")
	}

	// Get "connected" value
	connectedValue, err := config.Get(UCIConfigSettingsSection, "settings", "connected")
	if err == nil {
		if strValue, ok := connectedValue.(string); ok {
			c.Connected = strValue == "1"
		}
	} else {
		c.Connected = false
	}

	// Get "Key" value
	keyValue, err := config.Get(UCIConfigSettingsSection, "settings", "Key")
	if err == nil {
		if strValue, ok := keyValue.(string); ok {
			c.Key = strValue
		}
	}

	// Get "refreshToken" value
	refreshTokenValue, err := config.Get(UCIConfigSettingsSection, "settings", "refreshToken")
	if err == nil {
		if strValue, ok := refreshTokenValue.(string); ok {
			c.RefreshToken = strValue
		}
	}

	// Get "accessToken" value
	accessTokenValue, err := config.Get(UCIConfigSettingsSection, "settings", "accessToken")
	if err == nil {
		if strValue, ok := accessTokenValue.(string); ok {
			c.AccessToken = strValue
		}
	}

	// Get "IperfServer" value
	iperfServerValue, err := config.Get(UCIConfigSettingsSection, "settings", "IperfServer")
	if err == nil {
		if strValue, ok := iperfServerValue.(string); ok {
			c.IperfServer = strValue
		}
	}

	// Get and convert "ListenerPort" value
	listenerPortValue, err := config.Get(UCIConfigSettingsSection, "settings", "ListenerPort")
	if err == nil {
		if strValue, ok := listenerPortValue.(string); ok {
			if port, err := strconv.Atoi(strValue); err == nil {
				c.ListenerPort = port
			} else {
				c.ListenerPort = 443 // Default if parsing fails
			}
		}
	} else {
		c.ListenerPort = 443 // Default if missing
	}

	// Get and convert "updateInterval" value
	updateIntervalValue, err := config.Get(UCIConfigSettingsSection, "settings", "updateInterval")
	if err == nil {
		if strValue, ok := updateIntervalValue.(string); ok {
			if interval, err := strconv.Atoi(strValue); err == nil {
				c.UpdateInterval = interval
			} else {
				c.UpdateInterval = 1 // Default if parsing fails
			}
		}
	} else {
		c.UpdateInterval = 1 // Default if missing
	}

	// Get "pingTargets" list
	pingTargetsValue, err := config.Get(UCIConfigSettingsSection, "settings", "pingTargets")
	if err == nil {
		if listValue, ok := pingTargetsValue.([]string); ok {
			c.PingTargets = listValue
		}
	}

	return nil
}

// loadOverview loads the overview section
func (c *UCIConfig) loadOverview(uci *muci.UCI) error {
	config := uci.Configs[UCIConfigPackage]

	// Get "enabled" value
	enabledValue, err := config.Get(UCIConfigOverviewSection, "overview", "enabled")
	if err != nil {
		return fmt.Errorf("failed to get 'enabled' value in overview section: %v", err)
	}
	if strValue, ok := enabledValue.(string); ok {
		c.OverviewEnabled = strValue == "1"
	}

	// Get "highrefresh" value
	highRefreshValue, err := config.Get(UCIConfigOverviewSection, "overview", "highrefresh")
	if err == nil {
		if strValue, ok := highRefreshValue.(string); ok {
			c.HighRefresh = strValue == "1"
		}
	}

	// Get and convert "refreshInterval" value
	refreshIntervalValue, err := config.Get(UCIConfigOverviewSection, "overview", "refreshInterval")
	if err == nil {
		if strValue, ok := refreshIntervalValue.(string); ok {
			if interval, err := strconv.Atoi(strValue); err == nil {
				c.RefreshInterval = interval
			} else {
				c.RefreshInterval = 0 // Default if parsing fails
			}
		}
	} else {
		c.RefreshInterval = 0 // Default if missing
	}

	return nil
}

// saveSettings saves the settings section
func (c *UCIConfig) saveSettings(uci *muci.UCI) error {
	config := uci.Configs[UCIConfigPackage]

	// Save all settings
	settingsMap := map[string]string{
		"enabled":        boolToUCI(c.Enabled),
		"login":          c.Login,
		"Domain":         c.Domain,
		"ListenerPort":   strconv.Itoa(c.ListenerPort),
		"connected":      boolToUCI(c.Connected),
		"Key":            c.Key,
		"refreshToken":   c.RefreshToken,
		"accessToken":    c.AccessToken,
		"updateInterval": strconv.Itoa(c.UpdateInterval),
		"IperfServer":    c.IperfServer,
	}

	for key, value := range settingsMap {
		if err := config.Set(UCIConfigSettingsSection, "settings", key, value, false); err != nil {
			return fmt.Errorf("failed to set %s: %v", key, err)
		}
	}

	// Handle ping targets as a list
	for _, target := range c.PingTargets {
		if err := config.Set(UCIConfigSettingsSection, "settings", "pingTargets", target, true); err != nil {
			return fmt.Errorf("failed to set pingTargets: %v", err)
		}
	}

	// Persist login to flash memory
	if c.Login != "" && c.Login != "00000000-0000-0000-0000-000000000000" {
		if err := persistToFlash("login", c.Login); err != nil {
			return fmt.Errorf("failed to persist login to flash memory: %v", err)
		}
	}

	// Persist values to flash memory only when they are valid
	if c.Key != "" {
		if err := persistToFlash("Key", c.Key); err != nil {
			return fmt.Errorf("failed to persist Key to flash memory: %v", err)
		}
	}

	if c.RefreshToken != "" {
		if err := persistToFlash("refreshToken", c.RefreshToken); err != nil {
			return fmt.Errorf("failed to persist refreshToken to flash memory: %v", err)
		}
	}

	if c.AccessToken != "" {
		if err := persistToFlash("accessToken", c.AccessToken); err != nil {
			return fmt.Errorf("failed to persist accessToken to flash memory: %v", err)
		}
	}

	if c.IperfServer != "" {
		if err := persistToFlash("IperfServer", c.IperfServer); err != nil {
			return fmt.Errorf("failed to persist IperfServer to flash memory: %v", err)
		}
	}

	if c.Domain != "" {
		if err := persistToFlash("Domain", c.Domain); err != nil {
			return fmt.Errorf("failed to persist Domain to flash memory: %v", err)
		}
	}

	return nil
}

// saveOverview saves the overview section
func (c *UCIConfig) saveOverview(uci *muci.UCI) error {
	config := uci.Configs[UCIConfigPackage]

	// Save all settings
	overviewMap := map[string]string{
		"enabled":         boolToUCI(c.OverviewEnabled),
		"highrefresh":     boolToUCI(c.HighRefresh),
		"refreshInterval": strconv.Itoa(c.RefreshInterval),
	}

	for key, value := range overviewMap {
		if err := config.Set(UCIConfigOverviewSection, "overview", key, value, false); err != nil {
			return fmt.Errorf("failed to set %s: %v", key, err)
		}
	}

	return nil
}

func (c *UCIConfig) getFromFlash(key string) error {
	value, err := getFromFlash(key)
	// find the exact field and set it
	if key == "login" {
		c.Login = value
	} else if key == "Key" {
		c.Key = value
	} else if key == "refreshToken" {
		c.RefreshToken = value
	} else if key == "accessToken" {
		c.AccessToken = value
	} else if key == "IperfServer" {
		c.IperfServer = value
	} else if key == "Domain" {
		c.Domain = value
	} else if key == "ListenerPort" {
		if port, err := strconv.Atoi(value); err == nil {
			c.ListenerPort = port
		}
	}
	if err != nil {
		return err
	}

	return nil
}

func getFromFlash(key string) (string, error) {
	cmd := exec.Command("fw_printenv", key)
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("fw_printenv failed: %v", err)
	}
	return strings.TrimSpace(string(output)), nil
}

// persistToFlash persists a key-value pair to flash memory using fw_setenv
func persistToFlash(key, value string) error {
	cmd := exec.Command("fw_setenv", key, value)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("fw_setenv failed: %v", err)
	}

	// Verify the value with fw_printenv
	verifyCmd := exec.Command("fw_printenv", key)
	output, err := verifyCmd.Output()
	if err != nil {
		return fmt.Errorf("fw_printenv failed: %v", err)
	}

	if !strings.Contains(string(output), fmt.Sprintf("%s=%s", key, value)) {
		return fmt.Errorf("verification failed: expected %s=%s, got %s", key, value, string(output))
	}

	return nil
}

// boolToUCI converts a bool to UCI format (0/1)
func boolToUCI(b bool) string {
	if b {
		return "1"
	}
	return "0"
}
