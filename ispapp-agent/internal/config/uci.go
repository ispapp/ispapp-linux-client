package config

import (
	"fmt"
	uci "ispapp-agent/internal/uci"
	"os"
	"strconv"
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
	ctx := uci.NewTree(c.uciRoot)
	defer ctx.Revert()

	// Load settings section
	if err := c.loadSettings(ctx); err != nil {
		return fmt.Errorf("failed to load settings: %v", err)
	}

	// Load overview section
	if err := c.loadOverview(ctx); err != nil {
		return fmt.Errorf("failed to load overview: %v", err)
	}

	c.log.Info("UCI configuration loaded successfully")
	return nil
}

// Save saves the UCI configuration
func (c *UCIConfig) Save() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.log.Info("Saving UCI configuration")

	// Create UCI instance
	ctx := uci.NewTree(c.uciRoot)
	defer ctx.Revert()

	// Save settings section
	if err := c.saveSettings(ctx); err != nil {
		return fmt.Errorf("failed to save settings: %v", err)
	}

	// Save overview section
	if err := c.saveOverview(ctx); err != nil {
		return fmt.Errorf("failed to save overview: %v", err)
	}

	// Commit changes
	if err := ctx.Commit(); err != nil {
		return fmt.Errorf("failed to commit UCI config: %v", err)
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

	// Create UCI instance
	ctx := uci.NewTree(c.uciRoot)
	defer ctx.Revert()

	// Set value
	connValue := "0"
	if connected {
		connValue = "1"
	}

	if ok := ctx.Set(UCIConfigPackage, UCIConfigSettingsSection, "connected", connValue); !ok {
		return fmt.Errorf("failed to set connected status")
	}

	// Commit changes
	if err := ctx.Commit(); err != nil {
		return fmt.Errorf("failed to commit UCI config: %v", err)
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
	ctx := uci.NewTree(c.uciRoot)
	defer ctx.Revert()

	// Set tokens
	if ok := ctx.Set(UCIConfigPackage, UCIConfigSettingsSection, "accessToken", accessToken); !ok {
		return fmt.Errorf("failed to set access token")
	}

	if ok := ctx.Set(UCIConfigPackage, UCIConfigSettingsSection, "refreshToken", refreshToken); !ok {
		return fmt.Errorf("failed to set refresh token")
	}

	// Commit changes
	if err := ctx.Commit(); err != nil {
		return fmt.Errorf("failed to commit UCI config: %v", err)
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
	ctx := uci.NewTree(c.uciRoot)
	defer ctx.Revert()

	// Create settings section
	if err := ctx.AddSection(UCIConfigPackage, UCIConfigSettingsSection, "settings"); err != nil {
		return fmt.Errorf("failed to create settings section: %v", err)
	}

	// Set default settings
	defaults := map[string]string{
		"enabled":        "1",
		"login":          "",
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
		if ok := ctx.Set(UCIConfigPackage, UCIConfigSettingsSection, key, value); !ok {
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
		if ok := ctx.Set(UCIConfigPackage, UCIConfigSettingsSection, "pingTargets", target); !ok {
			return fmt.Errorf("failed to add default ping target")
		}
	}

	// Create overview section
	if err := ctx.AddSection(UCIConfigPackage, UCIConfigOverviewSection, "overview"); err != nil {
		return fmt.Errorf("failed to create overview section: %v", err)
	}

	// Set default overview settings
	overviewDefaults := map[string]string{
		"enabled":         "1",
		"highrefresh":     "0",
		"refreshInterval": "0",
	}

	for key, value := range overviewDefaults {
		if ok := ctx.Set(UCIConfigPackage, UCIConfigOverviewSection, key, value); !ok {
			return fmt.Errorf("failed to set default overview setting %s", key)
		}
	}

	// Commit changes
	if err := ctx.Commit(); err != nil {
		return fmt.Errorf("failed to commit UCI config: %v", err)
	}

	return nil
}

// loadSettings loads the settings section
func (c *UCIConfig) loadSettings(ctx uci.Tree) error {
	// Get "enabled" value
	enabledValues, ok := ctx.Get(UCIConfigPackage, UCIConfigSettingsSection, "enabled")
	if !ok || len(enabledValues) == 0 {
		return fmt.Errorf("failed to get 'enabled' value got:%v", enabledValues)
	}
	c.Enabled = enabledValues[0] == "1"

	// Get "login" value
	loginValues, ok := ctx.Get(UCIConfigPackage, UCIConfigSettingsSection, "login")
	if ok && len(loginValues) > 0 {
		c.Login = loginValues[0]
	}

	// Get "Domain" value
	domainValues, ok := ctx.Get(UCIConfigPackage, UCIConfigSettingsSection, "Domain")
	if ok && len(domainValues) > 0 {
		c.Domain = domainValues[0]
	}

	// Get "connected" value
	connectedValues, ok := ctx.Get(UCIConfigPackage, UCIConfigSettingsSection, "connected")
	if ok && len(connectedValues) > 0 {
		c.Connected = connectedValues[0] == "1"
	}

	// Get "Key" value
	keyValues, ok := ctx.Get(UCIConfigPackage, UCIConfigSettingsSection, "Key")
	if ok && len(keyValues) > 0 {
		c.Key = keyValues[0]
	}

	// Get "refreshToken" value
	refreshTokenValues, ok := ctx.Get(UCIConfigPackage, UCIConfigSettingsSection, "refreshToken")
	if ok && len(refreshTokenValues) > 0 {
		c.RefreshToken = refreshTokenValues[0]
	}

	// Get "accessToken" value
	accessTokenValues, ok := ctx.Get(UCIConfigPackage, UCIConfigSettingsSection, "accessToken")
	if ok && len(accessTokenValues) > 0 {
		c.AccessToken = accessTokenValues[0]
	}

	// Get "IperfServer" value
	iperfServerValues, ok := ctx.Get(UCIConfigPackage, UCIConfigSettingsSection, "IperfServer")
	if ok && len(iperfServerValues) > 0 {
		c.IperfServer = iperfServerValues[0]
	}

	// Get and convert "ListenerPort" value
	listenerPortValues, ok := ctx.Get(UCIConfigPackage, UCIConfigSettingsSection, "ListenerPort")
	if ok && len(listenerPortValues) > 0 {
		if port, err := strconv.Atoi(listenerPortValues[0]); err == nil {
			c.ListenerPort = port
		} else {
			c.ListenerPort = 443 // Default if parsing fails
		}
	} else {
		c.ListenerPort = 443 // Default if missing
	}

	// Get and convert "updateInterval" value
	updateIntervalValues, ok := ctx.Get(UCIConfigPackage, UCIConfigSettingsSection, "updateInterval")
	if ok && len(updateIntervalValues) > 0 {
		if interval, err := strconv.Atoi(updateIntervalValues[0]); err == nil {
			c.UpdateInterval = interval
		} else {
			c.UpdateInterval = 1 // Default if parsing fails
		}
	} else {
		c.UpdateInterval = 1 // Default if missing
	}

	// Get "pingTargets" list
	pingTargetsValues, ok := ctx.Get(UCIConfigPackage, UCIConfigSettingsSection, "pingTargets")
	if ok {
		c.PingTargets = pingTargetsValues
	}

	return nil
}

// loadOverview loads the overview section
func (c *UCIConfig) loadOverview(ctx uci.Tree) error {
	// Get "enabled" value
	enabledValues, ok := ctx.Get(UCIConfigPackage, UCIConfigOverviewSection, "enabled")
	if !ok || len(enabledValues) == 0 {
		return fmt.Errorf("failed to get 'enabled' value in overview section")
	}
	c.OverviewEnabled = enabledValues[0] == "1"

	// Get "highrefresh" value
	highRefreshValues, ok := ctx.Get(UCIConfigPackage, UCIConfigOverviewSection, "highrefresh")
	if ok && len(highRefreshValues) > 0 {
		c.HighRefresh = highRefreshValues[0] == "1"
	}

	// Get and convert "refreshInterval" value
	refreshIntervalValues, ok := ctx.Get(UCIConfigPackage, UCIConfigOverviewSection, "refreshInterval")
	if ok && len(refreshIntervalValues) > 0 {
		if interval, err := strconv.Atoi(refreshIntervalValues[0]); err == nil {
			c.RefreshInterval = interval
		} else {
			c.RefreshInterval = 0 // Default if parsing fails
		}
	} else {
		c.RefreshInterval = 0 // Default if missing
	}

	return nil
}

// saveSettings saves the settings section
func (c *UCIConfig) saveSettings(ctx uci.Tree) error {
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
		if ok := ctx.Set(UCIConfigPackage, UCIConfigSettingsSection, key, value); !ok {
			return fmt.Errorf("failed to set %s", key)
		}
	}

	// Handle pingTargets list - first delete existing
	ctx.Del(UCIConfigPackage, UCIConfigSettingsSection, "pingTargets")

	// Add each ping target
	for _, target := range c.PingTargets {
		if ok := ctx.Set(UCIConfigPackage, UCIConfigSettingsSection, "pingTargets", target); !ok {
			return fmt.Errorf("failed to add ping target: %v", target)
		}
	}

	return nil
}

// saveOverview saves the overview section
func (c *UCIConfig) saveOverview(ctx uci.Tree) error {
	// Save all settings
	overviewMap := map[string]string{
		"enabled":         boolToUCI(c.OverviewEnabled),
		"highrefresh":     boolToUCI(c.HighRefresh),
		"refreshInterval": strconv.Itoa(c.RefreshInterval),
	}

	for key, value := range overviewMap {
		if ok := ctx.Set(UCIConfigPackage, UCIConfigOverviewSection, key, value); !ok {
			return fmt.Errorf("failed to set %s", key)
		}
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
