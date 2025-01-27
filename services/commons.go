package services

import "vault/openapi"

type VaultService struct{}

var _ openapi.StrictServerInterface = (*VaultService)(nil)

func NewVaultService() *VaultService {
	return &VaultService{}
}
