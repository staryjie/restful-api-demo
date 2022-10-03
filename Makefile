PROJECT_NAME=restful-api-demo
MAIN_FILE=main.go
PKG := "github.com/staryjie/$(PROJECT_NAME)"
MOD_DIR := $(shell go env GOMODCACHE)
PKG_LIST := $(shell go list ${PKG}/... | grep -v /vendor/)
GO_FILES := $(shell find . -name '*.go' | grep -v /vendor/ | grep -v _test.go)

BUILD_BRANCH := $(shell git rev-parse --abbrev-ref HEAD)
BUILD_COMMIT := ${shell git rev-parse HEAD}
BUILD_TIME := ${shell date '+%Y-%m-%d %H:%M:%S'}
BUILD_GO_VERSION := $(shell go version | grep -o  'go[0-9].[0-9].*')
VERSION_PATH := "${PKG}/version"

MCUBE_MODULE := "github.com/infraboard/mcube"
MCUBE_VERSION := $(shell go list -m "github.com/infraboard/mcube" | cut -d ' ' -f2)
MCUBE_PKG_PATH := ${MOD_DIR}/${MCUBE_MODULE}@${MCUBE_VERSION}

.PHONY: all dep lint vet test test-coverage build clean

all: build

dep: ## Get the dependencies
	@go mod tidy

lint: ## Lint Golang files
	@golint -set_exit_status ${PKG_LIST}

vet: ## Run go vet
	@go vet ${PKG_LIST}

test: ## Run unittests
	@go test -short ${PKG_LIST}

test-coverage: ## Run tests with coverage
	@go test -short -coverprofile cover.out -covermode=atomic ${PKG_LIST}
	@cat cover.out >> coverage.txt

build: dep ## Build the binary file
	@go build -ldflags "-s -w" -ldflags "-X ${VERSION_PATH}.GIT_TAG='0.0.1' -X '${VERSION_PATH}.GIT_BRANCH=${BUILD_BRANCH}' -X '${VERSION_PATH}.GIT_COMMIT=${BUILD_COMMIT}' -X '${VERSION_PATH}.BUILD_TIME=${BUILD_TIME}' -X '${VERSION_PATH}.GO_VERSION=${BUILD_GO_VERSION}'" -o dist/demo-api.exe $(MAIN_FILE)

linux: dep ## Build the binary file
	@GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -ldflags "-X '${VERSION_PATH}.GIT_BRANCH=${BUILD_BRANCH}' -X '${VERSION_PATH}.GIT_COMMIT=${BUILD_COMMIT}' -X '${VERSION_PATH}.BUILD_TIME=${BUILD_TIME}' -X '${VERSION_PATH}.GO_VERSION=${BUILD_GO_VERSION}'" -o dist/demo-api $(MAIN_FILE)

pb: ## Copy mcube protobuf files to common/pb
	@mkdir -p common/pb/github.com/infraboard/mcube/pb/
	@cp -r ${MCUBE_PKG_PATH}/pb/* common/pb/github.com/infraboard/mcube/pb/
	@rm -rf common/pb/github.com/infraboard/mcube/pb/*/*.go

gen: ## Init Service
	@protoc -I=. -I=common/pb -I=/usr/local/include/ --go_out=. --go_opt=module=${PKG} --go-grpc_out=. --go-grpc_opt=module=${PKG} apps/*/pb/*.proto
    ##  go install github.com/favadi/protoc-go-inject-tag@latest
	@protoc-go-inject-tag -input=apps/*/*.pb.go
    ##  go install github.com/infraboard/mcube/cmd/mcube@latest
	@mcube generate enum -p -m apps/*/*.pb.go
	@go mod tidy
	@go fmt ./...

run: # Run Develop server
	@go run $(MAIN_FILE) start -f etc/demo.toml

clean: ## Remove previous build
	@rm -f dist/*

push: ## push code to github
	@git push -u origin main

help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'