package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jossbnd/trainwatch/backend/internal/logger"
	"github.com/jossbnd/trainwatch/backend/internal/service"
)

type Input struct {
	Logger  logger.Logger
	Service service.Service
}

type handler struct {
	log     logger.Logger
	service service.Service
}

func New(input Input) *gin.Engine {
	r := gin.New()
	r.Use(requestLogger(input.Logger))
	r.Use(gin.Recovery())

	h := &handler{
		log:     input.Logger,
		service: input.Service,
	}

	r.GET("/health", healthHandler)
	r.GET("/next-train", h.GetNextTrainHandler)
	return r
}

func healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func requestLogger(log logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		log.Info("request",
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"status", c.Writer.Status(),
			"latency_ms", time.Since(start).Milliseconds(),
		)
	}
}
