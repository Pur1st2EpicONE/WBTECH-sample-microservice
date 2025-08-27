package slog_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/configs"
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/pkg/logger"
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/pkg/logger/slog"
)

func TestLogger_toDir(t *testing.T) {
	tmpDir := t.TempDir()
	config := configs.Logger{LogDir: tmpDir, Debug: false}
	firstLog, firstLogFile := slog.NewLogger(config)
	defer firstLogFile.Close()

	firstLog.LogInfo("very informative log")

	data, err := os.ReadFile(filepath.Join(tmpDir, "app.log"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), "very informative log") {
		t.Errorf("expected log in file, got: %s", string(data))
	}
	config = configs.Logger{LogDir: "/root", Debug: false}
	secondLog, secondLogFile := slog.NewLogger(config)

	if secondLogFile != os.Stdout {
		t.Errorf("expected stdout fallback, got: %v", secondLogFile)
	}
	secondLog.LogInfo("very informative log — stdout")
}

func TestLogger_toStdout(t *testing.T) {
	config := configs.Logger{LogDir: "/dev/full/test", Debug: false}
	log, logDest := logger.NewLogger(config)
	if logDest != os.Stdout {
		t.Errorf("expected stdout, got %v", logDest)
	}

	log.LogInfo("info message")
	log.LogError("error message", nil)
	log.LogError("error with err", os.ErrInvalid)
}

func TestLogger_DEBUG(t *testing.T) {
	config := configs.Logger{LogDir: "", Debug: true}
	log, _ := logger.NewLogger(config)
	log.Debug("debug message")
}
