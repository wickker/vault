package main

import (
	"context"
	"errors"
	"github.com/caarlos0/env/v11"
	"github.com/clerk/clerk-sdk-go/v2"
	clerkhttp "github.com/clerk/clerk-sdk-go/v2/http"
	"github.com/clerk/clerk-sdk-go/v2/user"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"vault/config"
	"vault/openapi"
	"vault/services"
)

func main() {
	envCfg := loadEnv()

	clerk.SetKey(envCfg.ClerkSecretKey)

	router := setupGin()

	vaultService := services.NewVaultService()
	vaultHandler := openapi.NewStrictHandler(vaultService, nil)
	openapi.RegisterHandlersWithOptions(router, vaultHandler, openapi.GinServerOptions{
		ErrorHandler: errorHandler,
	})

	server := &http.Server{
		Addr:    ":9000",
		Handler: router,
	}

	// initializing the server in a goroutine so that it does not block graceful shutdown
	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Err(err).Msg("Unable to initialise server.")
		}
	}()

	gracefulShutdown(server)
}

func hello() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		claims, ok := clerk.SessionClaimsFromContext(ctx)
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"access": "unauthorized"}`))
			return
		}

		usr, err := user.Get(ctx, claims.Subject)
		if err != nil {
			panic(err)
		}
		if usr == nil {
			w.Write([]byte("User does not exist"))
			return
		}

		w.Write([]byte("Hello " + *usr.FirstName))
	}
}

func setupGin() *gin.Engine {
	r := gin.Default()

	r.Use(CORS())
	r.Use(gin.WrapH(clerkhttp.WithHeaderAuthorization()(hello())))

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, "Vault is up!")
	})

	return r
}

func gracefulShutdown(server *http.Server) {
	channel := make(chan os.Signal)
	signal.Notify(channel, syscall.SIGINT, syscall.SIGTERM)
	<-channel
	log.Info().Msg("Shutting down server.")

	// the context is used to inform the server it has 10 seconds to finish the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal().Err(err).Msg("Unable to shutdown server.")
	}
}

func errorHandler(c *gin.Context, err error, statusCode int) {
	log.Err(err).Msgf("Error occurred on request %s %s", c.Request.Method, c.Request.URL.Path)

	c.JSON(statusCode, openapi.Error{
		Message: err.Error(),
	})
}

func loadEnv() config.EnvConfig {
	if err := godotenv.Load(); err != nil {
		log.Warn().Msg("Unable to read from .env file.")
	}
	var envCfg config.EnvConfig
	if err := env.Parse(&envCfg); err != nil {
		log.Err(err).Msg("Unable to parse environment variables to struct.")
	}
	return envCfg
}

func CORS() gin.HandlerFunc {
	// creds of otterlite
	corsConfig := cors.Config{
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
			"X-SC-Organization-ID",
		},
		MaxAge: 12 * time.Hour,
	}

	//if cfg.IsDev() {
	corsConfig.AllowOriginFunc = func(origin string) bool { return true }
	//} else {
	//	allowedOrigins := strings.Split(cfg.FrontendOrigin, ",")
	//	corsConfig.AllowOrigins = allowedOrigins
	//}

	return cors.New(corsConfig)
}
