package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port     string
	BotToken string
	ChatID   string
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	return &Config{
		Port:     port,
		BotToken: os.Getenv("TELEGRAM_BOT_TOKEN"),
		ChatID:   os.Getenv("TELEGRAM_CHAT_ID"),
	}, nil
}
