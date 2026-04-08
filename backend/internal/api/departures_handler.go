package api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/jossbnd/trainwatch/backend/internal/model"
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
	var q Query
	if err := c.ShouldBindQuery(&q); err != nil {
		h.log.Error("api: invalid query parameters", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.service.GetDepartures(c.Request.Context(), q.StopRef, q.LineRef, q.Direction, q.Limit)
	if err != nil {
		h.log.Error("api: service error", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	if len(result) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "no trains found"})
		return
	}
	c.JSON(http.StatusOK, departuresResponse{Departures: result})
}
