package services

import (
	"context"
	"github.com/jackc/pgx/v5/pgtype"
	"vault/db/sqlc"
	"vault/openapi"
	"vault/utils"
)

// (GET /items)
func (v *VaultService) GetItems(ctx context.Context, request openapi.GetItemsRequestObject) (openapi.GetItemsResponseObject, error) {
	logger := v.getLogger(ctx)
	user, err := v.getUser(ctx)
	if err != nil {
		return openapi.GetItems4XXJSONResponse{Body: openapi.Error{
			Message: err.Error(),
		}, StatusCode: 401}, nil
	}

	searchPhrase := utils.String{
		Pointer: request.Params.SearchPhrase,
		Like:    true,
	}
	categoryID := utils.Number{Pointer: request.Params.CategoryId}
	items, err := v.queries.ListItemsByUser(ctx, sqlc.ListItemsByUserParams{
		ClerkUserID: user.ID,
		Name:        searchPhrase.ToPgText(),
		OrderBy:     string(request.Params.OrderBy),
		CategoryID:  categoryID.ToPgInt4(),
	})
	if err != nil {
		logger.Err(err).Msgf("Unable to list items by user [UserID: %s][Request: %+v].", user.ID, request.Params)
		return openapi.GetItems5XXJSONResponse{Body: openapi.Error{
			Message: err.Error(),
		}, StatusCode: 500}, nil
	}

	result := openapi.GetItems200JSONResponse{}
	for _, item := range items {
		result = append(result, openapi.Item{
			Id:         item.ID,
			Name:       item.Name,
			CreatedAt:  item.CreatedAt.Time.String(),
			CategoryId: item.CategoryID.Int32,
		})
	}

	return result, nil
}

// (POST /items)
func (v *VaultService) CreateItem(ctx context.Context, request openapi.CreateItemRequestObject) (openapi.CreateItemResponseObject, error) {
	logger := v.getLogger(ctx)
	user, err := v.getUser(ctx)
	if err != nil {
		return openapi.CreateItem4XXJSONResponse{Body: openapi.Error{
			Message: err.Error(),
		}, StatusCode: 401}, nil
	}

	item, err := v.queries.CreateItem(ctx, sqlc.CreateItemParams{
		Name:        request.Body.Name,
		ClerkUserID: user.ID,
		CategoryID: pgtype.Int4{
			Int32: request.Body.CategoryId,
			Valid: true,
		},
	})
	if err != nil {
		logger.Err(err).Msgf("Unable to create item [Request: %+v][ClerkUserID: %s].", request.Body, user.ID)
		return openapi.CreateItem5XXJSONResponse{Body: openapi.Error{
			Message: err.Error(),
		}, StatusCode: 500}, nil
	}

	return openapi.CreateItem201JSONResponse{
		Name:       item.Name,
		Id:         item.ID,
		CreatedAt:  item.CreatedAt.Time.String(),
		CategoryId: item.CategoryID.Int32,
	}, nil
}

// (DELETE /items/{itemId})
func (v *VaultService) DeleteItem(ctx context.Context, request openapi.DeleteItemRequestObject) (openapi.DeleteItemResponseObject, error) {
	logger := v.getLogger(ctx)
	user, err := v.getUser(ctx)
	if err != nil {
		return openapi.DeleteItem4XXJSONResponse{Body: openapi.Error{
			Message: err.Error(),
		}, StatusCode: 401}, nil
	}

	tx, err := v.dbPool.Begin(ctx)
	if err != nil {
		logger.Err(err).Msg("Unable to begin transaction.")
		return openapi.DeleteItem5XXJSONResponse{Body: openapi.Error{
			Message: err.Error(),
		}, StatusCode: 500}, nil
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()
	qtx := v.queries.WithTx(tx)

	item, err := qtx.DeleteItem(ctx, sqlc.DeleteItemParams{
		ID:          request.ItemId,
		ClerkUserID: user.ID,
	})
	if err != nil {
		logger.Err(err).Msgf("Unable to delete item [ID: %v][ClerkUserID: %s].", request.ItemId, user.ID)
		return openapi.DeleteItem5XXJSONResponse{Body: openapi.Error{
			Message: err.Error(),
		}, StatusCode: 500}, nil
	}
	if item.ID == 0 {
		logger.Err(err).Msgf("Unable to find item to delete [ID: %v][ClerkUserID: %s].", request.ItemId, user.ID)
		return openapi.DeleteItem5XXJSONResponse{Body: openapi.Error{
			Message: "Item not found",
		}, StatusCode: 500}, nil
	}

	// delete associated records
	if _, err := qtx.DeleteRecords(ctx, request.ItemId); err != nil {
		logger.Err(err).Msgf("Unable to delete records [ItemID: %v].", request.ItemId)
		return openapi.DeleteItem5XXJSONResponse{Body: openapi.Error{
			Message: err.Error(),
		}, StatusCode: 500}, nil
	}

	if err := tx.Commit(ctx); err != nil {
		logger.Err(err).Msg("Unable to commit transaction.")
		return openapi.DeleteItem5XXJSONResponse{Body: openapi.Error{
			Message: err.Error(),
		}, StatusCode: 500}, nil
	}

	return openapi.DeleteItem204Response{}, nil
}

// (PUT /items/{itemId})
func (v *VaultService) UpdateItem(ctx context.Context, request openapi.UpdateItemRequestObject) (openapi.UpdateItemResponseObject, error) {
	logger := v.getLogger(ctx)
	user, err := v.getUser(ctx)
	if err != nil {
		return openapi.UpdateItem4XXJSONResponse{Body: openapi.Error{
			Message: err.Error(),
		}, StatusCode: 401}, nil
	}

	item, err := v.queries.UpdateItem(ctx, sqlc.UpdateItemParams{
		ID:          request.ItemId,
		Name:        request.Body.Name,
		ClerkUserID: user.ID,
		CategoryID: pgtype.Int4{
			Int32: request.Body.CategoryId,
			Valid: true,
		},
	})
	if err != nil {
		logger.Err(err).Msgf("Unable to update item [Request: %+v][ClerkUserID: %s].", request, user.ID)
		return openapi.UpdateItem5XXJSONResponse{Body: openapi.Error{
			Message: err.Error(),
		}, StatusCode: 500}, nil
	}
	if item.ID == 0 {
		logger.Err(err).Msgf("Unable to find item to update [Request: %+v][ClerkUserID: %s].", request, user.ID)
		return openapi.UpdateItem5XXJSONResponse{Body: openapi.Error{
			Message: "Item not found",
		}, StatusCode: 500}, nil
	}

	return openapi.UpdateItem200JSONResponse{
		Id:         item.ID,
		Name:       item.Name,
		CreatedAt:  item.CreatedAt.Time.String(),
		CategoryId: item.CategoryID.Int32,
	}, nil
}
