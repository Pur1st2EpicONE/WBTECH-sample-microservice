package server

import (
	"context"
	"net/http"

	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/configs"
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/logger"
)

type Server struct {
	HttpServer *http.Server
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
	return server
}

func (s *Server) Run(ctx context.Context, logger logger.Logger) error {
	logger.LogInfo("server — receiving requests")
	err := s.HttpServer.ListenAndServe()
	return err
}

func (s *Server) Shutdown(ctx context.Context, logger logger.Logger) {
	if err := s.HttpServer.Shutdown(ctx); err != nil {
		logger.LogError("server — failed to shutdown gracefully", err)
	} else {
		logger.LogInfo("server — shutdown complete")
	}

}
