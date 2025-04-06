package events

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
)

// Handler manages event operations for the agent
type Handler struct {
	Config       ConfigManager
	Log          *logrus.Logger
	HttpClient   HttpClient
	CommandRuner CommandRunner
}

// ConfigManager interface for accessing configuration
type ConfigManager interface {
	GetString(key string) string
	GetInt(key string) int
	GetBool(key string) bool
	Set(key string, value interface{}) error
	Save() error
	UpdateTokens(accessToken, refreshToken string) error
}

// HttpClient interface for HTTP operations
type HttpClient interface {
	Get(url string) ([]byte, error)
	Post(url string, contentType string, body []byte) ([]byte, error)
	Do(req *http.Request) ([]byte, error)
}

// CommandRunner interface for executing commands
type CommandRunner interface {
	Execute(command string, args ...string) ([]byte, error)
}

// DefaultCommandRunner implements the CommandRunner interface
type DefaultCommandRunner struct{}

// Execute runs a command and returns its output
func (r *DefaultCommandRunner) Execute(command string, args ...string) ([]byte, error) {
	cmd := exec.Command(command, args...)
	return cmd.CombinedOutput()
}

// New creates a new event handler
func New(config ConfigManager, log *logrus.Logger, client HttpClient) *Handler {
	return &Handler{
		Config:       config,
		Log:          log,
		HttpClient:   client,
		CommandRuner: &DefaultCommandRunner{},
	}
}

// DeviceInfo contains basic device information
type DeviceInfo struct {
	HardwareMake         string   `json:"hardwareMake"`
	HardwareModel        string   `json:"hardwareModel"`
	HardwareModelNumber  string   `json:"hardwareModelNumber"`
	HardwareSerialNumber string   `json:"hardwareSerialNumber"`
	Hostname             string   `json:"hostname"`
	Interfaces           []string `json:"interfaces"`
	OS                   string   `json:"os"`
	OSVersion            string   `json:"osVersion"`
	// Additional fields as needed
}

// GetDeviceInfo collects device information
func (h *Handler) GetDeviceInfo() (*DeviceInfo, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return nil, err
	}

	// This is a simplified implementation
	// You would need to implement proper methods to get these values
	return &DeviceInfo{
		HardwareMake:  "OpenWrt", // Use proper detection
		HardwareModel: "Generic", // Use proper detection
		Hostname:      hostname,
		OS:            "OpenWrt",
		OSVersion:     "Unknown", // Get from /etc/os-release
	}, nil
}

// SignUp registers the device with the cloud service
func (h *Handler) SignUp() error {
	domain := h.Config.GetString("Domain")
	port := h.Config.GetInt("ListenerPort")
	key := h.Config.GetString("Key")
	login := h.Config.GetString("login")

	if domain == "" {
		return fmt.Errorf("domain not configured")
	}

	// Generate or use existing login ID
	if login == "" || login == "00000000-0000-0000-0000-000000000000" {
		// Generate a unique ID or get from firmware
		deviceInfo, err := h.GetDeviceInfo()
		if err != nil {
			return err
		}
		login = deviceInfo.HardwareSerialNumber
		if login == "" {
			login = deviceInfo.Hostname
		}
		h.Config.Set("login", login)
		h.Config.Save()
	}

	// Collect device configuration
	deviceConfig, err := h.collectConfig()
	if err != nil {
		return err
	}

	// Post to init config endpoint
	url := fmt.Sprintf("https://%s:%d/initconfig?login=%s&key=%s", domain, port, login, key)
	configJson, err := json.Marshal(deviceConfig)
	if err != nil {
		return err
	}

	h.Log.Infof("Registering device with login: %s", login)
	response, err := h.HttpClient.Post(url, "application/json", configJson)
	if err != nil {
		return err
	}

	// Process response
	var tokenResponse struct {
		AccessToken  string `json:"accessToken"`
		RefreshToken string `json:"refreshToken"`
	}

	if err := json.Unmarshal(response, &tokenResponse); err != nil {
		return fmt.Errorf("failed to parse token response: %v", err)
	}

	// Save tokens
	if tokenResponse.AccessToken != "" && tokenResponse.RefreshToken != "" {
		err = h.Config.UpdateTokens(tokenResponse.AccessToken, tokenResponse.RefreshToken)
		if err != nil {
			return fmt.Errorf("failed to save tokens: %v", err)
		}
		h.Log.Info("Device registered successfully")
	}

	return nil
}

// CheckConnection validates and refreshes the authentication token
func (h *Handler) CheckConnection() error {
	domain := h.Config.GetString("Domain")
	port := h.Config.GetInt("ListenerPort")
	accessToken := h.Config.GetString("accessToken")

	if accessToken == "" {
		return fmt.Errorf("no access token available")
	}

	url := fmt.Sprintf("https://%s:%d/auth/refresh?accessToken=%s", domain, port, accessToken)
	response, err := h.HttpClient.Get(url)
	if err != nil {
		return err
	}

	var tokenResponse struct {
		AccessToken  string `json:"accessToken"`
		RefreshToken string `json:"refreshToken"`
	}

	if err := json.Unmarshal(response, &tokenResponse); err != nil {
		return fmt.Errorf("failed to parse token response: %v", err)
	}

	// Update tokens
	if tokenResponse.AccessToken != "" && tokenResponse.RefreshToken != "" {
		err = h.Config.UpdateTokens(tokenResponse.AccessToken, tokenResponse.RefreshToken)
		if err != nil {
			return fmt.Errorf("failed to update tokens: %v", err)
		}
		h.Log.Debug("Authentication tokens refreshed")
	}

	return nil
}

