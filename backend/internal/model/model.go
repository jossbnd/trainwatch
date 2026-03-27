package model

import "time"

// NextTrain is the enriched train departure returned to API clients.
type NextTrain struct {
	EstimatedAt  time.Time  `json:"estimated_at"`
	AimedAt      *time.Time `json:"aimed_at,omitempty"`
	Destination  string     `json:"destination,omitempty"`
	Status       string     `json:"status"`
	DelayMinutes int        `json:"delay_minutes"`
}

// NewNextTrain builds a NextTrain, computing DelayMinutes as the difference
// between estimatedAt and aimedAt. Returns 0 if aimedAt is nil or on time.
func NewNextTrain(estimatedAt time.Time, aimedAt *time.Time, destination, status string) NextTrain {
	delay := 0
	if aimedAt != nil {
		d := int(estimatedAt.Sub(*aimedAt).Minutes())
		if d > 0 {
			delay = d
		}
	}
	return NextTrain{
		EstimatedAt:  estimatedAt,
		AimedAt:      aimedAt,
		Destination:  destination,
		Status:       status,
		DelayMinutes: delay,
	}
}
