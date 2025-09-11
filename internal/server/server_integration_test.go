package server_test

import (
	"context"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/configs"
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/server"
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/pkg/logger"
)

func TestServer_RunAndShutdown_WithLogger(t *testing.T) {
	tmpDir := t.TempDir()
	config := configs.Logger{LogDir: tmpDir, Debug: false}
	logger, logFile := logger.NewLogger(config)
	defer func() { _ = logFile.Close() }()

	handlerCalled := false
	fakeHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true
		w.WriteHeader(http.StatusOK)
	})

	srv := &server.Server{
		HttpServer: &http.Server{
			Addr:    "localhost:8086",
			Handler: fakeHandler,
		},
	}

	go func() {
		_ = srv.Run(context.Background(), logger)
	}()

	time.Sleep(100 * time.Millisecond)

	resp, err := http.Get("http://" + srv.HttpServer.Addr)
	if err != nil {
		t.Fatalf("failed to GET: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200 OK, got %d", resp.StatusCode)
	}
	if !handlerCalled {
		t.Error("handler was not called")
	}

	srv.Shutdown(context.Background(), logger)

	logPath := filepath.Join(tmpDir, "app.log")
	data, err := os.ReadFile(logPath)
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
