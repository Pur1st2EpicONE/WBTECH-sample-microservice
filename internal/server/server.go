package server

import (
	"context"
	"net/http"
	"time"

	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/cache"
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/handler"
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/repository"
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/service"
)

type Server struct {
	httpServer *http.Server
}

func InitServer(port string, storage *repository.Storage) *Server {
	cache := cache.NewCache(24 * time.Hour)
	service := service.NewService(storage, cache)
	handler := handler.NewHandler(service)
	router := handler.InitRoutes()
	server := new(Server)
	server.serverPrep(port, router)
	return server
}

func (s *Server) serverPrep(port string, handler http.Handler) {
	s.httpServer = &http.Server{
		Addr:           ":" + port,
		Handler:        handler,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
}

func (s *Server) Run() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
