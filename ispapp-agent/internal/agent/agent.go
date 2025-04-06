package agent

import (
	"fmt"

	"ispapp-agent/internal/config"
	"ispapp-agent/internal/handlers"

	"github.com/sirupsen/logrus"
)

// Agent represents the OpenWrt agent
type Agent struct {
	config   *config.Config
	log      *logrus.Logger
	handlers []handlers.Handler
	running  bool
}

// New creates a new agent instance
func New(config *config.Config, log *logrus.Logger) (*Agent, error) {
	a := &Agent{
		config:   config,
		log:      log,
		handlers: make([]handlers.Handler, 0),
	}

	// Register handlers
	a.registerHandlers()

	return a, nil
}

func (a *Agent) registerHandlers() {
	// Register WebSocket handler first for comms
	a.handlers = append(a.handlers, handlers.NewWebSocketHandler(a.config, a.log))

	// Register other handlers
	a.handlers = append(a.handlers, handlers.NewSystemHandler(a.log))
	a.handlers = append(a.handlers, handlers.NewNetworkHandler(a.log))
}

// Start begins the agent operations
func (a *Agent) Start() error {
	a.log.Info("Starting agent...")

	// Initialize and start all handlers
	for _, h := range a.handlers {
		if err := h.Start(); err != nil {
			return fmt.Errorf("failed to start handler %s: %v", h.Name(), err)
		}
	}

	a.running = true
	a.log.Info("Agent started successfully")
	return nil
}

// Stop gracefully shuts down the agent
func (a *Agent) Stop() error {
	a.log.Info("Stopping agent...")

	for _, h := range a.handlers {
		if err := h.Stop(); err != nil {
			a.log.Errorf("Error stopping handler %s: %v", h.Name(), err)
		}
	}

	a.running = false
	return nil
}
