package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Query struct {
	Stop      string `form:"stop" binding:"required"`
	Line      string `form:"line" binding:"required"`
	Direction string `form:"direction"`
	Limit     int    `form:"limit"`
}

func (h *handler) GetNextTrainHandler(c *gin.Context) {
	var q Query
	if err := c.ShouldBindQuery(&q); err != nil {
		h.log.Error("api: invalid query parameters", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.service.GetNextTrains(c.Request.Context(), q.Stop, q.Line, q.Direction, q.Limit)
	if err != nil {
		h.log.Error("api: service error", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	if len(result) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "no trains found"})
		return
	}
	c.JSON(http.StatusOK, result)
}
