package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/jossbnd/trainwatch/backend/internal/logger"
	"github.com/jossbnd/trainwatch/backend/internal/prim"
)

// --- mock prim client ---

type mockClient struct {
	visits []prim.StopVisit
	err    error
}

func (m *mockClient) FetchStopVisits(_ context.Context, _, _ string) ([]prim.StopVisit, int, error) {
	return m.visits, -1, m.err
}

// --- helpers ---

func newService(visits []prim.StopVisit, err error) Service {
	return New(Input{
		Logger:     logger.NewDiscard(),
		PrimClient: &mockClient{visits: visits, err: err},
	})
}

// futureVisit builds a StopVisit departing minutesFromNow minutes from now,
// delayed by delayMin minutes (0 = on time).
func futureVisit(minutesFromNow, delayMin int, destination, directionRef, status string) prim.StopVisit {
	aimed := time.Now().Add(time.Duration(minutesFromNow) * time.Minute)
	expected := aimed.Add(time.Duration(delayMin) * time.Minute)
	return prim.StopVisit{
		DirectionRef:    prim.TextValue{Value: directionRef},
		DestinationName: []prim.TextValue{{Value: destination}},
		Timing: prim.Timing{
			AimedDepartureTime:    &aimed,
			ExpectedDepartureTime: &expected,
			DepartureStatus:       status,
		},
	}
}

// --- tests ---

// Test 1: prim client error is propagated
func TestGetDepartures_ClientError(t *testing.T) {
	svc := newService(nil, errors.New("network error"))
	_, err := svc.GetDepartures(context.Background(), "stop1", "line1", "", 5)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

// Test 2: empty visits returns empty slice (not nil)
func TestGetDepartures_EmptyVisits(t *testing.T) {
	svc := newService([]prim.StopVisit{}, nil)
	departures, err := svc.GetDepartures(context.Background(), "stop1", "line1", "", 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(departures) != 0 {
		t.Fatalf("expected 0 departures, got %d", len(departures))
	}
}

// Test 3: past departures (> 1 min ago) are filtered out
func TestGetDepartures_PastDeparturesFiltered(t *testing.T) {
	past := time.Now().Add(-5 * time.Minute)
	visits := []prim.StopVisit{
		{Timing: prim.Timing{ExpectedDepartureTime: &past}},
	}
	svc := newService(visits, nil)
	departures, err := svc.GetDepartures(context.Background(), "s", "l", "", 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(departures) != 0 {
		t.Fatalf("expected 0 departures, got %d", len(departures))
	}
}

// Test 4: direction filter by DirectionRef (exact, case-insensitive)
func TestGetDepartures_DirectionFilterByRef(t *testing.T) {
	visits := []prim.StopVisit{
		futureVisit(10, 0, "Dest A", "A", "onTime"),
		futureVisit(15, 0, "Dest B", "B", "onTime"),
	}
	svc := newService(visits, nil)
	departures, err := svc.GetDepartures(context.Background(), "s", "l", "a", 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(departures) != 1 {
		t.Fatalf("expected 1 departure, got %d", len(departures))
	}
}

// Test 5: direction filter by DestinationName (substring, case-insensitive)
func TestGetDepartures_DirectionFilterByDestination(t *testing.T) {
	visits := []prim.StopVisit{
		futureVisit(5, 0, "Gare du Nord", "", "onTime"),
		futureVisit(8, 0, "Chatelet", "", "onTime"),
	}
	svc := newService(visits, nil)
	departures, err := svc.GetDepartures(context.Background(), "s", "l", "nord", 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(departures) != 1 {
		t.Fatalf("expected 1 departure, got %d", len(departures))
	}
	if departures[0].Destination != "Gare du Nord" {
		t.Errorf("expected destination 'Gare du Nord', got %q", departures[0].Destination)
	}
}

// Test 6: results are sorted ascending by departure time
func TestGetDepartures_SortedAscending(t *testing.T) {
	visits := []prim.StopVisit{
		futureVisit(20, 0, "B", "", ""),
		futureVisit(5, 0, "A", "", ""),
		futureVisit(10, 0, "C", "", ""),
	}
	svc := newService(visits, nil)
	departures, err := svc.GetDepartures(context.Background(), "s", "l", "", 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(departures) != 3 {
		t.Fatalf("expected 3 departures, got %d", len(departures))
	}
	if !departures[0].EstimatedAt.Before(departures[1].EstimatedAt) || !departures[1].EstimatedAt.Before(departures[2].EstimatedAt) {
		t.Error("trains are not sorted ascending")
	}
}

// Test 7: delay_minutes reflects gap between expected and aimed departure
func TestGetDepartures_DelayMinutes(t *testing.T) {
	visits := []prim.StopVisit{
		futureVisit(10, 3, "A", "", "delayed"), // 3 min late
		futureVisit(20, 0, "B", "", "onTime"),  // on time
	}
	svc := newService(visits, nil)
	departures, err := svc.GetDepartures(context.Background(), "s", "l", "", 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if departures[0].DelayMinutes != 3 {
		t.Errorf("expected delay 3, got %d", departures[0].DelayMinutes)
	}
	if departures[1].DelayMinutes != 0 {
		t.Errorf("expected delay 0, got %d", departures[1].DelayMinutes)
	}
}

// Test 8: limit is enforced
func TestGetDepartures_LimitEnforced(t *testing.T) {
	visits := []prim.StopVisit{
		futureVisit(5, 0, "A", "", ""),
		futureVisit(10, 0, "B", "", ""),
		futureVisit(15, 0, "C", "", ""),
		futureVisit(20, 0, "D", "", ""),
	}
	svc := newService(visits, nil)
	departures, err := svc.GetDepartures(context.Background(), "s", "l", "", 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(departures) != 2 {
		t.Fatalf("expected 2 departures (limit=2), got %d", len(departures))
	}
}
