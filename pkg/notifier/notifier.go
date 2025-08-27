package notifier

import (
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/configs"
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/pkg/notifier/telegram"
)

type Notifier interface {
	Notify(message string) error
}

func NewNotifier(config configs.Notifier) Notifier {
	return telegram.NewNotifier(config)
}
