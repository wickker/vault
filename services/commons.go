package services

import (
	"context"
	"errors"
	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"vault/db/sqlc"
	"vault/openapi"
)

type VaultService struct {
	queries *sqlc.Queries
}

var _ openapi.StrictServerInterface = (*VaultService)(nil)

func NewVaultService(queries *sqlc.Queries) *VaultService {
	return &VaultService{
		queries: queries,
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
	user, ok := c.Value("user").(*clerk.User)
	if !ok {
		return nil, errors.New("failed to parse clerk user from context")
	}
	return user, nil
}
