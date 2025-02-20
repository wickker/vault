package services

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5/pgtype"
	"vault/db/sqlc"
	"vault/openapi"
)

// (GET /categories)
func (v *VaultService) GetCategories(ctx context.Context, _ openapi.GetCategoriesRequestObject) (openapi.GetCategoriesResponseObject, error) {
	logger := v.getLogger(ctx)
	user, err := v.getUser(ctx)
	if err != nil {
		return openapi.GetCategories4XXJSONResponse{Body: openapi.Error{
			Message: err.Error(),
		}, StatusCode: 401}, nil
	}

	categories, err := v.queries.ListCategoriesByUser(ctx, user.ID)
	if err != nil {
		logger.Err(err).Msgf("Unable to list categories by user [UserID: %s].", user.ID)
		return openapi.GetCategories5XXJSONResponse{Body: openapi.Error{
			Message: err.Error(),
		}, StatusCode: 500}, nil
	}

	result := openapi.GetCategories200JSONResponse{}
	for _, category := range categories {
		result = append(result, openapi.Category{
			Id:    category.ID,
			Name:  category.Name,
			Color: category.Color,
		})
	}

	return result, nil
}

// (POST /categories)
func (v *VaultService) CreateCategory(ctx context.Context, request openapi.CreateCategoryRequestObject) (openapi.CreateCategoryResponseObject, error) {
	logger := v.getLogger(ctx)
	user, err := v.getUser(ctx)
	if err != nil {
		return openapi.CreateCategory4XXJSONResponse{Body: openapi.Error{
			Message: err.Error(),
		}, StatusCode: 401}, nil
	}

	category, err := v.queries.CreateCategory(ctx, sqlc.CreateCategoryParams{
		Name:        request.Body.Name,
		Color:       request.Body.Color,
		ClerkUserID: user.ID,
	})
	if err != nil {
		logger.Err(err).Msgf("Unable to create category [Request: %+v][ClerkUserID: %s].", request.Body, user.ID)
		return openapi.CreateCategory5XXJSONResponse{Body: openapi.Error{
			Message: err.Error(),
		}, StatusCode: 500}, nil
	}

	return openapi.CreateCategory201JSONResponse{
		Name:  category.Name,
		Id:    category.ID,
		Color: category.Color,
	}, nil
}

// (DELETE /categories/{categoryId})
func (v *VaultService) DeleteCategory(ctx context.Context, request openapi.DeleteCategoryRequestObject) (openapi.DeleteCategoryResponseObject, error) {
	logger := v.getLogger(ctx)
	user, err := v.getUser(ctx)
	if err != nil {
		return openapi.DeleteCategory4XXJSONResponse{Body: openapi.Error{
			Message: err.Error(),
		}, StatusCode: 401}, nil
	}

	itemsWithCategory, err := v.queries.ListItemsByCategory(ctx, sqlc.ListItemsByCategoryParams{
		CategoryID: pgtype.Int4{
			Int32: request.CategoryId,
			Valid: true,
		},
		ClerkUserID: user.ID,
	})
	if err != nil {
		logger.Err(err).Msgf("Unable to get items by category [CategoryID: %v][ClerkUserID: %s].", request.CategoryId, user.ID)
		return openapi.DeleteCategory5XXJSONResponse{Body: openapi.Error{
			Message: err.Error(),
		}, StatusCode: 500}, nil
	}
	if len(itemsWithCategory) > 0 {
		err := errors.New("category still has items associated with it")
		logger.Err(err).Msgf("Unable to delete category [CategoryID: %v]", request.CategoryId)
		return openapi.DeleteCategory5XXJSONResponse{Body: openapi.Error{
			Message: err.Error(),
		}, StatusCode: 500}, nil
	}

	category, err := v.queries.DeleteCategory(ctx, sqlc.DeleteCategoryParams{
		ID:          request.CategoryId,
		ClerkUserID: user.ID,
	})
	if err != nil {
		logger.Err(err).Msgf("Unable to delete category [CategoryID: %v][ClerkUserID: %s].", request.CategoryId, user.ID)
		return openapi.DeleteCategory5XXJSONResponse{Body: openapi.Error{
			Message: err.Error(),
		}, StatusCode: 500}, nil
	}
	if category.ID == 0 {
		logger.Err(err).Msgf("Unable to find category to delete [CategoryID: %v][ClerkUserID: %s].", request.CategoryId, user.ID)
		return openapi.DeleteCategory5XXJSONResponse{Body: openapi.Error{
			Message: err.Error(),
		}, StatusCode: 500}, nil
	}

	return openapi.DeleteCategory204Response{}, nil
}

// (PUT /categories/{categoryId})
func (v *VaultService) UpdateCategory(ctx context.Context, request openapi.UpdateCategoryRequestObject) (openapi.UpdateCategoryResponseObject, error) {
	logger := v.getLogger(ctx)
	user, err := v.getUser(ctx)
	if err != nil {
		return openapi.UpdateCategory4XXJSONResponse{Body: openapi.Error{
			Message: err.Error(),
		}, StatusCode: 401}, nil
	}

	category, err := v.queries.UpdateCategory(ctx, sqlc.UpdateCategoryParams{
		ID:          request.CategoryId,
		Name:        request.Body.Name,
		Color:       request.Body.Color,
		ClerkUserID: user.ID,
	})
	if err != nil {
		logger.Err(err).Msgf("Unable to update category [Request: %+v][ClerkUserID: %s].", request, user.ID)
		return openapi.UpdateCategory5XXJSONResponse{Body: openapi.Error{
			Message: err.Error(),
		}, StatusCode: 500}, nil
	}
	if category.ID == 0 {
		logger.Err(err).Msgf("Unable to find category to update [Request: %+v][ClerkUserID: %s].", request, user.ID)
		return openapi.UpdateCategory5XXJSONResponse{Body: openapi.Error{
			Message: "Item not found",
		}, StatusCode: 500}, nil
	}

	return openapi.UpdateCategory200JSONResponse{
		Id:    category.ID,
		Name:  category.Name,
		Color: category.Color,
	}, nil
}
