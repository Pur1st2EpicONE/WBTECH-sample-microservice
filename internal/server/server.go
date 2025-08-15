package server

import (
	"context"
	"net/http"
	"time"

	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/cache"
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/handler"
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/logger"
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/repository"
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/service"
)

type Server struct {
	httpServer *http.Server
}

func NewServer(port string, cache *cache.Cache, storage *repository.Storage) *Server {
	service := service.NewService(storage, cache)
	handler := handler.NewHandler(service)
	router := handler.InitRoutes()
	server := new(Server)
	server.serverConfig(port, router)
	return server
}

func (s *Server) serverConfig(port string, handler http.Handler) {
	s.httpServer = &http.Server{
		Addr:           ":" + port,
		Handler:        handler,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
}

func (s *Server) Run(ctx context.Context) {
	logger.LogInfo("server — receiving requests")
	err := s.httpServer.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		logger.LogFatal("server — run failed", err)
	}
}

func (s *Server) Shutdown(ctx context.Context) error {
	logger.LogInfo("server — shutting down")
	return s.httpServer.Shutdown(ctx)
}
