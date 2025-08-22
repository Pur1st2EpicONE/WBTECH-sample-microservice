package server_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/configs"
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/server"
)

func TestNewServer_ConfigApplied(t *testing.T) {
	cfg := configs.Server{
		Port:           "1234",
		ReadTimeout:    2 * time.Second,
		WriteTimeout:   3 * time.Second,
		MaxHeaderBytes: 4096,
	}

	fakeHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	srv := server.NewServer(cfg, fakeHandler)

	if srv.HttpServer == nil {
		t.Fatal("expected HttpServer initialized, got nil")
	}
	if srv.HttpServer.Addr != ":1234" {
		t.Errorf("expected Addr :1234, got %s", srv.HttpServer.Addr)
	}
	if srv.HttpServer.Handler == nil {
		t.Error("expected handler set, got nil")
	}
	if srv.HttpServer.ReadTimeout != 2*time.Second {
		t.Errorf("expected ReadTimeout=2s, got %v", srv.HttpServer.ReadTimeout)
	}
	if srv.HttpServer.WriteTimeout != 3*time.Second {
		t.Errorf("expected WriteTimeout=3s, got %v", srv.HttpServer.WriteTimeout)
	}
	if srv.HttpServer.MaxHeaderBytes != 4096 {
		t.Errorf("expected MaxHeaderBytes=4096, got %d", srv.HttpServer.MaxHeaderBytes)
	}
}
