package middleware

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"
	"vault/config"
)

func Cors(envCfg config.EnvConfig) gin.HandlerFunc {
	corsCfg := cors.Config{
		AllowMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodDelete,
			http.MethodPatch,
			http.MethodOptions,
			http.MethodHead,
		},
		AllowCredentials: true,
		AllowHeaders: []string{
			"Accept",
			"Authorization",
			"Access-Control-Allow-Credentials",
			"Access-Control-Allow-Headers",
			"Access-Control-Allow-Origin",
			"Access-Control-Request-Headers",
			"Access-Control-Request-Method",
			"Origin",
			"Content-Length",
			"Content-Type",
			"X-Request-Id",
			"X-Requested-With",
		},
		MaxAge: 12 * time.Hour,
	}

	if envCfg.IsDev() {
		corsCfg.AllowOriginFunc = func(origin string) bool { return true }
	} else {
		corsCfg.AllowOrigins = strings.Split(envCfg.FrontendOrigins, ",")
	}

	return cors.New(corsCfg)
}