// CheckForUpdates polls the server for commands or configuration updates
func (h *Handler) CheckForUpdates() error {
	domain := h.Config.GetString("Domain")
	port := h.Config.GetInt("ListenerPort")

	url := fmt.Sprintf("https://%s:%d/update", domain, port)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+h.Config.GetString("accessToken"))

	response, err := h.HttpClient.Do(req)
	if err != nil {
		return err
	}

	// Parse response for commands
	var updateResponse struct {
		ExecuteSpeedtest bool   `json:"executeSpeedtest"`
		FwStatus         string `json:"fwStatus"`
		UpdateFast       bool   `json:"updateFast"`
		Reboot           string `json:"reboot"`
	}

	if err := json.Unmarshal(response, &updateResponse); err != nil {
		return fmt.Errorf("failed to parse update response: %v", err)
	}

	// Execute commands based on response
	if updateResponse.ExecuteSpeedtest {
		go h.RunSpeedTest()
	}

	if updateResponse.FwStatus == "upgrade" {
		// Implement firmware upgrade logic
		h.Log.Info("Firmware upgrade requested")
	}

	if updateResponse.UpdateFast {
		h.Config.Set("updateInterval", 5)
		h.Config.Save()
		h.Log.Info("Update interval set to fast mode (5 seconds)")
	}

	if updateResponse.Reboot == "1" {
		h.Log.Info("Reboot requested by server, initiating reboot")
		go func() {
			time.Sleep(3 * time.Second) // Give some time for response
			h.CommandRuner.Execute("reboot")
		}()
	}

	return nil
}

// RunSpeedTest performs a network speed test
func (h *Handler) RunSpeedTest() error {
	h.Log.Info("Running speed test")

	// Determine which speed test tool to use
	var testTool string
	if _, err := exec.LookPath("iperf3"); err == nil {
		testTool = "iperf3"
	} else if _, err := exec.LookPath("iperf"); err == nil {
		testTool = "iperf"
	} else {
		return fmt.Errorf("no speed test tool available")
	}

	// Get server from config or use default
	server := h.Config.GetString("IperfServer")
	if server == "" {
		server = "iperf.longshot-router.com" // Default server
	}

	// Run the test
	startTime := time.Now()
	var output []byte
	var err error

	if testTool == "iperf3" {
		output, err = h.CommandRuner.Execute(testTool, "-c", server, "-i", "1", "-t", "1", "-P", "5", "-f", "m", "-J")
	} else {
		output, err = h.CommandRuner.Execute(testTool, "-c", server, "-i", "1", "-t", "1", "-P", "5", "-f", "m")
	}

	if err != nil {
		return fmt.Errorf("speed test failed: %v", err)
	}

	duration := time.Since(startTime)

	// Parse results
	result := map[string]interface{}{
		"up":   0.0,
		"down": 0.0,
	}

	if testTool == "iperf3" {
		var jsonResult map[string]interface{}
		if err := json.Unmarshal(output, &jsonResult); err == nil {
			if end, ok := jsonResult["end"].(map[string]interface{}); ok {
				if sumSent, ok := end["sum_sent"].(map[string]interface{}); ok {
					if bits, ok := sumSent["bits_per_second"].(float64); ok {
						result["up"] = bits / 1024 / 1024 // Convert to Mbps
					}
				}
				if sumReceived, ok := end["sum_received"].(map[string]interface{}); ok {
					if bits, ok := sumReceived["bits_per_second"].(float64); ok {
						result["down"] = bits / 1024 / 1024 // Convert to Mbps
					}
				}
			}
		}
	} else {
		// Parse iperf output
		outputStr := string(output)
		upRe := regexp.MustCompile(`SUM.*sender.*\s+(\d+\.\d+)\s+Mbits/sec`)
		downRe := regexp.MustCompile(`SUM.*receiver.*\s+(\d+\.\d+)\s+Mbits/sec`)

		upMatches := upRe.FindStringSubmatch(outputStr)
		if len(upMatches) >= 2 {
			result["up"], _ = strconv.ParseFloat(upMatches[1], 64)
		}

		downMatches := downRe.FindStringSubmatch(outputStr)
		if len(downMatches) >= 2 {
			result["down"], _ = strconv.ParseFloat(downMatches[1], 64)
		}
	}

	// Send results to server
	domain := h.Config.GetString("Domain")
	port := h.Config.GetInt("ListenerPort")

	bandwidthData := map[string]interface{}{
		"date":       time.Now().Format("2006-01-02"),
		"time":       time.Now().Format("15:04:05"),
		"txAvg":      int(result["up"].(float64) * 1024),
		"rxAvg":      int(result["down"].(float64) * 1024),
		"txDuration": fmt.Sprintf("%ds", int(duration.Seconds())),
		"rxDuration": fmt.Sprintf("%ds", int(duration.Seconds())),
		"server":     server,
	}

	bandwidthJson, err := json.Marshal(bandwidthData)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("https://%s:%d/bandwidth", domain, port)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bandwidthJson))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+h.Config.GetString("accessToken"))

	resp, err := h.HttpClient.Do(req)
	if err != nil {
		return err
	}

	// Check if server sent a new Iperf server to use
	var serverResponse struct {
		SetServer string `json:"setserver"`
	}

	if err := json.Unmarshal(resp, &serverResponse); err == nil {
		if serverResponse.SetServer != "" {
			h.Config.Set("IperfServer", serverResponse.SetServer)
			h.Config.Save()
			h.Log.Infof("Updated speed test server to: %s", serverResponse.SetServer)
		}
	}

	h.Log.Infof("Speed test completed: Up: %.2f Mbps, Down: %.2f Mbps", result["up"], result["down"])
	return nil
}

