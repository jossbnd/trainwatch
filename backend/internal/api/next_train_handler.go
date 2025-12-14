package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jossbnd/trainwatch/backend/internal/service"
)

// GET /next-train?type=metro&line=1&station=Chatelet&direction=A
func GetNextTrainHandler(c *gin.Context) {
	transportType := c.Query("type")
	line := c.Query("line")
	station := c.Query("station")
	direction := c.Query("direction")

	if transportType == "" || line == "" || station == "" || direction == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing required query params: type, line, station, direction"})
		return
	}

	result, err := service.GetNextTrain(transportType, line, station, direction)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}
