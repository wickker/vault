package services

import (
	"context"
	"errors"
	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"vault/db/sqlc"
	"vault/openapi"
)

type VaultService struct {
	queries       *sqlc.Queries
	dbPool        *pgxpool.Pool
	encryptionKey string
}

var _ openapi.StrictServerInterface = (*VaultService)(nil)

func NewVaultService(queries *sqlc.Queries, dbPool *pgxpool.Pool, encryptionKey string) *VaultService {
	return &VaultService{
		queries:       queries,
		dbPool:        dbPool,
		encryptionKey: encryptionKey,
	}
}

func (v *VaultService) getLogger(c context.Context) zerolog.Logger {
	logger, ok := c.Value("logger").(zerolog.Logger)
	if !ok {
		log.Warn().Msg("Unable to get logger from Gin context, defaulting to global logger.")
		return log.Logger
	}
	return logger
}

func (v *VaultService) getUser(c context.Context) (*clerk.User, error) {
	logger := v.getLogger(c)
	user, ok := c.Value("user").(*clerk.User)
	if !ok {
		err := errors.New("failed to parse clerk user from context")
		logger.Err(err).Msgf("Unable to parse clerk user from context [User: %+v].", user)
		return nil, err
	}
	return user, nil
}
