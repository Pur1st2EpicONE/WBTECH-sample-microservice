package slog

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/configs"
)

type Logger struct {
	logger *slog.Logger
}

func NewLogger(config configs.Logger) (*Logger, *os.File) {
	var logDest *os.File
	if config.LogDir == "" {
		logDest = os.Stdout
	} else {
		logDest = openFile(config.LogDir)
		if logDest == nil {
			logDest = os.Stdout
		}
	}
	var level slog.Level
	if config.Debug {
		level = slog.LevelDebug
	} else {
		level = slog.LevelInfo
	}
	handler := slog.NewJSONHandler(logDest, &slog.HandlerOptions{Level: level})
	logger := &Logger{logger: slog.New(handler)}
	slog.SetDefault(logger.logger)
	return logger, logDest
}

func openFile(logDir string) *os.File {
	if err := os.MkdirAll(logDir, 0777); err != nil {
		fmt.Fprintf(os.Stderr, "logger — failed to create log directory switching to stdout: %v\n", err)
		return nil
	}
	logPath := filepath.Join(logDir, "app.log")
	logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0777)
	if err != nil {
		fmt.Fprintf(os.Stderr, "logger — failed to create log file switching to stdout: %v\n", err)
		return nil
	}
	return logFile
}

func (l *Logger) LogFatal(msg string, err error, args ...any) {
	if err != nil {
		args = append(args, "err", err.Error())
	}
	slog.Error(msg, args...)
	os.Exit(1)
}

func (l *Logger) LogError(msg string, err error, args ...any) {
	if err != nil {
		args = append(args, "err", err.Error())
	}
	slog.Error(msg, args...)
}

func (l *Logger) LogInfo(msg string, args ...any) {
	slog.Info(msg, args...)
}

func (l *Logger) Debug(msg string, args ...any) {
	l.logger.Debug(msg, args...)
}
