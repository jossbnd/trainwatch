package service

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/jossbnd/trainwatch/backend/internal/model"
	"github.com/jossbnd/trainwatch/backend/internal/prim"
	"github.com/jossbnd/trainwatch/backend/internal/sentry"
)

const defaultLimit = 5

// GetDepartures fetches the next departures for the given stop, line and
// direction. If direction is empty, returns all directions. At most limit
// results are returned; if limit <= 0 the default of 5 is used.
func (s *service) GetDepartures(ctx context.Context, stop, line, direction string, limit int) ([]model.Departure, error) {
	s.log.Debugc(ctx, fmt.Sprintf("GetDepartures called stop=%s line=%s direction=%s limit=%d", stop, line, direction, limit))

	if limit <= 0 {
		limit = defaultLimit
	}

	visits, credits, err := s.primClient.FetchStopVisits(ctx, stop, line)
	if err != nil {
		var respErr *prim.ResponseError
		if errors.As(err, &respErr) && respErr.IsClientError() {
			s.log.Warnc(ctx, fmt.Sprintf("service: prim rejected request with status %d: %s", respErr.StatusCode, respErr.Body),
				"stop", stop,
				"line", line,
			)
			return nil, fmt.Errorf("%w: %s", ErrInvalidRequest, err)
		}
		s.log.Errorc(ctx, fmt.Sprintf("service: prim request failed stop=%s line=%s: %s", stop, line, err))
		return nil, err
	}
	if credits >= 0 {
		sentry.SendGauge(ctx, sentry.MetricPrimCreditsRemainingDay, float64(credits))
	}
	s.log.Debugc(ctx, fmt.Sprintf("fetched visits count=%d stop=%s line=%s credits=%d", len(visits), stop, line, credits))

	now := time.Now()
	var departures []model.Departure
	for _, visit := range visits {
		dirFilter := strings.TrimSpace(direction)
		if !matchDirection(visit, dirFilter) {
			continue
		}

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

		if estimatedAt.Before(now.Add(-1 * time.Minute)) {
			continue
		}

		destination := ""
		if len(visit.DestinationName) > 0 {
			destination = visit.DestinationName[0].Value
		}

		status := visit.Timing.DepartureStatus
		departures = append(departures, model.NewDeparture(estimatedAt, aimedAt, destination, status))
	}

	sort.Slice(departures, func(i, j int) bool { return departures[i].EstimatedAt.Before(departures[j].EstimatedAt) })

	if len(departures) > limit {
		departures = departures[:limit]
	}

	s.log.Infoc(ctx, fmt.Sprintf("service: returning upcoming departures count=%d stop=%s line=%s direction=%s",
		len(departures), stop, line, direction,
	))
	return departures, nil
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
