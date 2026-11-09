package main

import (
	"log"
	"net/http"

	"github.com/mcypher-alt/telemetry_collector_go/internal/config"
	"github.com/mcypher-alt/telemetry_collector_go/internal/handler"
	"github.com/mcypher-alt/telemetry_collector_go/internal/telegram"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Config error: %v", err)
	}

	bot := telegram.NewBot(cfg.BotToken, cfg.ChatID)
	telemetryHandler := handler.NewTelemetryHandler(bot)

	http.HandleFunc("/api/v1/telemetry", telemetryHandler.Handle)

	log.Printf("Go telemetry server listening on 0.0.0.0:%s...\n", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
