package sentry

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/getsentry/sentry-go"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
)

type Input struct {
	Enabled     bool
	DSN         string
	Environment string
	EnableLogs  bool
}

func Init(input Input) (func(time.Duration) bool, error) {
	err := sentry.Init(sentry.ClientOptions{
		Dsn:              input.DSN,
		Environment:      input.Environment,
		EnableLogs:       input.EnableLogs,
		EnableTracing:    true,
		TracesSampleRate: 1.0,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize sentry: %w", err)
	}

	return sentry.Flush, nil
}

// CaptureException sends an error to Sentry with the given tags.
func CaptureException(ctx context.Context, err error, tags map[string]string) {
	hub := sentry.GetHubFromContext(ctx)
	if hub == nil {
		return
	}
	for k, v := range tags {
		hub.Scope().SetTag(k, v)
	}
	hub.CaptureException(err)
}

// GinMiddleware returns the sentrygin middleware for use in Gin routers.
func GinMiddleware() gin.HandlerFunc {
	return sentrygin.New(sentrygin.Options{Repanic: true})
}

const MetricPrimCreditsRemainingDay = "prim.credits_remaining_day"

// SendGauge emits a gauge metric to Sentry.
func SendGauge(ctx context.Context, key string, value float64) {
	sentry.NewMeter(ctx).Gauge(key, value)
}

// EmitLog sends a log entry to Sentry at the appropriate level.
func EmitLog(ctx context.Context, level slog.Level, msg string, attrs map[string]string) {
	sl := sentry.NewLogger(ctx)
	entry := toSentryEntry(sl, level).WithCtx(ctx)
	for k, v := range attrs {
		entry = entry.String(k, v)
	}
	entry.Emit(msg)
}

func toSentryEntry(sl sentry.Logger, l slog.Level) sentry.LogEntry {
	switch {
	case l >= slog.LevelError:
		return sl.Error()
	case l >= slog.LevelWarn:
		return sl.Warn()
	case l >= slog.LevelInfo:
		return sl.Info()
	default:
		return sl.Debug()
	}
}
