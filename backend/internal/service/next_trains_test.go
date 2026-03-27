package service

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"testing"
	"time"

	"github.com/jossbnd/trainwatch/backend/internal/prim"
)

// --- mock prim client ---

type mockClient struct {
	visits []prim.StopVisit
	err    error
}

func (m *mockClient) FetchStopVisits(_ context.Context, _, _ string) ([]prim.StopVisit, error) {
	return m.visits, m.err
}

// --- helpers ---

func discardLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

func newService(visits []prim.StopVisit, err error) Service {
	return New(Input{
		Logger:     discardLogger(),
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
func TestGetNextTrains_ClientError(t *testing.T) {
	svc := newService(nil, errors.New("network error"))
	_, err := svc.GetNextTrains(context.Background(), "stop1", "line1", "", 5)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

// Test 2: empty visits returns empty slice (not nil)
func TestGetNextTrains_EmptyVisits(t *testing.T) {
	svc := newService([]prim.StopVisit{}, nil)
	trains, err := svc.GetNextTrains(context.Background(), "stop1", "line1", "", 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(trains) != 0 {
		t.Fatalf("expected 0 trains, got %d", len(trains))
	}
}

// Test 3: past departures (> 1 min ago) are filtered out
func TestGetNextTrains_PastDeparturesFiltered(t *testing.T) {
	past := time.Now().Add(-5 * time.Minute)
	visits := []prim.StopVisit{
		{Timing: prim.Timing{ExpectedDepartureTime: &past}},
	}
	svc := newService(visits, nil)
	trains, err := svc.GetNextTrains(context.Background(), "s", "l", "", 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(trains) != 0 {
		t.Fatalf("expected 0 trains, got %d", len(trains))
	}
}

// Test 4: direction filter by DirectionRef (exact, case-insensitive)
func TestGetNextTrains_DirectionFilterByRef(t *testing.T) {
	visits := []prim.StopVisit{
		futureVisit(10, 0, "Dest A", "A", "onTime"),
		futureVisit(15, 0, "Dest B", "B", "onTime"),
	}
	svc := newService(visits, nil)
	trains, err := svc.GetNextTrains(context.Background(), "s", "l", "a", 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(trains) != 1 {
		t.Fatalf("expected 1 train, got %d", len(trains))
	}
}

// Test 5: direction filter by DestinationName (substring, case-insensitive)
func TestGetNextTrains_DirectionFilterByDestination(t *testing.T) {
	visits := []prim.StopVisit{
		futureVisit(5, 0, "Gare du Nord", "", "onTime"),
		futureVisit(8, 0, "Chatelet", "", "onTime"),
	}
	svc := newService(visits, nil)
	trains, err := svc.GetNextTrains(context.Background(), "s", "l", "nord", 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(trains) != 1 {
		t.Fatalf("expected 1 train, got %d", len(trains))
	}
	if trains[0].Destination != "Gare du Nord" {
		t.Errorf("expected destination 'Gare du Nord', got %q", trains[0].Destination)
	}
}

// Test 6: results are sorted ascending by departure time
func TestGetNextTrains_SortedAscending(t *testing.T) {
	visits := []prim.StopVisit{
		futureVisit(20, 0, "B", "", ""),
		futureVisit(5, 0, "A", "", ""),
		futureVisit(10, 0, "C", "", ""),
	}
	svc := newService(visits, nil)
	trains, err := svc.GetNextTrains(context.Background(), "s", "l", "", 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(trains) != 3 {
		t.Fatalf("expected 3 trains, got %d", len(trains))
	}
	if !trains[0].EstimatedAt.Before(trains[1].EstimatedAt) || !trains[1].EstimatedAt.Before(trains[2].EstimatedAt) {
		t.Error("trains are not sorted ascending")
	}
}

// Test 7: delay_minutes reflects gap between expected and aimed departure
func TestGetNextTrains_DelayMinutes(t *testing.T) {
	visits := []prim.StopVisit{
		futureVisit(10, 3, "A", "", "delayed"), // 3 min late
		futureVisit(20, 0, "B", "", "onTime"),  // on time
	}
	svc := newService(visits, nil)
	trains, err := svc.GetNextTrains(context.Background(), "s", "l", "", 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if trains[0].DelayMinutes != 3 {
		t.Errorf("expected delay 3, got %d", trains[0].DelayMinutes)
	}
	if trains[1].DelayMinutes != 0 {
		t.Errorf("expected delay 0, got %d", trains[1].DelayMinutes)
	}
}

// Test 9: limit is enforced
func TestGetNextTrains_LimitEnforced(t *testing.T) {
	visits := []prim.StopVisit{
		futureVisit(5, 0, "A", "", ""),
		futureVisit(10, 0, "B", "", ""),
		futureVisit(15, 0, "C", "", ""),
		futureVisit(20, 0, "D", "", ""),
	}
	svc := newService(visits, nil)
	trains, err := svc.GetNextTrains(context.Background(), "s", "l", "", 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(trains) != 2 {
		t.Fatalf("expected 2 trains (limit=2), got %d", len(trains))
	}
}
