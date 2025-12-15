package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *handler) GetNextTrainHandler(c *gin.Context) {
	stop := c.Query("stop")
	line := c.Query("line")
	direction := ""

	if line == "" || stop == "" {
		fmt.Println("missing required query params: line, stop", line, stop)
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing required query params: line, stop"})
		return
	}

	result, err := h.service.GetNextTrains(stop, line, direction)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(result) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "no upcoming trains found"})
		return
	}
	c.JSON(http.StatusOK, result)
}
