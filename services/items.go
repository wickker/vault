package services

import (
	"context"
	"vault/openapi"
)

// (GET /items)
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

// (POST /items)
func (v *VaultService) CreateItem(ctx context.Context, request openapi.CreateItemRequestObject) (openapi.CreateItemResponseObject, error) {
	return openapi.CreateItem201JSONResponse{}, nil
}

// (DELETE /items/{itemId})
func (v *VaultService) DeleteItem(ctx context.Context, request openapi.DeleteItemRequestObject) (openapi.DeleteItemResponseObject, error) {
	return openapi.DeleteItem204Response{}, nil
}

// (PUT /items/{itemId})
func (v *VaultService) UpdateItem(ctx context.Context, request openapi.UpdateItemRequestObject) (openapi.UpdateItemResponseObject, error) {
	return openapi.UpdateItem200JSONResponse{}, nil
}
