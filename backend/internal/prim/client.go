package prim

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
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

// ResponseError represents a non-2xx response from the PRIM API.
type ResponseError struct {
	StatusCode int
	Body       string
}

func (e *ResponseError) Error() string {
	return fmt.Sprintf("prim returned status %d: %s", e.StatusCode, e.Body)
}

func (e *ResponseError) IsClientError() bool {
	return e.StatusCode >= 400 && e.StatusCode < 500
}

// StopVisit represents a train stopping at a station, as returned by the PRIM API.
type StopVisit struct {
	DirectionRef    TextValue   `json:"DirectionRef"`
	DirectionName   []TextValue `json:"DirectionName"`
	DestinationName []TextValue `json:"DestinationName"`
	Timing          Timing      `json:"MonitoredCall"`
}

// Client defines the interface for interacting with the PRIM stop-monitoring API.
type Client interface {
	// FetchStopVisits queries the PRIM stop-monitoring endpoint using the given
	// stopRef and lineRef, parses the SIRI response, and returns a slice of
	// StopVisit and the remaining daily API credits (-1 if unavailable).
	FetchStopVisits(ctx context.Context, stopRef, lineRef string) ([]StopVisit, int, error)
}

type client struct {
	baseURL string
	apiKey  string
	http    *http.Client
}

func New(baseURL, apiKey string) (Client, error) {
	if !strings.HasPrefix(baseURL, "http://") && !strings.HasPrefix(baseURL, "https://") {
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

func (c *client) FetchStopVisits(ctx context.Context, stopRef, lineRef string) ([]StopVisit, int, error) {
	endpoint := fmt.Sprintf("%s/marketplace/stop-monitoring", c.baseURL)
	req, err := http.NewRequestWithContext(ctx, "GET", endpoint, nil)
	if err != nil {
		return nil, -1, fmt.Errorf("failed to build request: %w", err)
	}

	q := req.URL.Query()
	q.Set("MonitoringRef", stopRef)
	q.Set("LineRef", lineRef)
	req.URL.RawQuery = q.Encode()

	req.Header.Set("Accept", "application/json")
	if c.apiKey != "" {
		req.Header.Set("apiKey", c.apiKey)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, -1, fmt.Errorf("request failed: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	// Parse remaining daily credits from response header.
	credits := -1
	if v := resp.Header.Get("x-ratelimit-remaining-day"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			credits = n
		}
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		b, _ := io.ReadAll(resp.Body)
		return nil, credits, &ResponseError{StatusCode: resp.StatusCode, Body: string(b)}
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, credits, fmt.Errorf("failed to read response body: %w", err)
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
		return nil, credits, fmt.Errorf("failed to parse prim response: %w", err)
	}

	var visits []StopVisit
	for _, delivery := range wrapper.Siri.ServiceDelivery.StopMonitoringDelivery {
		for _, sv := range delivery.MonitoredStopVisit {
			visits = append(visits, sv.MonitoredVehicleJourney)
		}
	}

	return visits, credits, nil
}
