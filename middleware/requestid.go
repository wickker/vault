package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.URL.Path == "/" {
			return
		}
		requestId := c.GetHeader("X-Request-Id")
		if requestId == "" {
			requestId = uuid.Must(uuid.NewV7()).String()
		}
		c.Set("requestId", requestId)
		c.Header("X-Request-Id", requestId)

		logger := log.With().Str("request_id", requestId).Str("path", c.Request.URL.Path).Logger()
		c.Set("logger", logger)

		c.Next()
	}
}
