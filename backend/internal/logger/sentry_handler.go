package logger

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/jossbnd/trainwatch/backend/internal/sentry"
)

type sentryHandler struct {
	inner       slog.Handler
	staticAttrs map[string]string
}

func (h *sentryHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.inner.Enabled(ctx, level)
}

func (h *sentryHandler) Handle(ctx context.Context, r slog.Record) error {
	attrs := make(map[string]string, len(h.staticAttrs))
	for k, v := range h.staticAttrs {
		attrs[k] = v
	}
	r.Attrs(func(a slog.Attr) bool {
		attrs[a.Key] = fmt.Sprintf("%v", a.Value.Any())
		return true
	})
	sentry.EmitLog(ctx, r.Level, r.Message, attrs)

	return h.inner.Handle(ctx, r)
}

func (h *sentryHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	merged := make(map[string]string, len(h.staticAttrs)+len(attrs))
	for k, v := range h.staticAttrs {
		merged[k] = v
	}
	for _, a := range attrs {
		merged[a.Key] = fmt.Sprintf("%v", a.Value.Any())
	}
	return &sentryHandler{inner: h.inner.WithAttrs(attrs), staticAttrs: merged}
}

func (h *sentryHandler) WithGroup(name string) slog.Handler {
	return &sentryHandler{inner: h.inner.WithGroup(name), staticAttrs: h.staticAttrs}
}
