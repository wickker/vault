package middleware

import (
	"github.com/DataDog/dd-trace-go/v2/ddtrace/tracer"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

var ContextKeys = struct {
	Logger     string
	RequestID  string
	GinContext string
	User       string
}{
	Logger:     "logger",
	RequestID:  "requestId",
	GinContext: "ginContext",
	User:       "user",
}

func Trace() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.URL.Path == "/" {
			return
		}
		requestId := c.GetHeader("X-Request-Id")
		if requestId == "" {
			requestId = uuid.Must(uuid.NewV7()).String()
		}
		c.Header("X-Request-Id", requestId)

		t, _ := tracer.SpanFromContext(c.Request.Context())
		logger := log.With().
			Str("path", c.Request.URL.Path).
			Str("trace_id", t.Context().TraceID()).
			Uint64("span_id", t.Context().SpanID()).
			Str("request_id", requestId).
			Logger()

		c.Set(ContextKeys.RequestID, requestId)
		c.Set(ContextKeys.Logger, logger)
		c.Set(ContextKeys.GinContext, c.Request.Context())

		c.Next()
	}
}
