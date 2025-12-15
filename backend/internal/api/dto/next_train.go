package dto

import "time"

type NextTrain struct {
	EstimatedAt time.Time `json:"estimated_at"`
	Status      string    `json:"status"`
}
