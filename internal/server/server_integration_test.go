package server_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/logger/slog"
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/server"
)

func TestServer_RunAndShutdown_WithSlog(t *testing.T) {
	logFile, err := os.CreateTemp("", "server_log_*.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(logFile.Name())
	defer logFile.Close()

	logger, _ := slog.NewLogger("")
	handlerCalled := false
	fakeHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true
		w.WriteHeader(http.StatusOK)
	})

	ts := httptest.NewServer(fakeHandler)
	defer ts.Close()

	srv := &server.Server{
		HttpServer: &http.Server{
			Addr:    ts.Listener.Addr().String(),
			Handler: fakeHandler,
		},
	}

	go func() {
		_ = srv.Run(context.Background(), logger)
	}()

	time.Sleep(50 * time.Millisecond)

	resp, err := http.Get("http://" + ts.Listener.Addr().String())
	if err != nil {
		t.Fatalf("failed to GET: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200 OK, got %d", resp.StatusCode)
	}
	if !handlerCalled {
		t.Error("handler was not called")
	}

	srv.Shutdown(context.Background(), logger)

	data, err := os.ReadFile(logFile.Name())
	if err != nil {
		t.Fatal(err)
	}
	content := string(data)
	if !strings.Contains(content, "server — receiving requests") {
		t.Error("log does not contain 'server — receiving requests'")
	}
	if !strings.Contains(content, "server — shutdown complete") {
		t.Error("log does not contain 'server — shutdown complete'")
	}
}
