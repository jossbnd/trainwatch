package prim

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type TextValue struct {
	Value string `json:"value"`
}

type Timing struct {
	ExpectedDepartureTime *time.Time `json:"ExpectedDepartureTime,omitempty"`
	AimedDepartureTime    *time.Time `json:"AimedDepartureTime,omitempty"`
	ExpectedArrivalTime   *time.Time `json:"ExpectedArrivalTime,omitempty"`
	AimedArrivalTime      *time.Time `json:"AimedArrivalTime,omitempty"`
	DepartureStatus       string     `json:"DepartureStatus,omitempty"`
	ArrivalStatus         string     `json:"ArrivalStatus,omitempty"`
}

// StopVisit represents a train stopping at a station, as returned by the PRIM API.
type StopVisit struct {
	DirectionRef    TextValue     `json:"DirectionRef"`
	DirectionName   []TextValue   `json:"DirectionName"`
	DestinationName []TextValue   `json:"DestinationName"`
	Timing   Timing `json:"Timing"`
}

// Client defines the interface for interacting with the PRIM stop-monitoring API.
type Client interface {
	// FetchStopVisits queries the PRIM stop-monitoring endpoint using the given
	// stopRef and lineRef, parses the SIRI response, and returns a slice of
	// StopVisit for trains stopping at the station.
	FetchStopVisits(ctx context.Context, stopRef, lineRef string) ([]StopVisit, error)
}

type client struct {
	baseURL string
	apiKey  string
	http    *http.Client
}

func New(baseURL, apiKey string) (Client, error) {
	// validate scheme
	if !(strings.HasPrefix(baseURL, "http://") || strings.HasPrefix(baseURL, "https://")) {
		return nil, fmt.Errorf("invalid PRIM base URL scheme")
	}

	// normalize baseURL: strip trailing slash
	normalizedBaseURL := strings.TrimRight(baseURL, "/")

	return &client{
		baseURL: normalizedBaseURL,
		apiKey:  apiKey,
		http:    &http.Client{Timeout: 10 * time.Second},
	}, nil
}

func (c *client) FetchStopVisits(ctx context.Context, stopRef, lineRef string) ([]StopVisit, error) {
	// Build request
	endpoint := fmt.Sprintf("%s/marketplace/stop-monitoring", c.baseURL)
	req, err := http.NewRequestWithContext(ctx, "GET", endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to build request: %w", err)
	}

	// Set query parameters
	q := req.URL.Query()
	q.Set("MonitoringRef", stopRef)
	q.Set("LineRef", lineRef)
	req.URL.RawQuery = q.Encode()

	// Set headers
	req.Header.Set("Accept", "application/json")
	if c.apiKey != "" {
		req.Header.Set("apiKey", c.apiKey)
	}

	// Execute request
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Check for non-2xx status codes
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		b, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("prim returned status %d: %s", resp.StatusCode, string(b))
	}

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Unmarshal only the useful fields from the SIRI response.
	var wrapper struct {
		Siri struct {
			ServiceDelivery struct {
				StopMonitoringDelivery []struct {
					MonitoredStopVisit []struct {
						MonitoredVehicleJourney StopVisit `json:"MonitoredVehicleJourney"`
					} `json:"MonitoredStopVisit"`
				} `json:"StopMonitoringDelivery"`
			} `json:"ServiceDelivery"`
		} `json:"Siri"`
	}
	if err := json.Unmarshal(body, &wrapper); err != nil {
		return nil, fmt.Errorf("failed to parse prim response: %w", err)
	}

	// Aggregate stop visits from all deliveries.
	var visits []StopVisit
	for _, delivery := range wrapper.Siri.ServiceDelivery.StopMonitoringDelivery {
		for _, sv := range delivery.MonitoredStopVisit {
			visits = append(visits, sv.MonitoredVehicleJourney)
		}
	}

	return visits, nil
}
