package services

import (
	"context"
	"fmt"
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
	for i, r := range recordsByItem {
		records[i] = openapi.Record{
			Id:    r.RecordID,
			Name:  r.RecordName,
			Value: r.RecordValue,
		}
	}
	var id int32
	var name string
	if len(records) > 0 {
		id = records[0].Id
		name = records[0].Name
	}
	return openapi.GetRecordsByItem200JSONResponse{
		Id:      id,
		Name:    name,
		Records: records,
	}, nil
}

// (POST /records)
func (v *VaultService) CreateRecord(ctx context.Context, request openapi.CreateRecordRequestObject) (openapi.CreateRecordResponseObject, error) {
	logger := v.getLogger(ctx)
	user, err := v.getUser(ctx)
	if err != nil {
		return openapi.CreateRecord4XXJSONResponse{Body: openapi.Error{
			Message: err.Error(),
		}, StatusCode: 401}, nil
	}

	// check that item belongs to user
	item, err := v.queries.GetItem(ctx, request.Body.ItemId)
	if err != nil {
		logger.Err(err).Msgf("Unable to get item [ItemID: %v].", request.Body.ItemId)
		return openapi.CreateRecord5XXJSONResponse{Body: openapi.Error{
			Message: err.Error(),
		}, StatusCode: 500}, nil
	}
	if item.ClerkUserID != user.ID {
		err = fmt.Errorf("requested item does not belong to user")
		logger.Err(err).Msgf("Unable to match user IDs [itemUserID: %s][currentUserID: %s].", item.ClerkUserID, user.ID)
		return openapi.CreateRecord5XXJSONResponse{Body: openapi.Error{
			Message: err.Error(),
		}, StatusCode: 500}, nil
	}

	record, err := v.queries.CreateRecord(ctx, sqlc.CreateRecordParams{Name: request.Body.Name, Value: request.Body.Value, ItemID: request.Body.ItemId})
	if err != nil {
		logger.Err(err).Msgf("Unable to create record [Name: %s][Value: %s][ItemID: %v].", request.Body.Name, request.Body.Value, request.Body.ItemId)
		return openapi.CreateRecord5XXJSONResponse{Body: openapi.Error{
			Message: err.Error(),
		}, StatusCode: 500}, nil
	}

	return openapi.CreateRecord201JSONResponse{
		Id:    record.ID,
		Name:  record.Name,
		Value: record.Value,
	}, nil
}

// (DELETE /records/{recordId})
func (v *VaultService) DeleteRecord(ctx context.Context, request openapi.DeleteRecordRequestObject) (openapi.DeleteRecordResponseObject, error) {
	logger := v.getLogger(ctx)
	user, err := v.getUser(ctx)
	if err != nil {
		return openapi.DeleteRecord4XXJSONResponse{Body: openapi.Error{
			Message: err.Error(),
		}, StatusCode: 401}, nil
	}

	// check that record belongs to user
	recordUserID, err := v.queries.GetRecordUserID(ctx, request.RecordId)
	if err != nil {
		logger.Err(err).Msgf("Unable to get record userID [recordID: %v].", request.RecordId)
		return openapi.DeleteRecord5XXJSONResponse{Body: openapi.Error{
			Message: err.Error(),
		}, StatusCode: 500}, nil
	}
	if recordUserID != user.ID {
		err = fmt.Errorf("requested item does not belong to user")
		logger.Err(err).Msgf("Unable to match user IDs [itemUserID: %s][currentUserID: %s].", recordUserID, user.ID)
		return openapi.DeleteRecord5XXJSONResponse{Body: openapi.Error{
			Message: err.Error(),
		}, StatusCode: 500}, nil
	}

	record, err := v.queries.DeleteRecord(ctx, request.RecordId)
	if err != nil || record.ID == 0 {
		logger.Err(err).Msgf("Unable to delete record [recordID: %v].", request.RecordId)
		return openapi.DeleteRecord5XXJSONResponse{Body: openapi.Error{
			Message: err.Error(),
		}, StatusCode: 500}, nil
	}

	return openapi.DeleteRecord204Response{}, nil
}

// (PUT /records/{recordId})
func (v *VaultService) UpdateRecord(ctx context.Context, request openapi.UpdateRecordRequestObject) (openapi.UpdateRecordResponseObject, error) {
	logger := v.getLogger(ctx)
	user, err := v.getUser(ctx)
	if err != nil {
		return openapi.UpdateRecord4XXJSONResponse{Body: openapi.Error{
			Message: err.Error(),
		}, StatusCode: 401}, nil
	}

	// check that record belongs to user
	recordUserID, err := v.queries.GetRecordUserID(ctx, request.RecordId)
	if err != nil {
		logger.Err(err).Msgf("Unable to get record userID [recordID: %v].", request.RecordId)
		return openapi.UpdateRecord5XXJSONResponse{Body: openapi.Error{
			Message: err.Error(),
		}, StatusCode: 500}, nil
	}
	if recordUserID != user.ID {
		err = fmt.Errorf("requested item does not belong to user")
		logger.Err(err).Msgf("Unable to match user IDs [itemUserID: %s][currentUserID: %s].", recordUserID, user.ID)
		return openapi.UpdateRecord5XXJSONResponse{Body: openapi.Error{
			Message: err.Error(),
		}, StatusCode: 500}, nil
	}

	record, err := v.queries.UpdateRecord(ctx, sqlc.UpdateRecordParams{
		Name:  request.Body.Name,
		Value: request.Body.Value,
		ID:    request.RecordId,
	})
	if err != nil || record.ID == 0 {
		logger.Err(err).Msgf("Unable to update record [ID: %v][Name: %s][Value: %s].", request.RecordId, request.Body.Name, request.Body.Value)
		return openapi.UpdateRecord5XXJSONResponse{Body: openapi.Error{
			Message: err.Error(),
		}, StatusCode: 500}, nil
	}

	return openapi.UpdateRecord200JSONResponse{
		Id:    record.ID,
		Value: record.Value,
		Name:  record.Name,
	}, nil
}
