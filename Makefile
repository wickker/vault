server:
	go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen --config=openapi/serverconfig.yaml openapi/openapi.yaml

models:
	go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen --config=openapi/modelsconfig.yaml openapi/openapi.yaml