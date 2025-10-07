package middleware

import (
	"github.com/DataDog/dd-trace-go/v2/ddtrace/tracer"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type GinContextKey string

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

		t, _ := tracer.SpanFromContext(c.Request.Context())
		logger := log.With().
			Str("path", c.Request.URL.Path).
			Str("trace_id", t.Context().TraceID()).
			Uint64("span_id", t.Context().SpanID()).
			Str("request_id", requestId).
			Logger()
		c.Set("logger", logger)

		c.Set("ginContext", c.Request.Context())

		c.Next()
	}
}
