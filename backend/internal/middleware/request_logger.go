package middleware

import (
	"time"

	"github.com/gin-gonic/gin"

	"github.com/jossbnd/trainwatch/backend/internal/logger"
)

func RequestLogger(log *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		ctx := logger.ContextWithRequestAttrs(c.Request.Context(), logger.RequestAttrs{
			Method:    c.Request.Method,
			Path:      c.Request.URL.Path,
			Query:     c.Request.URL.RawQuery,
			ClientIP:  c.ClientIP(),
			UserAgent: c.Request.UserAgent(),
		})
		c.Request = c.Request.WithContext(ctx)
		log.Infoc(ctx, "handling request")
		c.Next()

		status := c.Writer.Status()
		attrs := []any{
			"status", status,
			"latency_ms", time.Since(start).Milliseconds(),
		}

		switch {
		case status >= 500:
			log.Errorc(ctx, "request completed", attrs...)
		case status >= 400:
			log.Warnc(ctx, "request completed", attrs...)
		default:
			log.Infoc(ctx, "request completed", attrs...)
		}
	}
}
