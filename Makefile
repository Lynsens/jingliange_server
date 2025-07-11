PROJECT_NAME := "jingliange_server"
PKG := "github.com/lynsens/jingliange_server"

.DEFAULT_GOAL := help
.PHONY: all build build-linux build-windows build-darwin build-all test test-full clean run swagger dev fmt lint deps install help

## help: 显示帮助信息
help:
	@echo "净莲阁后端项目 Makefile"
	@echo ""
	@echo "可用命令:"
	@echo "  build         - 编译项目 (当前平台)"
	@echo "  build-linux   - 编译Linux版本"
	@echo "  build-windows - 编译Windows版本"
	@echo "  build-darwin  - 编译macOS版本"
	@echo "  build-all     - 编译所有平台版本"
	@echo "  test          - 运行单元测试"
	@echo "  test-full     - 运行完整测试套件"
	@echo "  clean         - 清理构建文件"
	@echo "  run           - 运行开发服务器"
	@echo "  swagger       - 生成Swagger文档"
	@echo "  dev           - 开发模式 (清理+编译+运行)"
	@echo "  fmt           - 格式化代码"
	@echo "  lint          - 代码检查"
	@echo "  deps          - 更新依赖"
	@echo "  install       - 安装开发工具"
	@echo "  all           - 执行完整的开发流程"
	@echo ""

## build: 编译项目 (当前平台)
build:
	@echo "编译项目 (当前平台)..."
	@go build -o bin/jingliange_server cmd/main.go
	@echo "编译完成: bin/jingliange_server"

## build-linux: 编译Linux版本
build-linux:
	@echo "编译Linux版本..."
	@GOOS=linux GOARCH=amd64 go build -o bin/jingliange_server cmd/main.go
	@echo "Linux版本编译完成: bin/jingliange_server_linux"

## build-windows: 编译Windows版本
build-windows:
	@echo "编译Windows版本..."
	@GOOS=windows GOARCH=amd64 go build -o bin/jingliange_server.exe cmd/main.go
	@echo "Windows版本编译完成: bin/jingliange_server.exe"

## build-darwin: 编译macOS版本
build-darwin:
	@echo "编译macOS版本..."
	@GOOS=darwin GOARCH=amd64 go build -o bin/jingliange_server cmd/main.go
	@echo "macOS版本编译完成: bin/jingliange_server_darwin"

## build-all: 编译所有平台版本
build-all: build-linux build-windows build-darwin
	@echo "所有平台版本编译完成"

## test: 运行快速单元测试
test:
	@echo "运行单元测试..."
	@go test ./pkg/util/jwt_simple_test.go ./pkg/util/jwt.go -v
	@go test ./internal/router/api/v1/api_simple_test.go -v
	@echo "单元测试完成"

## test-full: 运行完整测试套件
test-full:
	@echo "运行完整测试套件..."
	@./run_tests.sh

## clean: 清理构建文件
clean:
	@echo "清理构建文件..."
	@rm -f bin/jingliange_server
	@rm -f bin/jingliange_server_linux
	@rm -f bin/jingliange_server.exe
	@rm -f bin/jingliange_server_darwin
	@rm -f main
	@echo "清理完成"

## run: 运行开发服务器
run: build
	@echo "启动开发服务器..."
	@./bin/jingliange_server

## swagger: 生成Swagger文档
swagger:
	@echo "生成Swagger文档..."
	@swag init -g cmd/main.go -o docs/
	@echo "Swagger文档生成完成"

## dev: 开发模式 (清理+编译+运行)
dev: clean build
	@echo "开发模式启动..."
	@./bin/jingliange_server

## fmt: 格式化代码
fmt:
	@echo "格式化代码..."
	@go fmt ./...
	@echo "代码格式化完成"

## lint: 代码检查 (需要安装golangci-lint)
lint:
	@echo "代码检查..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint 未安装，跳过代码检查"; \
		echo "安装命令: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

## deps: 更新和整理依赖
deps:
	@echo "更新依赖..."
	@go mod tidy
	@go mod download
	@echo "依赖更新完成"

## install: 安装开发工具
install:
	@echo "安装开发工具..."
	@go install github.com/swaggo/swag/cmd/swag@latest
	@echo "开发工具安装完成"

## all: 执行完整的开发流程
all: clean fmt test build
	@echo "完整开发流程执行完成"