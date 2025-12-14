package service

import (
	"github.com/jossbnd/trainwatch/backend/internal/api/dto"
	"github.com/jossbnd/trainwatch/backend/internal/prim"
)

func primToDTO(d prim.NextTrain) dto.NextTrain {
	return dto.NextTrain{
		EstimatedAt: d.EstimatedAt,
		Status:      d.Status,
	}
}
