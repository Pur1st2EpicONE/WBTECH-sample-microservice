// Package notifier provides an abstraction for sending notifications in the service.
// Currently, it supports sending messages via Telegram, but can be extended for other channels.
package notifier

import (
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/configs"
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/pkg/notifier/telegram"
)

// Notifier defines the interface for sending notifications.
type Notifier interface {
	// Notify sends a message and returns an error if the notification fails.
	Notify(message string) error
}

// NewNotifier creates a new Notifier instance based on the provided configuration.
// Currently, it returns a Telegram notifier.
func NewNotifier(config configs.Notifier) Notifier {
	return telegram.NewNotifier(config)
}
