package server

import (
	"context"
	"net/http"

	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/cache"
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/configs"
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/handler"
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/logger"
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/repository"
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/service"
)

type Server struct {
	httpServer *http.Server
}

func NewServer(config configs.Server, cache cache.Cache, storage *repository.Storage) *Server {
	service := service.NewService(storage, cache)
	handler := handler.NewHandler(service)
	router := handler.InitRoutes()
	server := new(Server)
	server.serverConfig(config, router)
	return server
}

func (s *Server) serverConfig(config configs.Server, handler http.Handler) {
	s.httpServer = &http.Server{
		Addr:           ":" + config.Port,
		Handler:        handler,
		ReadTimeout:    config.ReadTimeout,
		WriteTimeout:   config.WriteTimeout,
		MaxHeaderBytes: config.MaxHeaderBytes,
	}
}

func (s *Server) Run(ctx context.Context) error {
	logger.LogInfo("server — receiving requests")
	err := s.httpServer.ListenAndServe()
	return err
}

func (s *Server) Shutdown(ctx context.Context) {
	if err := s.httpServer.Shutdown(ctx); err != nil {
		logger.LogError("server — failed to shutdown gracefully", err)
	} else {
		logger.LogInfo("server — shutdown complete")
	}

}
