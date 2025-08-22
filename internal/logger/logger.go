package logger

import (
	"os"

	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/logger/slog"
)

type Logger interface {
	LogFatal(msg string, err error)
	LogError(msg string, err error)
	LogInfo(msg string, args ...any)
}

func NewLogger(logDir string) (Logger, *os.File) {
	return slog.NewLogger(logDir)
}
