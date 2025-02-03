package services

import (
	"context"
	"vault/db/sqlc"
	"vault/openapi"
)

// (GET /items)
func (v *VaultService) GetItems(ctx context.Context, _ openapi.GetItemsRequestObject) (openapi.GetItemsResponseObject, error) {
	return openapi.GetItems200JSONResponse{}, nil
}

// (POST /items)
func (v *VaultService) CreateItem(ctx context.Context, request openapi.CreateItemRequestObject) (openapi.CreateItemResponseObject, error) {
	user, err := v.getUser(ctx)
	if err != nil {
		return openapi.CreateItem4XXJSONResponse{Body: openapi.Error{
			Message: err.Error(),
		}, StatusCode: 401}, nil
	}

	item, err := v.queries.CreateItem(ctx, sqlc.CreateItemParams{
		Name:        request.Body.Name,
		ClerkUserID: user.ID,
	})
	if err != nil {
		return openapi.CreateItem5XXJSONResponse{Body: openapi.Error{
			Message: err.Error(),
		}, StatusCode: 500}, nil
	}

	return openapi.CreateItem201JSONResponse{
		Name: item.Name,
		Id:   item.ID,
	}, nil
}

// (DELETE /items/{itemId})
func (v *VaultService) DeleteItem(ctx context.Context, request openapi.DeleteItemRequestObject) (openapi.DeleteItemResponseObject, error) {
	return openapi.DeleteItem204Response{}, nil
}

// (PUT /items/{itemId})
func (v *VaultService) UpdateItem(ctx context.Context, request openapi.UpdateItemRequestObject) (openapi.UpdateItemResponseObject, error) {
	return openapi.UpdateItem200JSONResponse{}, nil
}
