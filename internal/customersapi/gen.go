package customersapi

//go:generate go get github.com/deepmap/oapi-codegen/cmd/oapi-codegen
//go:generate env GOBIN=$PWD/bin GO111MODULE=on go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen
//go:generate $PWD/bin/oapi-codegen -templates ../../codegen -generate types,server,spec,client -package customersapi -o customersapi.gen.go ../../openapi/customers.yaml
