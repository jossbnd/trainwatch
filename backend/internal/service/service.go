package service

import (
	"context"
	"errors"

	"github.com/jossbnd/trainwatch/backend/internal/logger"
	"github.com/jossbnd/trainwatch/backend/internal/model"
	"github.com/jossbnd/trainwatch/backend/internal/prim"
)

var ErrInvalidRequest = errors.New("invalid request")

type Input struct {
	Logger     *logger.Logger
	PrimClient prim.Client
}

type Service interface {
	// GetDepartures returns upcoming departures for the given stop, line and
	// direction. If direction is empty, returns all directions. At most limit
	// results are returned (0 means use default of 5).
	GetDepartures(ctx context.Context, stop, line, direction string, limit int) ([]model.Departure, error)
}

type service struct {
	log        *logger.Logger
	primClient prim.Client
}

func New(input Input) Service {
	return &service{
		log:        input.Logger,
		primClient: input.PrimClient,
	}
}
