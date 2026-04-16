package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/jossbnd/trainwatch/backend/internal/logger"
)

const RequestIDKey = "request_id"

func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := uuid.New().String()
		c.Set(RequestIDKey, id)
		c.Header("X-Request-ID", id)
		c.Request = c.Request.WithContext(logger.ContextWithRequestID(c.Request.Context(), id))
		c.Next()
	}
}
