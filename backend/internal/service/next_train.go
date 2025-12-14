package service

import (
	"github.com/jossbnd/trainwatch/backend/internal/api/dto"
	"github.com/jossbnd/trainwatch/backend/internal/prim"
)

// GetNextTrain fetches the next train information for given params.
func GetNextTrain(transportType, line, station, direction string) (dto.NextTrain, error) {
	trains, err := prim.ListNextTrains(transportType, line, station, direction)
	if err != nil {
		return dto.NextTrain{}, err
	}

	if len(trains) == 0 {
		return dto.NextTrain{}, nil
	}

	// Map prim -> domain -> dto via mapper functions.
	return primToDTO(trains[0]), nil
}
