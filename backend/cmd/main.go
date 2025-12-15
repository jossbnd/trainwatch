package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jossbnd/trainwatch/backend/internal/api"
	"github.com/jossbnd/trainwatch/backend/internal/config"
	"github.com/jossbnd/trainwatch/backend/internal/logger"
	"github.com/jossbnd/trainwatch/backend/internal/prim"
	"github.com/jossbnd/trainwatch/backend/internal/service"
)

func main() {
	// Load config
	c, err := config.Load()
	if err != nil {
		logger.New(logger.Input{
			Level: "error",
		}).Error("failed to load config", "error", err)
		os.Exit(1)
	}

	// Set logger
	logger := logger.New(logger.Input{
		Level: c.LogLevel,
	})
	logger.Info(fmt.Sprintf("logger initialized with level=%s", c.LogLevel))

	// Set gin mode from config
	gin.SetMode(c.GinMode)

	// Initialize prim client
	primClient, err := prim.New(c.PrimBaseURL, c.PrimAPIKey)
	if err != nil {
		logger.Error("failed to initialize prim client", "error", err)
		os.Exit(1)
	}
	logger.Info(fmt.Sprintf("prim client initialized with baseURL=%s", c.PrimBaseURL))

	// Initialize service
	svc := service.New(service.Input{
		Logger:     logger,
		PrimClient: primClient,
	})
	logger.Info("service initialized")

	// Initialize API
	r := api.New(api.Input{
		Logger:  logger,
		Service: svc,
	})

	addr := fmt.Sprintf(":%s", c.Port)
	logger.Info(fmt.Sprintf("starting server on %s", addr))
	if err := r.Run(addr); err != nil {
		logger.Error("server exited with error", "error", err)
		os.Exit(1)
	}
}
