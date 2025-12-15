package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jossbnd/trainwatch/backend/internal/api"
	"github.com/jossbnd/trainwatch/backend/internal/config"
	"github.com/jossbnd/trainwatch/backend/internal/prim"
	"github.com/jossbnd/trainwatch/backend/internal/service"
)

func main() {
	// Load config
	c, err := config.Load()
	if err != nil {
		log.Fatal("failed to load config: ", err)
	}
	log.Printf("config loaded: env=%s, gin_mode=%s, port=%s", c.Env, c.GinMode, c.Port)

	// Set gin mode from config
	gin.SetMode(c.GinMode)

	// Initialize prim client
	primClient, err := prim.New(c.PrimBaseURL, c.PrimAPIKey)
	if err != nil {
		log.Fatal("failed to initialize prim client: ", err)
	}
	log.Printf("prim client initialized with baseURL=%s", c.PrimBaseURL)

	// Initialize service
	svc := service.New(service.Input{
		PrimClient: primClient,
	})
	log.Printf("service initialized")

	// Initialize API
	r := api.New(api.Input{
		Service: svc,
	})

	addr := fmt.Sprintf(":%s", c.Port)
	log.Printf("starting server on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("server exited with error: %v", err)
	}
}
