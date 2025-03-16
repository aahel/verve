package server

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"verve/config"
	"verve/stats"
)

// Server represents the HTTP server handling requests
type Server struct {
	server    *http.Server
	collector *stats.Collector
	logger    stats.Logger
}

// NewServer creates a new server instance
func NewServer(cfg *config.Config, collector *stats.Collector, logger stats.Logger) *Server {
	mux := http.NewServeMux()

	server := &Server{
		server: &http.Server{
			Addr:         cfg.Server.Addr,
			Handler:      mux,
			ReadTimeout:  cfg.Server.ReadTimeout,
			WriteTimeout: cfg.Server.WriteTimeout,
		},
		collector: collector,
		logger:    logger,
	}

	// Register routes
	mux.HandleFunc("/api/verve/accept", server.handleAccept)

	return server
}

// Start starts the HTTP server
func (s *Server) Start() error {
	return s.server.ListenAndServe()
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	return s.server.Shutdown(ctx)
}

// handleAccept handles the /api/verve/accept endpoint
func (s *Server) handleAccept(w http.ResponseWriter, r *http.Request) {
	// Check request method
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get the id parameter
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "Missing required id parameter", http.StatusBadRequest)
		return
	}

	// Parse id as an integer
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid id parameter, must be an integer", http.StatusBadRequest)
		return
	}

	// Get the optional endpoint parameter
	endpoint := r.URL.Query().Get("endpoint")

	// Process the request
	err = s.collector.ProcessRequest(id, endpoint)
	if err != nil {
		s.logger.Printf("Error processing request: %v", err)
		fmt.Fprintln(w, "failed")
		return
	}

	fmt.Fprintln(w, "ok")
}
