package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"ispapp-agent/internal/agent"
	"ispapp-agent/internal/config"
	"ispapp-agent/internal/tools/constants"

	"github.com/sirupsen/logrus"
)

func main() {
	// Setup logger
	log := setupLogger()

	// Load configuration
	cfg, err := loadConfig(log)
	if err != nil {
		log.Fatalf("Fatal error: %v", err)
	}

	// Store config in constants for legacy code that needs it
	constants.Cfg = cfg

	// Initialize agent
	a, err := agent.New(cfg, log)
	if err != nil {
		log.Fatalf("Failed to initialize agent: %v", err)
	}

	// Start agent
	if err := a.Start(); err != nil {
		log.Fatalf("Failed to start agent: %v", err)
	}

	// Set up signal handling for graceful shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	// Block until we receive a signal
	sig := <-sigCh
	log.Infof("Received signal %v, shutting down...", sig)

	// Stop the agent with a timeout
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	shutdownCh := make(chan struct{})
	go func() {
		if err := a.Stop(); err != nil {
			log.Errorf("Error during shutdown: %v", err)
		}
		close(shutdownCh)
	}()

	select {
	case <-shutdownCtx.Done():
		log.Warn("Shutdown timed out")
	case <-shutdownCh:
		log.Info("Shutdown completed successfully")
	}
}

func setupLogger() *logrus.Logger {
	log := logrus.New()
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
		DisableColors: false,
	})

	// Default to info level, will update from config later
	log.SetLevel(logrus.InfoLevel)

	return log
}

func loadConfig(log *logrus.Logger) (*config.Config, error) {
	log.Info("Loading configuration...")

	cfg, err := config.Load(log)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	// Update log level from config
	if cfg != nil {
		if logLevel, err := logrus.ParseLevel(cfg.LogLevel); err == nil {
			log.SetLevel(logLevel)
			log.Infof("Log level set to: %s", logLevel)
		}

		log.Infof("Using device ID: %s", cfg.DeviceID)
		log.Infof("API endpoint: %s", cfg.ServerURL)
	}

	return cfg, nil
}
