package service

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/jossbnd/trainwatch/backend/internal/model"
	"github.com/jossbnd/trainwatch/backend/internal/prim"
)

const defaultLimit = 5

// GetNextTrains fetches the next departing trains for the given stop, line and
// direction. If direction is empty, returns all directions. At most limit
// results are returned; if limit <= 0 the default of 5 is used.
func (s *service) GetNextTrains(ctx context.Context, stop, line, direction string, limit int) ([]model.NextTrain, error) {
	s.log.Debug(fmt.Sprintf("GetNextTrains called stop=%s line=%s direction=%s limit=%d", stop, line, direction, limit))

	if limit <= 0 {
		limit = defaultLimit
	}

	// Fetch next visits from prim
	visits, err := s.primClient.FetchStopVisits(ctx, stop, line)
	if err != nil {
		s.log.Error(fmt.Sprintf("service: failed to fetch stop visits stop=%s line=%s", stop, line), "error", err)
		return nil, err
	}
	s.log.Debug(fmt.Sprintf("fetched visits count=%d stop=%s line=%s", len(visits), stop, line))

	// Process visits
	now := time.Now()
	var trains []model.NextTrain
	for _, visit := range visits {

		// Match direction (case-insensitive contains match against
		// DirectionRef, DirectionName, DestinationName). If direction is
		// empty, accept all.
		dirFilter := strings.TrimSpace(direction)
		if !matchDirection(visit, dirFilter) {
			continue
		}

		// Choose expected departure time if available, fall back to aimed.
		var estimatedAt time.Time
		var aimedAt *time.Time
		if visit.Timing.AimedDepartureTime != nil {
			aimedAt = visit.Timing.AimedDepartureTime
		}
		if visit.Timing.ExpectedDepartureTime != nil {
			estimatedAt = *visit.Timing.ExpectedDepartureTime
		} else if visit.Timing.AimedDepartureTime != nil {
			estimatedAt = *visit.Timing.AimedDepartureTime
		}
		if estimatedAt.IsZero() {
			continue
		}

		// Skip past departures (more than 1 minute ago)
		if estimatedAt.Before(now.Add(-1 * time.Minute)) {
			continue
		}

		// Extract destination from first DestinationName entry.
		destination := ""
		if len(visit.DestinationName) > 0 {
			destination = visit.DestinationName[0].Value
		}

		status := visit.Timing.DepartureStatus
		trains = append(trains, model.NewNextTrain(estimatedAt, aimedAt, destination, status))
	}

	// Sort ascending by departure time.
	sort.Slice(trains, func(i, j int) bool { return trains[i].EstimatedAt.Before(trains[j].EstimatedAt) })

	// Apply limit.
	if len(trains) > limit {
		trains = trains[:limit]
	}

	s.log.Info(fmt.Sprintf("service: returning upcoming trains count=%d stop=%s line=%s direction=%s",
		len(trains), stop, line, direction,
	))
	return trains, nil
}

// matchDirection returns true if visit matches the given direction filter.
// An empty filter matches everything.
func matchDirection(visit prim.StopVisit, dirFilter string) bool {
	if dirFilter == "" {
		return true
	}
	if visit.DirectionRef.Value != "" && strings.EqualFold(visit.DirectionRef.Value, dirFilter) {
		return true
	}
	for _, dn := range visit.DirectionName {
		if strings.Contains(strings.ToLower(dn.Value), strings.ToLower(dirFilter)) {
			return true
		}
	}
	for _, dest := range visit.DestinationName {
		if strings.Contains(strings.ToLower(dest.Value), strings.ToLower(dirFilter)) {
			return true
		}
	}
	return false
}
