package model

import "time"

// Departure is the enriched train departure returned to API clients.
type Departure struct {
	EstimatedAt  time.Time  `json:"estimated_at"`
	AimedAt      *time.Time `json:"aimed_at,omitempty"`
	Destination  string     `json:"destination,omitempty"`
	Status       string     `json:"status"`
	DelayMinutes int        `json:"delay_minutes"`
}

// NewDeparture builds a Departure, computing DelayMinutes as the difference
// between estimatedAt and aimedAt. Returns 0 if aimedAt is nil or on time.
func NewDeparture(estimatedAt time.Time, aimedAt *time.Time, destination, status string) Departure {
	delay := 0
	if aimedAt != nil {
		d := int(estimatedAt.Sub(*aimedAt).Minutes())
		if d > 0 {
			delay = d
		}
	}
	return Departure{
		EstimatedAt:  estimatedAt,
		AimedAt:      aimedAt,
		Destination:  destination,
		Status:       status,
		DelayMinutes: delay,
	}
}
