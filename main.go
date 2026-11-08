package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	// Загружаем переменные из .env
	if err := godotenv.Load(); err != nil {
		log.Println("Предупреждение: .env файл не найден, читаем системные переменные")
	}

	botToken := os.Getenv("BOT_TOKEN")
	chatID := os.Getenv("chat_ID")
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Сервис инициализирован на порту :%s", port)
	if botToken != "" {
		log.Println("Telegram токен успешно считан")
	}
	if chatID != "" {
		log.Println("ChatID успешно считан")
	}
}
