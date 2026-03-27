package model

import "time"

// NextTrain is the enriched train departure returned to API clients.
type NextTrain struct {
	EstimatedAt time.Time  `json:"estimated_at"`
	AimedAt     *time.Time `json:"aimed_at,omitempty"`
	Destination string     `json:"destination,omitempty"`
	Status      string     `json:"status"`
	WaitMinutes int        `json:"wait_minutes"`
}

// NewNextTrain builds a NextTrain, computing WaitMinutes automatically.
func NewNextTrain(estimatedAt time.Time, aimedAt *time.Time, destination, status string) NextTrain {
	wait := int(time.Until(estimatedAt).Minutes())
	if wait < 0 {
		wait = 0
	}
	return NextTrain{
		EstimatedAt: estimatedAt,
		AimedAt:     aimedAt,
		Destination: destination,
		Status:      status,
		WaitMinutes: wait,
	}
}
