APPNAME := realworld-aws-api
STAGE ?= dev
BRANCH ?= master

GOLANGCI_VERSION = 1.31.0

GIT_HASH := $(shell git rev-parse --short HEAD)
BUILD_DATE := $(shell date -u '+%Y%m%dT%H%M%S')

default: clean generate build archive package deploy

ci: clean generate lint test
.PHONY: ci

LDFLAGS := -ldflags="-s -w -X github.com/wolfeidau/realworld-aws-api/internal/app.BuildDate=${BUILD_DATE} -X github.com/wolfeidau/realworld-aws-api/internal/app.Commit=${GIT_HASH}"

bin/golangci-lint: bin/golangci-lint-${GOLANGCI_VERSION}
	@ln -sf golangci-lint-${GOLANGCI_VERSION} bin/golangci-lint
bin/golangci-lint-${GOLANGCI_VERSION}:
	@curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | BINARY=golangci-lint bash -s -- v${GOLANGCI_VERSION}
	@mv bin/golangci-lint $@

bin/mockgen:
	@env GOBIN=$$PWD/bin GO111MODULE=on go install github.com/golang/mock/mockgen

bin/gcov2lcov:
	@env GOBIN=$$PWD/bin GO111MODULE=on go install github.com/jandelgado/gcov2lcov

clean:
	@echo "--- clean all the things"
	@rm -rf ./dist
	@rm -f ./handler.zip
	@rm -f ./*.out.yaml
.PHONY: clean

lint: bin/golangci-lint
	@echo "--- lint all the things"
	@bin/golangci-lint run
.PHONY: lint

test: bin/gcov2lcov
	@echo "--- test all the things"
	@go test -v -covermode=count -coverprofile=coverage.txt ./internal/...
	@bin/gcov2lcov -infile=coverage.txt -outfile=coverage.lcov
.PHONY: test

generate:
	@echo "--- generate all the things"
	@go generate ./...
.PHONY: generate

build:
	@echo "--- build all the things"
	@mkdir -p dist
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -trimpath -o dist ./cmd/...
.PHONY: build

archive:
	@echo "--- build an archive"
	@cd dist && zip -X -9 -r ../handler.zip *-lambda
.PHONY: archive

package:
	@echo "--- uploading cloudformation assets to $(S3_BUCKET)"
	@aws cloudformation package \
		--template-file sam/api.yaml \
		--output-template-file api.out.yaml \
		--s3-bucket $(S3_BUCKET) \
		--s3-prefix sam
.PHONY: package

deploy:
	@echo "--- deploy stack $(APPNAME)-$(STAGE)-$(BRANCH)"
	@aws cloudformation deploy \
		--no-fail-on-empty-changeset \
		--template-file api.out.yaml \
		--capabilities CAPABILITY_IAM \
		--stack-name $(APPNAME)-$(STAGE)-$(BRANCH) \
		--parameter-overrides AppName=$(APPNAME) Stage=$(STAGE) Branch=$(BRANCH)
.PHONY: deploy
