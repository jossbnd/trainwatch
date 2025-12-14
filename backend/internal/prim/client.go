package prim

import (
	"time"
)

// NextTrain is a simplified response returned to the API consumer.
type NextTrain struct {
	EstimatedAt time.Time `json:"estimated_at"`
	Status      string    `json:"status"`
}

// ListNextTrains fetches the next train information from the data source.
func ListNextTrains(transportType, line, station, direction string) ([]NextTrain, error) {
	return []NextTrain{
		{
			EstimatedAt: time.Now().Add(5 * time.Minute),
			Status:      "On time",
		},
	}, nil
}
