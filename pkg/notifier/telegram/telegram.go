package telegram

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/configs"
)

type Telegram struct {
	Token  string
	ChatID string
	Client *http.Client
}

func NewNotifier(config configs.Notifier) *Telegram {
	return &Telegram{Token: config.Token, ChatID: config.Receiver}
}

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
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("telegram API returned status %s", resp.Status)
	}
	return nil
}
