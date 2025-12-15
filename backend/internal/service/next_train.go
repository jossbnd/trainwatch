package service

import (
	"log"
	"sort"
	"strings"
	"time"

	"github.com/jossbnd/trainwatch/backend/internal/api/dto"
)

// GetNextTrains fetches the next train information for given stop, line and direction.
func (s *service) GetNextTrains(stop, line, direction string) ([]dto.NextTrain, error) {
	visits, err := s.primClient.FetchStopVisits(stop, line)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	var out []dto.NextTrain
	for _, visit := range visits {
		mvj := visit.MonitoredVehicleJourney

		// direction matching (do case-insensitive contains match against
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

		// choose expected time
		var tstr string
		if mvj.MonitoredCall.ExpectedDepartureTime != "" {
			tstr = mvj.MonitoredCall.ExpectedDepartureTime
		} else if mvj.MonitoredCall.ExpectedArrivalTime != "" {
			tstr = mvj.MonitoredCall.ExpectedArrivalTime
		}
		if tstr == "" {
			continue
		}
		ts, err := time.Parse(time.RFC3339, tstr)
		if err != nil {
			ts, err = time.Parse("2006-01-02T15:04:05Z07:00", tstr)
			if err != nil {
				continue
			}
		}
		if ts.Before(now.Add(-1 * time.Minute)) {
			continue
		}

		status := mvj.MonitoredCall.DepartureStatus
		if status == "" {
			status = mvj.MonitoredCall.ArrivalStatus
		}

		out = append(out, dto.NextTrain{EstimatedAt: ts, Status: status})
	}
	// sort ascending
	sort.Slice(out, func(i, j int) bool { return out[i].EstimatedAt.Before(out[j].EstimatedAt) })

	log.Printf("service: returning %d upcoming trains for stop=%s line=%s direction=%s", len(out), stop, line, direction)
	return out, nil
}
