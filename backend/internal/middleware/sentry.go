package middleware

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"github.com/jossbnd/trainwatch/backend/internal/sentry"
)

func SentryCapture() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if c.Writer.Status() < 500 {
			return
		}

		ctx := c.Request.Context()
		tags := map[string]string{"request_id": c.GetString(RequestIDKey)}
		if len(c.Errors) > 0 {
			for _, ginErr := range c.Errors {
				sentry.CaptureException(ctx, ginErr.Err, tags)
			}
		} else {
			sentry.CaptureException(ctx, fmt.Errorf("HTTP %d on %s %s", c.Writer.Status(), c.Request.Method, c.Request.URL.Path), tags)
		}
	}
}
