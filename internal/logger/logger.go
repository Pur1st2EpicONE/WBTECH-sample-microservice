package logger

import (
	"log/slog"
	"os"
)

func OpenFile() *os.File {
	if err := os.MkdirAll("./logs", 0777); err != nil {
		LogFatal("failed to create log directory: %v", err)
	}
	logFile, err := os.OpenFile("./logs/app.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0777)
	if err != nil {
		LogFatal("failed to create log file: %v", err)
	}

	handler := slog.NewJSONHandler(logFile, nil)
	logger := slog.New(handler)
	slog.SetDefault(logger)

	return logFile
}

func CloseFile(file *os.File) {
	if err := file.Close(); err != nil {
		slog.Error("log file failed to close properly")
	}
}

func LogFatal(msg string, err error) {
	slog.Error(msg, slog.String("critical error", err.Error()))
	os.Exit(1)
}

func LogError(msg string, err error) {
	if err != nil {
		slog.Error(msg, slog.String("error", err.Error()))
	} else {
		slog.Error(msg)
	}
}

func LogInfo(msg string, args ...any) {
	slog.Info(msg, args...)
}
