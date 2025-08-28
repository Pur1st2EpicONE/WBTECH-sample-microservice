package server

import (
	"context"
	"net/http"
	"time"

	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/configs"
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/pkg/logger"
)

type Server struct {
	HttpServer      *http.Server
	ShutdownTimeout time.Duration
}

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

func (s *Server) Run(ctx context.Context, logger logger.Logger) error {
	logger.LogInfo("server — receiving requests", "layer", "server")
	return s.HttpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context, logger logger.Logger) {
	if err := s.HttpServer.Shutdown(ctx); err != nil {
		logger.LogError("server — failed to shutdown gracefully", err, "layer", "server")
	} else {
		logger.LogInfo("server — shutdown complete", "layer", "server")
	}
}
