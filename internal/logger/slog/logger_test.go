package slog_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/logger"
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/logger/slog"
)

func TestLogger_toDir(t *testing.T) {
	tmpDir := t.TempDir()
	l1, logFile1 := slog.NewLogger(tmpDir)
	defer logFile1.Close()

	l1.LogInfo("very informative log")

	data, err := os.ReadFile(filepath.Join(tmpDir, "app.log"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), "very informative log") {
		t.Errorf("expected log in file, got: %s", string(data))
	}
	l2, logFile2 := slog.NewLogger("/root")

	if logFile2 != os.Stdout {
		t.Errorf("expected stdout fallback, got: %v", logFile2)
	}
	l2.LogInfo("very informative log — stdout")
}

func TestLogger_toStdout(t *testing.T) {
	l, logDest := logger.NewLogger("")
	if logDest != os.Stdout {
		t.Errorf("expected stdout, got %v", logDest)
	}

	l.LogInfo("info message")
	l.LogError("error message", nil)
	l.LogError("error with err", os.ErrInvalid)
}
