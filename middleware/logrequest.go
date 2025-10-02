package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

func LogRequest() gin.HandlerFunc {
	return func(c *gin.Context) {
		//start := time.Now()
		//raw := c.Request.URL.RawQuery

		c.Next()

		//latency := time.Since(start)
		//status := c.Writer.Status()

		//Int("status", status).
		//Dur("latency", latency).
		//Str("client_ip", c.ClientIP())

		logger := zerolog.Ctx(c.Request.Context())
		event := logger.Info()

		if len(c.Errors) > 0 {
			event = event.Str("errors", c.Errors.String())
		}

		event.Msg("Request completed.")
	}
}
