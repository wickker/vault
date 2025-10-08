package services

import (
	"context"
	"database/sql"
	"errors"
	"vault/db/sqlc"
	"vault/middleware"
	"vault/openapi"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type VaultService struct {
	queries       *sqlc.Queries
	dbPool        *sql.DB
	encryptionKey string
}

type ContextValues struct {
	Logger     zerolog.Logger
	User       *clerk.User
	GinContext context.Context
}

var _ openapi.StrictServerInterface = (*VaultService)(nil)

func NewVaultService(queries *sqlc.Queries, dbPool *sql.DB, encryptionKey string) *VaultService {
	return &VaultService{
		queries:       queries,
		dbPool:        dbPool,
		encryptionKey: encryptionKey,
	}
}

func (v *VaultService) getContextValues(c context.Context) (ContextValues, error) {
	var err error

	logger, ok := c.Value(middleware.ContextKeys.Logger).(zerolog.Logger)
	if !ok {
		log.Warn().Msg("Unable to get logger from Gin context, defaulting to global logger.")
		logger = log.Logger
	}

	user, ok := c.Value(middleware.ContextKeys.User).(*clerk.User)
	if !ok {
		err = errors.New("failed to parse clerk user from context")
		logger.Err(err).Msgf("Unable to parse clerk user from context [User: %+v].", user)
	}

	ginCtx, ok := c.Value("ginContext").(context.Context)
	if !ok {
		logger.Warn().Msg("Unable to extract gin context from context.")
		ginCtx = c
	}

	return ContextValues{
		GinContext: ginCtx,
		User:       user,
		Logger:     logger,
	}, err
}
