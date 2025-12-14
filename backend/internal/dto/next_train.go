package dto

import "time"

// NextTrain is the HTTP response DTO for the next train endpoint.
type NextTrain struct {
	EstimatedAt time.Time `json:"estimated_at"`
	Status      string    `json:"status"`
}
