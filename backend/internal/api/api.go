package api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/jossbnd/trainwatch/backend/internal/logger"
	"github.com/jossbnd/trainwatch/backend/internal/middleware"
	"github.com/jossbnd/trainwatch/backend/internal/service"
)

type Input struct {
	Logger  *logger.Logger
	Service service.Service
	APIKey  string
}

type handler struct {
	log     *logger.Logger
	service service.Service
}

func New(input Input) *gin.Engine {
	r := gin.New()
	r.Use(middleware.RequestID())
	r.Use(middleware.RequestLogger(input.Logger))
	r.Use(gin.Recovery())

	h := &handler{
		log:     input.Logger,
		service: input.Service,
	}

	r.GET("/health", healthHandler)

	auth := r.Group("/")
	auth.Use(middleware.APIKeyAuth(input.APIKey))
	auth.GET("/departures", h.getDeparturesHandler)
	return r
}

func healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
