package api

import (
	"github.com/gin-gonic/gin"
	"github.com/jossbnd/trainwatch/backend/internal/service"
)

type Input struct {
	Service service.Service
}

type handler struct {
	service service.Service
}

func New(i Input) *gin.Engine {
	r := gin.Default()

	h := &handler{
		service: i.Service,
	}

	r.GET("/next-train", h.GetNextTrainHandler)
	return r
}
