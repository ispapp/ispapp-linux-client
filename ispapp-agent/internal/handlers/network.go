package handlers

import (
	"os/exec"
	"strings"

	"github.com/sirupsen/logrus"
)

// NetworkHandler handles network information
type NetworkHandler struct {
	log *logrus.Logger
}

// NewNetworkHandler creates a new network handler
func NewNetworkHandler(log *logrus.Logger) *NetworkHandler {
	return &NetworkHandler{
		log: log,
	}
}

// Name returns the handler's name
func (h *NetworkHandler) Name() string {
	return "network"
}

// Start initializes the handler
func (h *NetworkHandler) Start() error {
	h.log.Info("Network handler started")
	h.collectNetworkInfo()
	return nil
}

// Stop shutdowns the handler
func (h *NetworkHandler) Stop() error {
	h.log.Info("Network handler stopped")
	return nil
}

// collectNetworkInfo gathers network information
func (h *NetworkHandler) collectNetworkInfo() {
	// Get network interfaces
	cmd := exec.Command("ip", "addr")
	output, err := cmd.CombinedOutput()
	if err != nil {
		h.log.Errorf("Failed to get network interfaces: %v", err)
		return
	}

	h.log.Debugf("Network interfaces: %s", strings.TrimSpace(string(output)))
}
