PROJECT_NAME := "jingliange_server"
PKG := "github.com/lynsens/jingliange_server"
PKG_LIST := $(shell go mod tidy && go list ${PKG}/...)
BUS_PKG_LIST := $(shell go mod tidy && go list ${PKG}/... | grep -v stub_test)
GO_FILES := $(shell find . -name '*.go' | grep -v _test.go)
COVERAGE_PKG_LIST := $(shell go mod tidy && go list ${PKG}/... | grep -E -v "cmd|stub_test|pb")
COVERAGE_PKGS := $(shell echo ${COVERAGE_PKG_LIST} | sed  "s/ /,/g")
UNITTEST_COV = "unittest_cov"


.DEFAULT_GOAL := default
.PHONY: all

all: fmt lint vet test race build

clean-gen: ## Clean generated files
	@echo "clean generated files..."
	@	rm pb/*
	# @   rm swagger/*

gen: ## Generate protoc
	@echo "go generate..."
	@protoc --proto_path=proto proto/*.proto  --go_out=:pb --go-grpc_out=:pb --grpc-gateway_out=:pb  
	# @protoc --proto_path=proto proto/*.proto  --go_out=:pb --go-grpc_out=:pb --grpc-gateway_out=:pb --openapiv2_out=:swagger

dep: ## Get dependencies
	@echo "go dep..."
	@go mod tidy

fmt: dep ## Format code
	@echo "go fmt..."
	@go fmt $(PKG_LIST)

vet: dep ## Vet check
	@echo "go vet..."
	@go vet -all $(PKG_LIST)

test: dep ## Run unittests
	@echo "go test..."
	@go test -short -v -count=1 -p=1 `go list ./... | grep -v stub_test` 2>&1||true

race: dep ## Run data race detector
	@echo "go test race..."
	@go test -gcflags=all=-l -race -short -v -count=1 ${BUS_PKG_LIST}

build: dep fmt ## Build frpc project
	@echo "go build..."
	@CGO_ENABLED=1 go build -v -buildmode=default -o bin/${PROJECT_NAME} cmd/main.go cmd/metrics.go cmd/clients.go
#	@CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -v -gcflags=all="-N -l" -o bin/${PROJECT_NAME} cmd/main.go
	@chmod +x bin/${PROJECT_NAME}

build-dev: dep fmt ## Build  project
	@echo "go build..."
	@go build -v -gcflags=all="-N -l" -o bin/${PROJECT_NAME} cmd/main.go
	@chmod +x bin/${PROJECT_NAME}

run: dep fmt ## Run  project
	@echo "go run..."
#   @go run cmd/main.go --config=conf/conf.toml
	@go run cmd/main.go 
 
stub_test: fmt ## Run stub test
	@echo "go run stub_test..."
	@go test -gcflags=all=-l -v -count=1 ./stub_test/...

help: ## Display this help screen
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

default: help