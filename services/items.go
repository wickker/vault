package services

import (
	"context"
	"vault/openapi"
)

func (v *VaultService) GetItems(c context.Context, _ openapi.GetItemsRequestObject) (openapi.GetItemsResponseObject, error) {
	logger := v.getLogger(c)
	logger.Info().Msg("Test")

	return openapi.GetItems200JSONResponse{
		{
			Id:   1,
			Name: "Hello",
		},
		{
			Id:   2,
			Name: "World",
		},
	}, nil
}
