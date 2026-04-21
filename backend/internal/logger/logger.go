package logger

import (
	"context"
	"io"
	"log/slog"
	"os"
	"strings"
)

type Logger struct {
	l *slog.Logger
}

type Input struct {
	Level        string
	EnableSentry bool
}

// Context keys for request enrichment.
type (
	requestIDKey    struct{}
	requestAttrsKey struct{}
)

type RequestAttrs struct {
	Method    string
	Path      string
	Query     string
	ClientIP  string
	UserAgent string
}

func ContextWithRequestAttrs(ctx context.Context, attrs RequestAttrs) context.Context {
	return context.WithValue(ctx, requestAttrsKey{}, attrs)
}

func ContextWithRequestID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, requestIDKey{}, id)
}

// Standard methods (no context enrichment).

func (lg *Logger) Debug(msg string, args ...any) { lg.l.Debug(msg, args...) }
func (lg *Logger) Info(msg string, args ...any)  { lg.l.Info(msg, args...) }
func (lg *Logger) Warn(msg string, args ...any)  { lg.l.Warn(msg, args...) }
func (lg *Logger) Error(msg string, args ...any) { lg.l.Error(msg, args...) }

// Context-aware methods

func (lg *Logger) Debugc(ctx context.Context, msg string, args ...any) {
	lg.l.DebugContext(ctx, msg, args...)
}

func (lg *Logger) Infoc(ctx context.Context, msg string, args ...any) {
	lg.l.InfoContext(ctx, msg, args...)
}

func (lg *Logger) Warnc(ctx context.Context, msg string, args ...any) {
	lg.l.WarnContext(ctx, msg, args...)
}

func (lg *Logger) Errorc(ctx context.Context, msg string, args ...any) {
	lg.l.ErrorContext(ctx, msg, args...)
}

// contextHandler wraps an slog.Handler and enriches log records with the
// request_id stored in the context.
type contextHandler struct {
	inner slog.Handler
}

func (h *contextHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.inner.Enabled(ctx, level)
}

func (h *contextHandler) Handle(ctx context.Context, r slog.Record) error {
	if id, ok := ctx.Value(requestIDKey{}).(string); ok && id != "" {
		r.AddAttrs(slog.String("request_id", id))
	}
	if a, ok := ctx.Value(requestAttrsKey{}).(RequestAttrs); ok {
		r.AddAttrs(
			slog.String("method", a.Method),
			slog.String("path", a.Path),
			slog.String("query", a.Query),
			slog.String("client_ip", a.ClientIP),
			slog.String("user_agent", a.UserAgent),
		)
	}
	return h.inner.Handle(ctx, r)
}

func (h *contextHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &contextHandler{inner: h.inner.WithAttrs(attrs)}
}

func (h *contextHandler) WithGroup(name string) slog.Handler {
	return &contextHandler{inner: h.inner.WithGroup(name)}
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

func New(input Input) *Logger {
	jsonHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: mapLogLevel(input.Level),
	})
	var inner slog.Handler = jsonHandler
	if input.EnableSentry {
		inner = &sentryHandler{inner: jsonHandler}
	}
	return &Logger{l: slog.New(&contextHandler{inner: inner})}
}

func NewDiscard() *Logger {
	return &Logger{l: slog.New(slog.NewTextHandler(io.Discard, nil))}
}