// ExecuteTerminalCommands fetches and executes commands from the server
func (h *Handler) ExecuteTerminalCommands() error {
	domain := h.Config.GetString("Domain")
	port := h.Config.GetInt("ListenerPort")

	url := fmt.Sprintf("https://%s:%d/terminal", domain, port)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte("{}")))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+h.Config.GetString("accessToken"))

	response, err := h.HttpClient.Do(req)
	if err != nil {
		return err
	}

	// Parse command queue
	var terminalResponse struct {
		Queue []struct {
			ID  string `json:"id"`
			Cmd string `json:"cmd"`
		} `json:"queue"`
	}

	if err := json.Unmarshal(response, &terminalResponse); err != nil {
		return fmt.Errorf("failed to parse terminal response: %v", err)
	}

	if len(terminalResponse.Queue) == 0 {
		h.Log.Debug("No terminal commands to execute")
		return nil
	}

	// Execute commands and collect results
	results := make([]map[string]interface{}, 0, len(terminalResponse.Queue))

	for _, cmd := range terminalResponse.Queue {
		h.Log.Infof("Executing command: %s", cmd.Cmd)

		// Execute command
		output, err := h.executeCommand(cmd.Cmd)

		// Record result
		result := map[string]interface{}{
			"id":        cmd.ID,
			"cmd":       cmd.Cmd,
			"timestamp": time.Now().Unix(),
			"stdout":    string(output),
			"stderr":    "",
			"executed":  true,
		}

		if err != nil {
			result["stderr"] = err.Error()
			h.Log.Errorf("Command failed: %v", err)
		}

		results = append(results, result)
	}

	// Send results back to server
	resultsJson, err := json.Marshal(results)
	if err != nil {
		return err
	}

	req, err = http.NewRequest("POST", url, bytes.NewBuffer(resultsJson))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+h.Config.GetString("accessToken"))

	_, err = h.HttpClient.Do(req)
	if err != nil {
		return err
	}

	h.Log.Infof("Executed %d terminal commands", len(terminalResponse.Queue))
	return nil
}

// executeCommand runs a shell command and returns its output
func (h *Handler) executeCommand(command string) ([]byte, error) {
	// Create a temporary shell script
	tmpFile := filepath.Join(os.TempDir(), "ispapp_cmd_"+strconv.FormatInt(time.Now().UnixNano(), 10))

	script := "#!/bin/sh\n" + command

	if err := ioutil.WriteFile(tmpFile, []byte(script), 0700); err != nil {
		return nil, fmt.Errorf("failed to create temporary script: %v", err)
	}
	defer os.Remove(tmpFile)

	// Execute the script
	return h.CommandRuner.Execute(tmpFile)
}

// collectConfig gathers all device configuration and status information
func (h *Handler) collectConfig() (map[string]interface{}, error) {
	// This would be a comprehensive implementation to gather device information
	// Similar to the CollectConfigs function in the Lua script

	deviceInfo, err := h.GetDeviceInfo()
	if err != nil {
		return nil, err
	}

	// Get sequence number and timestamp
	sequenceNumber := h.Config.GetInt("sequenceNumber")
	lastConfigRequest := time.Now().Unix()

	// Increment sequence number
	h.Config.Set("sequenceNumber", sequenceNumber+1)
	h.Config.Set("lastConfigRequest", lastConfigRequest)
	h.Config.Save()

	// Basic configuration data
	config := map[string]interface{}{
		"agent":                "openwrt",
		"bandwidthTestSupport": true, // Check if iperf/iperf3 exists
		"hardwareMake":         deviceInfo.HardwareMake,
		"hardwareModel":        deviceInfo.HardwareModel,
		"hostname":             deviceInfo.Hostname,
		"os":                   deviceInfo.OS,
		"osVersion":            deviceInfo.OSVersion,
		"sequenceNumber":       sequenceNumber + 1,
		"lastConfigRequest":    lastConfigRequest,
		// Additional fields would be populated here
	}

	// Add collectors data if needed (like in AddCollectorsToConfig)
	// This would include network stats, wireless data, etc.

	return config, nil
}
