package service

import (
	"github.com/jossbnd/trainwatch/backend/internal/api/dto"
	"github.com/jossbnd/trainwatch/backend/internal/prim"
)

type Input struct {
	PrimClient prim.Client
}

type Service interface {
	// GetNextTrain returns upcoming trains for the given stop, line and
	// direction. If direction is empty, returns all directions.
	GetNextTrains(stop, line, direction string) ([]dto.NextTrain, error)
}

type service struct {
	primClient prim.Client
}

func New(i Input) Service {
	return &service{
		primClient: i.PrimClient,
	}
}
