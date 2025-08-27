package telegram_test

import (
	"net/http"
	"testing"

	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/configs"
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/pkg/notifier/telegram"
	"github.com/stretchr/testify/require"
)

func TestTelegram_Notify_HTTPError(t *testing.T) {
	tg := telegram.NewNotifier(configs.Notifier{
		Token:    "TOKEN",
		Receiver: "12345",
	})
	tg.Client = new(http.Client)
	err := tg.Notify("aboba")
	require.Error(t, err)
}
