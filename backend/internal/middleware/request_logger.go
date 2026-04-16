package middleware

import (
	"time"

	"github.com/gin-gonic/gin"

	"github.com/jossbnd/trainwatch/backend/internal/logger"
)

func RequestLogger(log *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		ctx := c.Request.Context()
		log.Infoc(ctx, "handling request",
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"query", c.Request.URL.RawQuery,
		)
		c.Next()

		status := c.Writer.Status()
		attrs := []any{
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"query", c.Request.URL.RawQuery,
			"status", status,
			"latency_ms", time.Since(start).Milliseconds(),
			"client_ip", c.ClientIP(),
			"user_agent", c.Request.UserAgent(),
		}

		switch {
		case status >= 500:
			log.Errorc(ctx, "request", attrs...)
		case status >= 400:
			log.Warnc(ctx, "request", attrs...)
		default:
			log.Infoc(ctx, "request", attrs...)
		}
	}
}
