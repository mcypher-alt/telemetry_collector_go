package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Bot struct {
	token      string
	chatID     string
	httpClient *http.Client
}

func NewBot(token, chatID string) *Bot {
	return &Bot{
		token:  token,
		chatID: chatID,
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (b *Bot) SendAlert(msg string) error {
	if b.chatID == "" || b.token == "" {
		return nil
	}

	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", b.token)
	body, _ := json.Marshal(map[string]string{
		"chat_id": b.chatID,
		"text":    msg,
	})

	resp, err := b.httpClient.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("telegram API returned status: %d", resp.StatusCode)
	}

	return nil
}
