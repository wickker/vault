package main

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"vault/openapi"
	"vault/services"
)

func main() {
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

func setupGin() *gin.Engine {
	r := gin.Default()

	// TODO: Add middleware

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
