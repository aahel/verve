package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"verve/config"
	"verve/server"
	"verve/stats"
)

// StdLogger implements stats.Logger interface
type StdLogger struct {
	logger *log.Logger
}

func (l *StdLogger) Printf(format string, v ...interface{}) {
	l.logger.Printf(format, v...)
}

func (l *StdLogger) Fatalf(format string, v ...interface{}) {
	l.logger.Fatalf(format, v...)
}

func main() {
	// Initialize logger
	stdLogger := log.New(os.Stdout, "verve: ", log.LstdFlags)
	logger := &StdLogger{logger: stdLogger}

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		logger.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize stats collector
	collector, err := stats.NewCollector(cfg, logger)
	if err != nil {
		logger.Fatalf("Failed to initialize stats collector: %v", err)
	}
	defer collector.Close()

	// Start the stats processing in background
	go collector.Run()

	// Initialize and start the HTTP server
	srv := server.NewServer(cfg, collector, logger)
	go func() {
		if err := srv.Start(); err != nil {
			logger.Fatalf("Server failed: %v", err)
		}
	}()

	logger.Printf("Server started on %s", cfg.Server.Addr)

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Printf("Shutting down server...")

	if err := srv.Shutdown(); err != nil {
		logger.Fatalf("Server shutdown failed: %v", err)
	}

	logger.Printf("Server stopped")
}
