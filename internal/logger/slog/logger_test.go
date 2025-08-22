package slog_test

import (
	"os"
	"strings"
	"sync"
	"testing"

	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/logger/slog"
)

func TestLogger_LogInfo_LogError(t *testing.T) {
	file, err := os.CreateTemp("", "log_test_*.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(file.Name())

	var wg sync.WaitGroup
	wg.Add(1)

	logger, _ := slog.NewLogger("")

	logger.LogInfo("info message", "arg1", 123)
	logger.LogError("error message", nil)
	logger.LogError("error with err", os.ErrInvalid)
	wg.Done()

	data, err := os.ReadFile(file.Name())
	if err != nil {
		t.Fatal(err)
	}

	content := string(data)
	if !strings.Contains(content, "info message") {
		t.Errorf("expected info message in log, got: %s", content)
	}
	if !strings.Contains(content, "error message") || !strings.Contains(content, "error with err") {
		t.Errorf("expected error messages in log, got: %s", content)
	}
}
