package logger

import (
	"log/slog"
	"os"
	"strings"
)

type Logger = *slog.Logger

type Input struct {
	Level string
}

func mapLogLevel(level string) slog.Level {
	switch strings.ToLower(level) {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

func New(input Input) Logger {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: mapLogLevel(input.Level),
	}))

	return logger
}
