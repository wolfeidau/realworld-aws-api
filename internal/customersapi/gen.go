package customersapi

//go:generate go run github.com/deepmap/oapi-codegen/cmd/oapi-codegen -generate types,server,spec,client -package customersapi -o customersapi.gen.go ../../openapi/customers.yaml
//go:generate go fmt customersapi.gen.go
