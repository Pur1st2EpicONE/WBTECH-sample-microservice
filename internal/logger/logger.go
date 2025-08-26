package logger

import (
	"os"

	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/configs"
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/logger/slog"
)

//go:generate mockgen -source=logger.go -destination=mocks/mock.go

type Logger interface {
	LogFatal(msg string, err error, args ...any)
	LogError(string, error, ...any)
	LogInfo(msg string, args ...any)
	Debug(msg string, args ...any)
}

func NewLogger(config configs.Logger) (Logger, *os.File) {
	return slog.NewLogger(config)
}
