package model

type TelemetryPayload struct {
	DeviceID string `json:"device_id"`
	Pin      int    `json:"pin"`
	State    string `json:"state"`
}
