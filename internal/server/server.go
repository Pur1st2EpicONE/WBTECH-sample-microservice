// Package server provides a simple HTTP server wrapper with configurable timeouts and graceful shutdown support.
// It uses the standard net/http server under the hood and logs events using the provided logger.
package server

import (
	"context"
	"net/http"
	"time"

	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/configs"
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/pkg/logger"
)

// Server wraps an http.Server and provides shutdown timeout configuration.
type Server struct {
	HttpServer      *http.Server
	ShutdownTimeout time.Duration
}

// NewServer creates a new Server instance with the provided configuration and HTTP handler.
func NewServer(config configs.Server, handler http.Handler) *Server {
	server := new(Server)
	server.HttpServer = &http.Server{
		Addr:           ":" + config.Port,
		Handler:        handler,
		ReadTimeout:    config.ReadTimeout,
		WriteTimeout:   config.WriteTimeout,
		MaxHeaderBytes: config.MaxHeaderBytes,
	}
	server.ShutdownTimeout = config.ShutdownTimeout
	return server
}

// Run starts the HTTP server and logs that it is receiving requests.
func (s *Server) Run(ctx context.Context, logger logger.Logger) error {
	logger.LogInfo("server — receiving requests", "layer", "server")
	return s.HttpServer.ListenAndServe()
}

// Shutdown gracefully shuts down the HTTP server and logs the result.
func (s *Server) Shutdown(ctx context.Context, logger logger.Logger) {
	if err := s.HttpServer.Shutdown(ctx); err != nil {
		logger.LogError("server — failed to shutdown gracefully", err, "layer", "server")
	} else {
		logger.LogInfo("server — shutdown complete", "layer", "server")
	}
}
