package api

import (
	"github.com/gin-gonic/gin"
	"github.com/jossbnd/trainwatch/backend/internal/logger"
	"github.com/jossbnd/trainwatch/backend/internal/service"
)

type Input struct {
	Logger  logger.Logger
	Service service.Service
}

type handler struct {
	logger  logger.Logger
	service service.Service
}

func New(i Input) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())

	h := &handler{
		logger:  i.Logger,
		service: i.Service,
	}

	r.GET("/next-train", h.GetNextTrainHandler)
	return r
}
