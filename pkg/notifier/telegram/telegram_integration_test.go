package telegram

import (
	"testing"

	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/configs"
	"github.com/stretchr/testify/require"
)

func TestTelegram_Notify_Integration(t *testing.T) {
	t.Skip("pkg/notifier/telegram/integration_test — test requires real Telegram bot token and chat ID")
	tg := NewNotifier(configs.Notifier{
		Token:    "TOKEN",
		Receiver: "805862991",
	})

	err := tg.Notify("tg integration test passed")
	require.NoError(t, err)
}
