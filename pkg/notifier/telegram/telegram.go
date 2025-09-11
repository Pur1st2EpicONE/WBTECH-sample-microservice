// Package telegram provides a Telegram notifier implementation for the service.
// It allows sending messages to a specified chat using the Telegram Bot API.
package telegram

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/configs"
)

// Telegram represents a Telegram bot notifier.
type Telegram struct {
	Token  string       // Bot token
	ChatID string       // Receiver chat ID
	Client *http.Client // Optional custom HTTP client
}

// NewNotifier creates a new Telegram notifier based on the provided configuration.
func NewNotifier(config configs.Notifier) *Telegram {
	return &Telegram{Token: config.Token, ChatID: config.Receiver}
}

// Notify sends a message to the configured Telegram chat using the Bot API.
// Returns an error if the request fails or the API responds with a non-OK status.
func (t *Telegram) Notify(message string) error {
	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", t.Token)
	data := url.Values{}
	data.Set("chat_id", t.ChatID)
	data.Set("text", message)

	client := t.Client
	if client == nil {
		client = new(http.Client)
	}

	resp, err := client.PostForm(apiURL, data)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("telegram API returned status %s", resp.Status)
	}
	return nil
}
