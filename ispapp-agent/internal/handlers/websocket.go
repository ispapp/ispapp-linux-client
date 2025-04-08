package handlers

import (
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"ispapp-agent/internal/config"
	"ispapp-agent/internal/websocket"

	"github.com/sirupsen/logrus"
)

// WebSocketHandler handles WebSocket connectivity
type WebSocketHandler struct {
	log    *logrus.Logger
	client *websocket.Client
	config *config.Config
}

// NewWebSocketHandler creates a new WebSocket handler
func NewWebSocketHandler(config *config.Config, log *logrus.Logger) *WebSocketHandler {
	return &WebSocketHandler{
		log:    log,
		config: config,
	}
}

// Name returns the handler's name
func (h *WebSocketHandler) Name() string {
	return "websocket"
}

// RegisterDevice sends device registration data to the server
func (h *WebSocketHandler) RegisterDevice() error {
	registrationData := map[string]interface{}{
		"hostname":  h.config.DeviceID,
		"device_id": h.config.DeviceID,
		"key":       h.config.AuthToken,
	}

	if err := h.client.ReportEvent("auth", registrationData); err != nil {
		h.log.Errorf("Failed to register device: %v", err)
		return err
	}

	h.log.Info("Device registered successfully")
	return nil
}

// Start initializes the handler with registration and self-healing
func (h *WebSocketHandler) Start() error {
	h.log.Info("WebSocket handler starting")

	// Create WebSocket client
	h.client = websocket.NewClient(
		h.config.WSEndpoint,
		h.log,
		h.config.DeviceID,
		h.config.AuthToken,
	)

	// Register message handlers
	h.registerMessageHandlers()

	// Start the WebSocket client
	if err := h.client.Start(); err != nil {
		return err
	}

	// Register the device
	if err := h.RegisterDevice(); err != nil {
		return err
	}

	// Start self-healing mechanism
	h.client.SelfHeal()

	// Send initial status message
	go func() {
		time.Sleep(2 * time.Second)
		h.sendStatusMessage()
	}()

	return nil
}

// Stop shutdowns the handler
func (h *WebSocketHandler) Stop() error {
	h.log.Info("WebSocket handler stopping")
	if h.client != nil {
		// Send offline status before stopping
		h.sendOfflineStatus()
		return h.client.Stop()
	}
	return nil
}

// Register handlers for different message types
func (h *WebSocketHandler) registerMessageHandlers() {
	// Handle ping messages
	h.client.RegisterHandler("ping", func(message websocket.Message) error {
		h.log.Debug("Received ping, sending pong")
		return h.client.Send("pong", nil)
	})

	// Handle command messages
	h.client.RegisterHandler("command", func(message websocket.Message) error {
		h.log.Info("Received command from server")
		// Process command
		// ...
		return nil
	})

	// Add more handlers as needed
}

// Send a status message with basic device information
func (h *WebSocketHandler) sendStatusMessage() {
	hostname, _ := os.Hostname()

	statusInfo := map[string]interface{}{
		"device_id": h.config.DeviceID,
		"hostname":  hostname,
		"status":    "online",
		"version":   "1.0.0", // This should be fetched from app version
		"uptime":    getSystemUptime(),
	}

	if err := h.client.Send("status", statusInfo); err != nil {
		h.log.Errorf("Failed to send status message: %v", err)
	}
}

// Send offline status before shutdown
func (h *WebSocketHandler) sendOfflineStatus() {
	statusInfo := map[string]interface{}{
		"device_id": h.config.DeviceID,
		"status":    "offline",
	}

	// Try to send but don't wait too long
	sendDone := make(chan bool)
	go func() {
		err := h.client.Send("status", statusInfo)
		if err != nil {
			h.log.Debugf("Failed to send offline status: %v", err)
		}
		sendDone <- true
	}()

	// Give it a second to send before continuing
	select {
	case <-sendDone:
		// Message sent
	case <-time.After(time.Second):
		// Timeout, continue with shutdown
		h.log.Debug("Timed out waiting to send offline status")
	}
}

// getSystemUptime returns the system uptime in seconds
func getSystemUptime() int64 {
	cmd := exec.Command("cat", "/proc/uptime")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return 0
	}

	// Output is like "12345.67 12345.68" - we only want the first number
	parts := strings.Split(string(output), " ")
	if len(parts) < 1 {
		return 0
	}

	// Convert to seconds (integer)
	uptimeFloat, err := strconv.ParseFloat(parts[0], 64)
	if err != nil {
		return 0
	}

	return int64(uptimeFloat)
}
