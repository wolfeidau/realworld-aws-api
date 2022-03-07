APPNAME := realworld-aws-api
STAGE ?= dev
BRANCH ?= master
SAR_VERSION ?= 1.0.0

WORKDIR = $(shell pwd)
GOLANGCI_LINT_TAG = v1.43.0

GIT_HASH := $(shell git rev-parse --short HEAD)
BUILD_DATE := $(shell date -u '+%Y%m%dT%H%M%S')

# This path is used to cache binaries used for development and can be overridden to avoid issues with osx vs linux
# binaries.
GOBIN ?= $(shell go env GOPATH)/bin

default: clean generate build archive deploy-bucket package deploy

ci: clean generate lint test
.PHONY: ci

LDFLAGS := -ldflags="-s -w -X github.com/wolfeidau/realworld-aws-api/internal/app.BuildDate=${BUILD_DATE} -X github.com/wolfeidau/realworld-aws-api/internal/app.Commit=${GIT_HASH}"

mocks:
	@echo "--- build all the mocks"
	@go run github.com/golang/mock/mockgen -destination=mocks/customers_store.go -package=mocks github.com/wolfeidau/realworld-aws-api/internal/stores Customers
.PHONY: mocks

clean:
	@echo "--- clean all the things"
	@rm -rf ./dist
	@rm -f ./handler.zip
	@rm -f ./*.out.yaml
.PHONY: clean

lint:
	@echo "--- lint all the things"
	@docker run --rm -v ${WORKDIR}:/app -w /app golangci/golangci-lint:$(GOLANGCI_LINT_TAG) golangci-lint run -v ./...
.PHONY: lint

lint-fix:
	@echo "--- lint all the things"
	@docker run --rm -v ${WORKDIR}:/app -w /app golangci/golangci-lint:$(GOLANGCI_LINT_TAG) golangci-lint run --fix -v ./...
.PHONY: lint-fix

test:
	@echo "--- test all the things"
	@go test -v -covermode=count -coverprofile=coverage.txt ./internal/...
.PHONY: test

generate:
	@echo "--- generate all the things"
	@go generate ./...
.PHONY: generate

$(GOBIN)/protoc-gen-go:
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@latest

proto: $(GOBIN)/protoc-gen-go proto/customers/storage/v1beta1/storage.pb.go

proto/customers/storage/v1beta1/storage.pb.go: proto/customers/storage/v1beta1/storage.proto
	protoc -I proto --go_out=paths=source_relative:proto --plugin=$(GOBIN)/protoc-gen-go proto/customers/storage/v1beta1/storage.proto

build:
	@echo "--- build all the things"
	@mkdir -p dist
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -trimpath -o dist ./cmd/...
.PHONY: build

archive:
	@echo "--- build an archive"
	@cd dist && zip -X -9 -r ../handler.zip *-lambda
.PHONY: archive

deploy-bucket:
	@sam deploy \
		--no-fail-on-empty-changeset \
		--template-file sam/deploy.yaml \
		--capabilities CAPABILITY_IAM \
		--stack-name $(APPNAME)-$(STAGE)-$(BRANCH)-deploybucket \
		--tags "environment=$(STAGE)" "branch=$(BRANCH)" "service=$(APPNAME)" \
		--parameter-overrides \
			AppName=$(APPNAME) \
			Stage=$(STAGE) \
			Branch=$(BRANCH)
.PHONY: deploy-bucket

package:
	@echo "--- uploading cloudformation assets to $(S3_BUCKET)"
	@sam package \
		--template-file sam/api.yaml \
		--output-template-file api.out.yaml \
		--s3-bucket $(shell aws ssm get-parameter --name "/config/$(STAGE)/$(BRANCH)/$(APPNAME)/deploy_bucket" --query 'Parameter.Value' --output text) \
		--s3-prefix sam/$(GIT_HASH)
.PHONY: package

publish:
	@echo "--- publish stack $(APPNAME)-$(STAGE)-$(BRANCH)"
	@sam publish \
		--template-file api.out.yaml \
		--semantic-version $(SAR_VERSION)
.PHONY: publish

deploy:
	@echo "--- deploy stack $(APPNAME)-$(STAGE)-$(BRANCH)"
	@sam deploy \
		--no-fail-on-empty-changeset \
		--template-file api.out.yaml \
		--capabilities CAPABILITY_IAM \
		--tags "environment=$(STAGE)" "branch=$(BRANCH)" "service=$(APPNAME)" \
		--stack-name $(APPNAME)-$(STAGE)-$(BRANCH) \
		--parameter-overrides AppName=$(APPNAME) Stage=$(STAGE) Branch=$(BRANCH)
.PHONY: deploy

watch:
	@echo "-- watch for changes and run local server"
	@go run github.com/cespare/reflex -s -r '\.go$$' go run cmd/api-server/main.go
.PHONY: deploy

docker-compose:
	@echo "-- run docker-compose in the foreground for local development"
	@docker-compose up
.PHONY: docker-compose
