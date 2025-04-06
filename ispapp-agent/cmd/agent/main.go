package main

import (
	"os"
	"os/signal"
	"syscall"

	"ispapp-agent/internal/agent"
	"ispapp-agent/internal/config"

	"github.com/sirupsen/logrus"
)

func main() {
	log := logrus.New()
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	// Set reasonable default log level
	log.SetLevel(logrus.InfoLevel)

	// Load configuration
	cfg, err := config.Load(log)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Update log level from config
	logLevel, err := logrus.ParseLevel(cfg.LogLevel)
	if err == nil {
		log.SetLevel(logLevel)
	}

	// Log config source and key settings
	log.Infof("Using device ID: %s", cfg.DeviceID)
	log.Infof("API endpoint: %s", cfg.ServerURL)

	// Initialize agent
	a, err := agent.New(cfg, log)
	if err != nil {
		log.Fatalf("Failed to initialize agent: %v", err)
	}

	// Set up signal handling for graceful shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	// Start the agent
	if err := a.Start(); err != nil {
		log.Fatalf("Failed to start agent: %v", err)
	}

	// Wait for shutdown signal
	sig := <-sigCh
	log.Infof("Received signal %v, shutting down...", sig)

	// Stop the agent gracefully
	if err := a.Stop(); err != nil {
		log.Errorf("Error during shutdown: %v", err)
		os.Exit(1)
	}

	log.Info("Agent stopped successfully")
}
