package services

import (
	"context"
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
