package services

import (
	"context"
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
