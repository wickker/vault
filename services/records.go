package services

import (
	"context"
	"vault/openapi"
)

// (GET /records)
func (v *VaultService) GetRecordsByItem(ctx context.Context, request openapi.GetRecordsByItemRequestObject) (openapi.GetRecordsByItemResponseObject, error) {
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
