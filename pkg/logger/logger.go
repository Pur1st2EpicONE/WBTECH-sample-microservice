// Package logger provides a logging abstraction for the service.
// It defines a Logger interface and a constructor for creating structured loggers.
package logger

import (
	"os"

	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/configs"
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/pkg/logger/slog"
)

// Logger defines the interface for structured logging with different severity levels.
type Logger interface {
	// LogFatal logs a fatal message with an error and optional key-value arguments.
	LogFatal(msg string, err error, args ...any)
	// LogError logs an error message with an error and optional key-value arguments.
	LogError(string, error, ...any)
	// LogInfo logs an informational message with optional key-value arguments.
	LogInfo(msg string, args ...any)
	// Debug logs a debug message with optional key-value arguments.
	Debug(msg string, args ...any)
}

// NewLogger creates a new Logger instance based on the provided configuration.
// Returns the logger and an *os.File if a file is used for logging.
func NewLogger(config configs.Logger) (Logger, *os.File) {
	return slog.NewLogger(config)
}
