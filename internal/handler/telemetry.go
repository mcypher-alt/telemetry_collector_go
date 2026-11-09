package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"github.com/mcypher-alt/telemetry_collector_go/internal/model"
	"github.com/mcypher-alt/telemetry_collector_go/internal/telegram"
)

type TelemetryHandler struct {
	bot *telegram.Bot
}

func NewTelemetryHandler(bot *telegram.Bot) *TelemetryHandler {
	return &TelemetryHandler{bot: bot}
}

func (h *TelemetryHandler) Handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var payload model.TelemetryPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	fmt.Printf("[EVENT] Device: %s | Pin: %d | State: %s\n", payload.DeviceID, payload.Pin, payload.State)

	if payload.State == "CLOSED" {
		alertMsg := fmt.Sprintf("⚠️ Внимание! Сработал датчик на %s (пин %d)", payload.DeviceID, payload.Pin)
		_ = h.bot.SendAlert(alertMsg)
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"status":"ok"}`))
}
