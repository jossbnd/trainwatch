package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/jossbnd/trainwatch/backend/internal/api"
	"github.com/jossbnd/trainwatch/backend/internal/config"
	"github.com/jossbnd/trainwatch/backend/internal/logger"
	"github.com/jossbnd/trainwatch/backend/internal/prim"
	"github.com/jossbnd/trainwatch/backend/internal/sentry"
	"github.com/jossbnd/trainwatch/backend/internal/service"
)

func main() {
	// Load config
	cfg, err := config.Load()
	if err != nil {
		logger.New(logger.Input{Level: "error"}).Error("failed to load config", "error", err)
		os.Exit(1)
	}

	// Initialize Sentry before the logger so log entries are forwarded from startup.
	if cfg.Sentry.DSN != "" {
		flush, err := sentry.Init(sentry.Input{
			DSN:         cfg.Sentry.DSN,
			Environment: cfg.Env,
			EnableLogs:  cfg.Sentry.EnableLogs,
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "sentry init failed: %v\n", err)
		}
		defer flush(2 * time.Second)
	}

	// Set logger
	log := logger.New(logger.Input{Level: cfg.LogLevel, EnableSentry: cfg.Sentry.DSN != ""})
	log.Info(fmt.Sprintf("logger initialized with level=%s sentry=%t", cfg.LogLevel, cfg.Sentry.DSN != ""))

	// Set gin mode from config
	gin.SetMode(cfg.GinMode)

	// Initialize prim client
	primClient, err := prim.New(cfg.Prim.BaseURL, cfg.Prim.APIKey)
	if err != nil {
		log.Error("failed to initialize prim client", "error", err)
		os.Exit(1)
	}
	log.Info(fmt.Sprintf("prim client initialized with baseURL=%s", cfg.Prim.BaseURL))

	// Initialize service
	svc := service.New(service.Input{
		Logger:     log,
		PrimClient: primClient,
	})
	log.Info("service initialized")

	// Initialize API
	r := api.New(api.Input{
		Logger:  log,
		Service: svc,
		APIKey:  cfg.APIKey,
	})

	addr := fmt.Sprintf(":%s", cfg.Port)
	srv := &http.Server{Addr: addr, Handler: r}

	// Start server in background goroutine.
	go func() {
		log.Info(fmt.Sprintf("application is running on %s", addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("server exited with error", "error", err)
			os.Exit(1)
		}
	}()

	// Wait for termination signal.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info("shutting down server")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Error("server shutdown error", "error", err)
		os.Exit(1)
	}
	log.Info("server stopped")
}
