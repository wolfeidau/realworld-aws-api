package customersapi

//go:generate go run github.com/deepmap/oapi-codegen/cmd/oapi-codegen --config=types.cfg.yaml ../../openapi/customers.yaml
//go:generate go run github.com/deepmap/oapi-codegen/cmd/oapi-codegen --config=server.cfg.yaml ../../openapi/customers.yaml
//go:generate go run github.com/deepmap/oapi-codegen/cmd/oapi-codegen --config=client.cfg.yaml ../../openapi/customers.yaml
