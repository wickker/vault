package services

import (
	"context"
	"vault/openapi"
)

func (v *VaultService) GetItems(_ context.Context, _ openapi.GetItemsRequestObject) (openapi.GetItemsResponseObject, error) {
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
