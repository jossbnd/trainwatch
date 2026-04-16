package api

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/jossbnd/trainwatch/backend/internal/middleware"
	"github.com/jossbnd/trainwatch/backend/internal/model"
	"github.com/jossbnd/trainwatch/backend/internal/service"
)

type departuresResponse struct {
	Departures []model.Departure `json:"departures"`
}

type Query struct {
	StopRef   string `form:"stop_ref" binding:"required"`
	LineRef   string `form:"line_ref" binding:"required"`
	Direction string `form:"direction"`
	Limit     int    `form:"limit"`
}

func (h *handler) getDeparturesHandler(c *gin.Context) {
	ctx := c.Request.Context()

	var q Query
	if err := c.ShouldBindQuery(&q); err != nil {
		h.log.Errorc(ctx, "api: invalid query parameters", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.service.GetDepartures(ctx, q.StopRef, q.LineRef, q.Direction, q.Limit)
	if err != nil {
		if errors.Is(err, service.ErrInvalidRequest) {
			h.log.Warnc(ctx, "api: invalid request", "error", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid stop or line reference"})
			return
		}
		h.log.Errorc(ctx, "api: service error", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error", "request_id": c.GetString(middleware.RequestIDKey)})
		return
	}

	if len(result) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "no trains found"})
		return
	}
	c.JSON(http.StatusOK, departuresResponse{Departures: result})
}
