# Makefile for AI Development Learning Project

.PHONY: help setup build run test clean docker-up docker-down

# Default target
help:
	@echo "🤖 AI应用开发学习项目"
	@echo "====================="
	@echo ""
	@echo "可用命令:"
	@echo "  setup      - 初始化项目环境"
	@echo "  build      - 构建Go应用"
	@echo "  run        - 运行Go API服务"
	@echo "  run-python - 运行Python AI服务"
	@echo "  test       - 运行API测试"
	@echo "  test-basic - 运行基础功能测试"
	@echo "  clean      - 清理构建文件"
	@echo ""

# 项目初始化
setup:
	@echo "🚀 初始化项目环境..."
	@./scripts/setup.sh

# 构建Go应用
build:
	@echo "🔨 构建Go应用..."
	@go mod tidy
	@go build -o bin/api cmd/api/main.go
	@echo "✅ 构建完成: bin/api"

# 运行Go API服务
run: build
	@echo "🚀 启动Go API服务..."
	@./bin/api

# 运行Python AI服务
run-python:
	@echo "🐍 启动Python AI服务..."
	@cd scripts && python api_server.py

# 运行API测试
test:
	@echo "🧪 运行API测试..."
	@./scripts/test_api.sh

# 运行基础测试
test-basic:
	@echo "🧪 运行基础功能测试..."
	@./scripts/test_basic.sh

# 清理构建文件
clean:
	@echo "🧹 清理构建文件..."
	@rm -rf bin/
	@rm -rf data/chroma_db/
	@echo "✅ 清理完成"

# 安装Go依赖
deps:
	@echo "📦 安装Go依赖..."
	@go mod download
	@go mod tidy

# 安装Python依赖
python-deps:
	@echo "🐍 安装Python依赖..."
	@pip install -r scripts/requirements.txt

# 数据库迁移
migrate:
	@echo "🗄️  运行数据库迁移..."
	@mysql -u ai_user -pai_pass ai_development < setup_database.sql

# 开发模式（实时重载）
dev:
	@echo "⚡ 启动开发模式..."
	@go run -ldflags="-X main.mode=dev" cmd/api/main.go

# 生产构建
prod-build:
	@echo "🏭 生产环境构建..."
	@CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/api-prod cmd/api/main.go

# 格式化代码
fmt:
	@echo "✨ 格式化Go代码..."
	@go fmt ./...
	@echo "✨ 格式化Python代码..."
	@cd scripts && python -m black .

# 代码检查
lint:
	@echo "🔍 运行代码检查..."
	@go vet ./...
	@golint ./...

# 性能测试
benchmark:
	@echo "⚡ 运行性能测试..."
	@go test -bench=. ./...

# 生成API文档
docs:
	@echo "📚 生成API文档..."
	@echo "文档已存在于 docs/ 目录"
	@echo "- API文档: docs/api.md"
	@echo "- 学习指南: docs/learning-guide.md"

# 快速启动（推荐用于初次使用）
quickstart: setup
	@echo ""
	@echo "🎉 快速启动完成！"
	@echo ""
	@echo "📚 下一步："
	@echo "1. 确保MySQL和Redis服务运行: brew services list"
	@echo "2. 在 .env 文件中设置 OpenAI API Key"
	@echo "3. 启动API服务: make run"
	@echo "4. 运行基础测试: make test-basic"
	@echo "5. 查看 docs/learning-guide.md 开始学习"
