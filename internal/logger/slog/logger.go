package slog

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
)

type Logger struct {
	logger *slog.Logger
}

func NewLogger(logDir string) (*Logger, *os.File) {
	var logDest *os.File
	if logDir == "" {
		logDest = os.Stdout
	} else {
		logDest = openFile(logDir)
		if logDest == nil {
			logDest = os.Stdout
		}
	}
	handler := slog.NewJSONHandler(logDest, nil)
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

func (l *Logger) LogFatal(msg string, err error) {
	slog.Error(msg, slog.String("critical error", err.Error()))
	os.Exit(1)
}

func (l *Logger) LogError(msg string, err error) {
	if err != nil {
		slog.Error(msg, slog.String("error", err.Error()))
	} else {
		slog.Error(msg)
	}
}

func (l *Logger) LogInfo(msg string, args ...any) {
	slog.Info(msg, args...)
}
