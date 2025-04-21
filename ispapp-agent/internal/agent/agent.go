package agent

import (
	"context"
	"fmt"
	"sync"
	"time"

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
	ctx      context.Context
	cancel   context.CancelFunc
	wg       sync.WaitGroup
}

// New creates a new agent instance
func New(config *config.Config, log *logrus.Logger) (*Agent, error) {
	if config == nil {
		return nil, fmt.Errorf("config cannot be nil")
	}
	if log == nil {
		return nil, fmt.Errorf("logger cannot be nil")
	}

	ctx, cancel := context.WithCancel(context.Background())
	
	a := &Agent{
		config:   config,
		log:      log,
		handlers: make([]handlers.Handler, 0),
		ctx:      ctx,
		cancel:   cancel,
	}

	// Register handlers
	a.registerHandlers()

	return a, nil
}

func (a *Agent) registerHandlers() {
	a.log.Debug("Registering handlers...")
	
	// Register WebSocket handler first for comms
	a.handlers = append(a.handlers, handlers.NewWebSocketHandler(a.config, a.log))

	// Register other handlers
	a.handlers = append(a.handlers, handlers.NewSystemHandler(a.log))
	a.handlers = append(a.handlers, handlers.NewNetworkHandler(a.log))
	
	a.log.Infof("Registered %d handlers", len(a.handlers))
}

// Start begins the agent operations
func (a *Agent) Start() error {
	a.log.Info("Starting agent...")

	if a.running {
		return fmt.Errorf("agent is already running")
	}

	// Initialize and start all handlers
	for _, h := range a.handlers {
		handlerName := h.Name()
		a.log.Debugf("Starting handler: %s", handlerName)
		
		a.wg.Add(1)
		go func(handler handlers.Handler) {
			defer a.wg.Done()
			
			if err := handler.Start(); err != nil {
				a.log.Errorf("Handler %s failed to start: %v", handler.Name(), err)
			}
			
			// Keep handler running until agent stops
			<-a.ctx.Done()
			
			if err := handler.Stop(); err != nil {
				a.log.Errorf("Error stopping handler %s: %v", handler.Name(), err)
			}
		}(h)
	}

	a.running = true
	a.log.Info("Agent started successfully")
	return nil
}

// Stop gracefully shuts down the agent
func (a *Agent) Stop() error {
	a.log.Info("Stopping agent...")

	if !a.running {
		return fmt.Errorf("agent is not running")
	}

	// Signal all handlers to stop
	a.cancel()
	
	// Wait for graceful shutdown with timeout
	done := make(chan struct{})
	go func() {
		a.wg.Wait()
		close(done)
	}()
	
	select {
	case <-done:
		a.log.Info("All handlers stopped successfully")
	case <-time.After(10 * time.Second):
		a.log.Warn("Some handlers did not stop gracefully within timeout")
	}

	a.running = false
	return nil
}

// Status returns the current status of the agent
func (a *Agent) Status() map[string]interface{} {
	status := map[string]interface{}{
		"running": a.running,
		"handlers": make([]map[string]string, 0, len(a.handlers)),
	}
	
	for _, h := range a.handlers {
		handlerStatus := map[string]string{
			"name": h.Name(),
			// We could add more handler status info if handlers exposed it
		}
		status["handlers"] = append(status["handlers"].([]map[string]string), handlerStatus)
	}
	
	return status
}
