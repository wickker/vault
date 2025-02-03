package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/caarlos0/env/v11"
	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	ginmiddleware "github.com/oapi-codegen/gin-middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"net/http"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"
	"time"
	"vault/config"
	"vault/db/sqlc"
	"vault/middleware"
	"vault/openapi"
	"vault/services"
)

func main() {
	setupLogger()

	envCfg := loadEnv()

	pool, err := pgxpool.New(context.Background(), envCfg.DatabaseURL)
	if err != nil {
		log.Fatal().Err(err).Msg("Unable to connect to database.")
	}
	defer pool.Close()
	queries := sqlc.New(pool)

	clerk.SetKey(envCfg.ClerkSecretKey)

	router := setupGin(envCfg)

	protectedRoutes := router.Group("protected", ginmiddleware.OapiRequestValidator(getSwagger()))
	vaultService := services.NewVaultService(queries)
	vaultHandler := openapi.NewStrictHandler(vaultService, nil)
	openapi.RegisterHandlersWithOptions(protectedRoutes, vaultHandler, openapi.GinServerOptions{
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

func setupLogger() {
	log.Logger = zerolog.New(os.Stdout).With().Timestamp().Caller().Logger()
	zerolog.ErrorStackMarshaler = func(err error) interface{} {
		return string(debug.Stack())
	}
}

func setupGin(envCfg config.EnvConfig) *gin.Engine {
	r := gin.Default()
	r.Use(middleware.Cors(envCfg))
	r.Use(middleware.RequestID())
	r.Use(middleware.Auth())

	//r.Use(f)

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

func getSwagger() *openapi3.T {
	spec, err := os.ReadFile("openapi/openapi.yaml")
	if err != nil {
		log.Err(err).Msg("Unable to open and read openapi.yaml")
		return nil
	}

	swagger, err := openapi3.NewLoader().LoadFromData(spec)
	if err != nil {
		log.Err(err).Msg("Unable to load Swagger")
		return nil
	}

	updatedPaths := &openapi3.Paths{}
	for k, v := range swagger.Paths.Map() {
		updatedPaths.Set(fmt.Sprintf("/protected%s", k), v)
	}
	swagger.Paths = updatedPaths
	return swagger
}
