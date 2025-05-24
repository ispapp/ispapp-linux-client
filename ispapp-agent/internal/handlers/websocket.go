package handlers

import (
	"context"
	"fmt"
	"math"
	"os/exec"
	"sync"
	"time"

	"ispapp-agent/internal/config"
	"ispapp-agent/internal/tools/constants"
	"ispapp-agent/internal/websocket"

	"github.com/sirupsen/logrus"
)

// WebSocketHandler handles WebSocket connectivity
type WebSocketHandler struct {
	BaseHandler
	log               *logrus.Logger
	client            *websocket.Client
	config            *config.Config
	ctx               context.Context
	cancel            context.CancelFunc
	wg                sync.WaitGroup
	connected         bool
	reconnectAttempts int
}

// NewWebSocketHandler creates a new WebSocket handler
func NewWebSocketHandler(config *config.Config, log *logrus.Logger) *WebSocketHandler {
	ctx, cancel := context.WithCancel(context.Background())

	return &WebSocketHandler{
		BaseHandler:       NewBaseHandler("websocket"),
		log:               log,
		config:            config,
		ctx:               ctx,
		cancel:            cancel,
		connected:         false,
		reconnectAttempts: 0,
	}
}

// RegisterDevice sends device registration data to the server
func (h *WebSocketHandler) RegisterDevice() error {
	if h.client == nil {
		return fmt.Errorf("websocket client not initialized")
	}

	h.log.Info("Registering device with server...")

	// Prepare registration data
	regData := map[string]interface{}{
		"device_id": h.config.DeviceID,
		"version":   constants.AppVersion,
		"platform":  "openwrt",
		// Add more registration data as needed
	}

	// Send registration message
	err := h.client.SendJSON("register", regData)
	if err != nil {
		return fmt.Errorf("failed to register device: %w", err)
	}

	h.log.Info("Device registered successfully")
	return nil
}

// Start initializes the handler with registration and self-healing
func (h *WebSocketHandler) Start() error {
	h.log.Info("Starting WebSocket handler...")

	// Create websocket client
	wsURL := fmt.Sprintf("wss://%s:%s/api/v1/ws",
		h.config.ServerURL,
		h.config.ListenerPort)

	h.client = websocket.NewClient(wsURL, h.config.DeviceID, h.log)

	// Setup message handlers
	h.setupMessageHandlers()

	// Start connection manager
	h.wg.Add(1)
	go h.connectionManager()

	return nil
}

func (h *WebSocketHandler) setupMessageHandlers() {
	// Register message handlers for different message types
	h.client.OnMessage("command", h.handleCommandMessage)
	h.client.OnMessage("config", h.handleConfigMessage)
	// Add more message handlers as needed
}

func (h *WebSocketHandler) handleCommandMessage(msg websocket.Message) error {
	h.log.Debugf("Received command: %v", msg.Content)

	data, ok := msg.Content.(map[string]interface{})
	if !ok {
		h.log.Error("Invalid command message format")
		return fmt.Errorf("invalid command message format")
	}

	cmd, ok := data["command"].(string)
	if !ok {
		h.log.Error("Invalid command format")
		return fmt.Errorf("invalid command format")
	}

	switch cmd {
	case "reboot":
		h.executeReboot()
	case "status":
		h.sendStatusResponse()
	// Add more command handlers
	default:
		h.log.Warnf("Unknown command: %s", cmd)
	}

	return nil
}

func (h *WebSocketHandler) handleConfigMessage(msg websocket.Message) error {
	h.log.Debugf("Received config update: %v", msg.Content)
	// Implement config update logic
	return nil
}

func (h *WebSocketHandler) executeReboot() {
	h.log.Info("Executing reboot command")
	cmd := exec.Command("reboot")
	if err := cmd.Run(); err != nil {
		h.log.Errorf("Reboot failed: %v", err)
	}
}

func (h *WebSocketHandler) sendStatusResponse() {
	// Implementation for sending status back to server
}

func (h *WebSocketHandler) connectionManager() {
	defer h.wg.Done()

	// Initial connection
	h.connect()

	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-h.ctx.Done():
			h.log.Debug("Connection manager shutting down")
			return

		case <-ticker.C:
			if !h.connected {
				h.reconnectWithBackoff()
			} else {
				// Send heartbeat/ping to keep connection alive
				if err := h.client.Ping(); err != nil {
					h.log.Warn("Connection appears to be down, reconnecting...")
					h.connected = false
				}
			}
		}
	}
}

func (h *WebSocketHandler) connect() {
	h.log.Info("Connecting to WebSocket server...")

	err := h.client.Connect()
	if err != nil {
		h.log.Errorf("Failed to connect: %v", err)
		h.connected = false
		return
	}

	// Register the device
	if err := h.RegisterDevice(); err != nil {
		h.log.Errorf("Registration failed: %v", err)
		h.client.Disconnect()
		h.connected = false
		return
	}

	h.connected = true
	h.reconnectAttempts = 0
	h.log.Info("Successfully connected to server")
}

func (h *WebSocketHandler) reconnectWithBackoff() {
	// Exponential backoff with max delay of 5 minutes
	maxDelay := 300.0
	delay := math.Min(30.0*math.Pow(1.5, float64(h.reconnectAttempts)), maxDelay)

	h.log.Infof("Reconnecting in %d seconds (attempt %d)...", int(delay), h.reconnectAttempts+1)
	time.Sleep(time.Duration(delay) * time.Second)

	h.reconnectAttempts++
	h.connect()
}

// Stop shutdowns the handler
func (h *WebSocketHandler) Stop() error {
	h.log.Info("Stopping WebSocket handler...")

	// Signal connection manager to stop
	h.cancel()

	// Disconnect WebSocket if connected
	if h.client != nil {
		h.client.Disconnect()
	}

	// Wait for all goroutines to finish
	h.wg.Wait()

	h.log.Info("WebSocket handler stopped")
	return nil
}
