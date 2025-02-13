package services

import (
	"context"
	"vault/db/sqlc"
	"vault/openapi"
)

// (GET /records)
func (v *VaultService) GetRecordsByItem(ctx context.Context, request openapi.GetRecordsByItemRequestObject) (openapi.GetRecordsByItemResponseObject, error) {
	logger := v.getLogger(ctx)
	user, err := v.getUser(ctx)
	if err != nil {
		return openapi.GetRecordsByItem4XXJSONResponse{Body: openapi.Error{
			Message: err.Error(),
		}, StatusCode: 401}, nil
	}

	recordsByItem, err := v.queries.ListRecordsByItem(ctx, sqlc.ListRecordsByItemParams{
		ID:          request.Params.ItemId,
		ClerkUserID: user.ID,
	})
	if err != nil {
		logger.Err(err).Msgf("Unable to list records by item [UserID: %s][ItemID: %v].", user.ID, request.Params.ItemId)
		return openapi.GetRecordsByItem5XXJSONResponse{Body: openapi.Error{
			Message: err.Error(),
		}, StatusCode: 500}, nil
	}

	records := make([]openapi.Record, len(recordsByItem))
	for i, record := range recordsByItem {
		records[i] = openapi.Record{
			Id:    record.RecordID,
			Name:  record.RecordName,
			Value: record.RecordValue,
		}
	}

	if len(records) > 0 {
		return openapi.GetRecordsByItem200JSONResponse{
			Id:      &records[0].Id,
			Name:    &records[0].Name,
			Records: &records,
		}, nil
	}

	return openapi.GetRecordsByItem200JSONResponse{}, nil
}

// (POST /records)
func (v *VaultService) CreateRecord(ctx context.Context, request openapi.CreateRecordRequestObject) (openapi.CreateRecordResponseObject, error) {
	return openapi.CreateRecord201JSONResponse{}, nil
}

// (DELETE /records/{recordId})
func (v *VaultService) DeleteRecord(ctx context.Context, request openapi.DeleteRecordRequestObject) (openapi.DeleteRecordResponseObject, error) {
	return openapi.DeleteRecord204Response{}, nil
}

// (PUT /records/{recordId})
func (v *VaultService) UpdateRecord(ctx context.Context, request openapi.UpdateRecordRequestObject) (openapi.UpdateRecordResponseObject, error) {
	return openapi.UpdateRecord200JSONResponse{}, nil
}
