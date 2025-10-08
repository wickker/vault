package services

import (
	"context"
	"fmt"
	"vault/db/sqlc"
	"vault/openapi"
	"vault/utils"
)

// (GET /records)
func (v *VaultService) GetRecordsByItem(ctx context.Context, request openapi.GetRecordsByItemRequestObject) (openapi.GetRecordsByItemResponseObject, error) {
	ctxValues, err := v.getContextValues(ctx)
	if err != nil {
		return openapi.GetRecordsByItem4XXJSONResponse{Body: openapi.Error{
			Message: err.Error(),
		}, StatusCode: 401}, nil
	}
	ctx = ctxValues.GinContext
	logger := ctxValues.Logger
	user := ctxValues.User

	// check that item belongs to user
	item, err := v.queries.GetItem(ctx, request.Params.ItemId)
	if err != nil {
		logger.Err(err).Msgf("Unable to get item [ItemID: %v].", request.Params.ItemId)
		return openapi.GetRecordsByItem5XXJSONResponse{Body: openapi.Error{
			Message: err.Error(),
		}, StatusCode: 500}, nil
	}
	if item.ClerkUserID != user.ID {
		err = fmt.Errorf("requested item does not belong to user")
		logger.Err(err).Msgf("Unable to match user IDs [itemUserID: %s][currentUserID: %s][ItemID: %v].", item.ClerkUserID, user.ID, request.Params.ItemId)
		return openapi.GetRecordsByItem5XXJSONResponse{Body: openapi.Error{
			Message: err.Error(),
		}, StatusCode: 500}, nil
	}

	recordsByItem, err := v.queries.ListRecordsByItemId(ctx, request.Params.ItemId)
	if err != nil {
		logger.Err(err).Msgf("Unable to list records by item [UserID: %s][ItemID: %v].", user.ID, request.Params.ItemId)
		return openapi.GetRecordsByItem5XXJSONResponse{Body: openapi.Error{
			Message: err.Error(),
		}, StatusCode: 500}, nil
	}

	records := make([]openapi.Record, len(recordsByItem))
	for i, r := range recordsByItem {
		decrypted, err := utils.Decrypt(r.Value, []byte(v.encryptionKey))
		if err != nil {
			logger.Err(err).Msgf("Unable to decrypt value [RecordID: %v].", r.ID)
			return openapi.GetRecordsByItem5XXJSONResponse{Body: openapi.Error{
				Message: err.Error(),
			}, StatusCode: 500}, nil
		}

		records[i] = openapi.Record{
			Id:    r.ID,
			Name:  r.Name,
			Value: decrypted,
		}
	}
	return openapi.GetRecordsByItem200JSONResponse{
		Id:      item.ID,
		Name:    item.Name,
		Records: records,
	}, nil
}

// (POST /records)
func (v *VaultService) CreateRecord(ctx context.Context, request openapi.CreateRecordRequestObject) (openapi.CreateRecordResponseObject, error) {
	ctxValues, err := v.getContextValues(ctx)
	if err != nil {
		return openapi.CreateRecord4XXJSONResponse{Body: openapi.Error{
			Message: err.Error(),
		}, StatusCode: 401}, nil
	}
	ctx = ctxValues.GinContext
	logger := ctxValues.Logger
	user := ctxValues.User

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
		logger.Err(err).Msgf("Unable to match user IDs [itemUserID: %s][currentUserID: %s][ItemID: %v].", item.ClerkUserID, user.ID, request.Body.ItemId)
		return openapi.CreateRecord5XXJSONResponse{Body: openapi.Error{
			Message: err.Error(),
		}, StatusCode: 500}, nil
	}

	encrypted, err := utils.Encrypt(request.Body.Value, []byte(v.encryptionKey))
	if err != nil {
		logger.Err(err).Msgf("Unable to encrypt record value [Request: %+v].", request)
		return openapi.CreateRecord5XXJSONResponse{Body: openapi.Error{
			Message: err.Error(),
		}, StatusCode: 500}, nil
	}
	record, err := v.queries.CreateRecord(ctx, sqlc.CreateRecordParams{Name: request.Body.Name, Value: encrypted, ItemID: request.Body.ItemId})
	if err != nil {
		logger.Err(err).Msgf("Unable to create record [Name: %s][Value: %s][ItemID: %v].", request.Body.Name, encrypted, request.Body.ItemId)
		return openapi.CreateRecord5XXJSONResponse{Body: openapi.Error{
			Message: err.Error(),
		}, StatusCode: 500}, nil
	}

	return openapi.CreateRecord201JSONResponse{
		Id:    record.ID,
		Name:  record.Name,
		Value: request.Body.Value,
	}, nil
}

// (DELETE /records/{recordId})
func (v *VaultService) DeleteRecord(ctx context.Context, request openapi.DeleteRecordRequestObject) (openapi.DeleteRecordResponseObject, error) {
	ctxValues, err := v.getContextValues(ctx)
	if err != nil {
		return openapi.DeleteRecord4XXJSONResponse{Body: openapi.Error{
			Message: err.Error(),
		}, StatusCode: 401}, nil
	}
	ctx = ctxValues.GinContext
	logger := ctxValues.Logger
	user := ctxValues.User

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
	ctxValues, err := v.getContextValues(ctx)
	if err != nil {
		return openapi.UpdateRecord4XXJSONResponse{Body: openapi.Error{
			Message: err.Error(),
		}, StatusCode: 401}, nil
	}
	ctx = ctxValues.GinContext
	logger := ctxValues.Logger
	user := ctxValues.User

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

	encrypted, err := utils.Encrypt(request.Body.Value, []byte(v.encryptionKey))
	if err != nil {
		logger.Err(err).Msgf("Unable to encrypt record value [Request: %+v].", request)
		return openapi.UpdateRecord5XXJSONResponse{Body: openapi.Error{
			Message: err.Error(),
		}, StatusCode: 500}, nil
	}
	record, err := v.queries.UpdateRecord(ctx, sqlc.UpdateRecordParams{
		Name:  request.Body.Name,
		Value: encrypted,
		ID:    request.RecordId,
	})
	if err != nil || record.ID == 0 {
		logger.Err(err).Msgf("Unable to update record [ID: %v][Name: %s][Value: %s].", request.RecordId, request.Body.Name, encrypted)
		return openapi.UpdateRecord5XXJSONResponse{Body: openapi.Error{
			Message: err.Error(),
		}, StatusCode: 500}, nil
	}

	return openapi.UpdateRecord200JSONResponse{
		Id:    record.ID,
		Value: request.Body.Value,
		Name:  record.Name,
	}, nil
}
