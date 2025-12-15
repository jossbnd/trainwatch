package service

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/jossbnd/trainwatch/backend/internal/api/dto"
)

// GetNextTrains fetches the next departing trains for the given stop, line and
// direction. If direction is empty, returns all directions.
func (s *service) GetNextTrains(stop, line, direction string) ([]dto.NextTrain, error) {
	s.logger.Debug(fmt.Sprintf("GetNextTrains called stop=%s line=%s direction=%s", stop, line, direction))

	// Fetch next stop visits from prim
	visits, err := s.primClient.FetchStopVisits(stop, line)
	if err != nil {
		s.logger.Error(fmt.Sprintf("service: failed to fetch stop visits stop=%s line=%s", stop, line), "error", err)
		return nil, err
	}
	s.logger.Debug(fmt.Sprintf("fetched stop visits count=%d stop=%s line=%s", len(visits), stop, line))

	// Process visits
	now := time.Now()
	var trains []dto.NextTrain
	for _, visit := range visits {
		mvj := visit.MonitoredVehicleJourney

		// Match direction (case-insensitive contains match against
		// DirectionRef, DirectionName, DestinationName). If direction is
		// empty, accept all.
		df := strings.TrimSpace(direction)
		matched := false
		if df == "" {
			matched = true
		} else {
			if mvj.DirectionRef.Value != "" && strings.EqualFold(mvj.DirectionRef.Value, df) {
				matched = true
			}
			if !matched {
				for _, dn := range mvj.DirectionName {
					if strings.Contains(strings.ToLower(dn.Value), strings.ToLower(df)) {
						matched = true
						break
					}
				}
			}
			if !matched {
				for _, dest := range mvj.DestinationName {
					if strings.Contains(strings.ToLower(dest.Value), strings.ToLower(df)) {
						matched = true
						break
					}
				}
			}
		}
		if !matched {
			continue
		}

		// Choose expected departure time if available
		var estimatedAt time.Time
		if mvj.MonitoredCall.ExpectedDepartureTime != nil {
			estimatedAt = *mvj.MonitoredCall.ExpectedDepartureTime
		} else if mvj.MonitoredCall.AimedDepartureTime != nil {
			estimatedAt = *mvj.MonitoredCall.AimedDepartureTime
		}
		if estimatedAt.IsZero() {
			continue
		}

		// Skip past departures (more than 1 minute ago)
		if estimatedAt.Before(now.Add(-1 * time.Minute)) {
			continue
		}

		status := mvj.MonitoredCall.DepartureStatus

		trains = append(trains, dto.NextTrain{EstimatedAt: estimatedAt, Status: status})
	}

	// Sort ascending
	sort.Slice(trains, func(i, j int) bool { return trains[i].EstimatedAt.Before(trains[j].EstimatedAt) })

	s.logger.Info(fmt.Sprintf("service: returning upcoming trains count=%d stop=%s line=%s direction=%s",
		len(trains),
		stop,
		line,
		direction,
	))
	return trains, nil
}
