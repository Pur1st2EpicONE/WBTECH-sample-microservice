package slog

import (
	"log"
	"log/slog"
	"os"
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

	}
	handler := slog.NewJSONHandler(logDest, nil)
	logger := &Logger{logger: slog.New(handler)}
	slog.SetDefault(logger.logger)
	return logger, logDest
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

func openFile(logDir string) *os.File {
	if err := os.MkdirAll(logDir, 0777); err != nil {
		log.Fatalf("logger — failed to create log directory: %v", err)
	}
	logFile, err := os.OpenFile("./logs/app.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0777)
	if err != nil {
		log.Fatalf("logger — failed to create log file: %v", err)
	}

	return logFile
}
