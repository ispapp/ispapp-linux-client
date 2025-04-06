package handlers

import (
	"os/exec"
	"strings"

	"github.com/sirupsen/logrus"
)

// SystemHandler handles system information
type SystemHandler struct {
	log *logrus.Logger
}

// NewSystemHandler creates a new system handler
func NewSystemHandler(log *logrus.Logger) *SystemHandler {
	return &SystemHandler{
		log: log,
	}
}

// Name returns the handler's name
func (h *SystemHandler) Name() string {
	return "system"
}

// Start initializes the handler
func (h *SystemHandler) Start() error {
	h.log.Info("System handler started")
	h.collectBasicInfo()
	return nil
}

// Stop shutdowns the handler
func (h *SystemHandler) Stop() error {
	h.log.Info("System handler stopped")
	return nil
}

// collectBasicInfo gathers basic system information
func (h *SystemHandler) collectBasicInfo() {
	// Get OpenWrt version
	cmd := exec.Command("cat", "/etc/openwrt_release")
	output, err := cmd.CombinedOutput()
	if err != nil {
		h.log.Errorf("Failed to get OpenWrt version: %v", err)
		return
	}

	h.log.Infof("OpenWrt information: %s", strings.TrimSpace(string(output)))
}
