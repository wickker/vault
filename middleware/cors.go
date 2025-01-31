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
			"Origin",
			"Content-Length",
			"Content-Type",
			"Access-Control-Allow-Headers",
			"Accept",
			"X-Requested-With",
			"Access-Control-Request-Method",
			"Access-Control-Request-Headers",
			"Access-Control-Allow-Credentials",
			"Access-Control-Allow-Origin",
			"Authorization",
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
