package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Query struct {
	Stop      string `form:"stop" binding:"required"`
	Line      string `form:"line" binding:"required"`
	Direction string `form:"direction"`
}

func (h *handler) GetNextTrainHandler(c *gin.Context) {
	var q Query
	if err := c.ShouldBindQuery(&q); err != nil {
		h.logger.Error("api: invalid query parameters", "error", err)
		c.JSON(http.StatusBadRequest, err)
		return
	}

	result, err := h.service.GetNextTrains(q.Stop, q.Line, q.Direction)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	if len(result) == 0 {
		c.JSON(http.StatusNotFound, err)
		return
	}
	c.JSON(http.StatusOK, result)
}
