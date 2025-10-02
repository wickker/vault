package middleware

import (
	"github.com/DataDog/dd-trace-go/v2/ddtrace/tracer"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

func SetLogTrace() gin.HandlerFunc {
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

		span, _ := tracer.SpanFromContext(c.Request.Context())

		logger := log.With().
			Str("path", c.Request.URL.Path).
			Str("trace_id", span.Context().TraceID()).
			Uint64("span_id", span.Context().SpanID()).
			Str("request_id", requestId).
			Logger()

		ctx := logger.WithContext(c.Request.Context())
		c.Request = c.Request.WithContext(ctx)

		// TODO: Find a way to remove this
		c.Set("logger", logger)

		c.Next()
	}
}
