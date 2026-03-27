package model

import "time"

// TextValue is a PRIM typed text field.
type TextValue struct {
	Value string `json:"value"`
}

// MonitoredCall holds departure/arrival timing fields from PRIM.
type MonitoredCall struct {
	ExpectedDepartureTime *time.Time `json:"ExpectedDepartureTime,omitempty"`
	AimedDepartureTime    *time.Time `json:"AimedDepartureTime,omitempty"`
	ExpectedArrivalTime   *time.Time `json:"ExpectedArrivalTime,omitempty"`
	AimedArrivalTime      *time.Time `json:"AimedArrivalTime,omitempty"`
	DepartureStatus       string     `json:"DepartureStatus,omitempty"`
	ArrivalStatus         string     `json:"ArrivalStatus,omitempty"`
}

// MonitoredVehicleJourney holds journey identification fields from PRIM.
type MonitoredVehicleJourney struct {
	DirectionRef    TextValue     `json:"DirectionRef"`
	DirectionName   []TextValue   `json:"DirectionName"`
	DestinationName []TextValue   `json:"DestinationName"`
	MonitoredCall   MonitoredCall `json:"MonitoredCall"`
}

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
