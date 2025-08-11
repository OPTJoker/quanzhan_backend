#!/bin/bash

# AI应用开发实战 - 快速启动脚本

echo "🤖 AI应用开发实战学习项目"
echo "================================"

# 检查Go环境
if ! command -v go &> /dev/null; then
    echo "❌ Go未安装，请先安装Go环境"
    exit 1
fi

echo "✅ Go环境检查通过"

# 检查MySQL
if ! command -v mysql &> /dev/null; then
    echo "❌ MySQL未安装，请先安装MySQL"
    echo "macOS用户可运行: brew install mysql"
    exit 1
fi

echo "✅ MySQL环境检查通过"

# 检查Redis
if ! command -v redis-server &> /dev/null; then
    echo "⚠️  Redis未安装，正在安装..."
    if command -v brew &> /dev/null; then
        brew install redis
    else
        echo "❌ 请手动安装Redis"
        exit 1
    fi
fi

echo "✅ Redis环境检查通过"

# 启动服务
echo "🚀 启动必要服务..."

# 启动MySQL（如果未运行）
if ! brew services list | grep mysql | grep -q started; then
    echo "启动MySQL服务..."
    brew services start mysql
fi

# 启动Redis（如果未运行）
if ! brew services list | grep redis | grep -q started; then
    echo "启动Redis服务..."
    brew services start redis
fi

echo "⏳ 等待服务启动..."
sleep 3

# 创建环境配置文件
echo "📝 创建配置文件..."
cat > .env << EOF
APP_PORT=8080
MYSQL_DSN=ai_user:ai_pass@tcp(localhost:3306)/ai_development?charset=utf8mb4&parseTime=True&loc=Local
REDIS_ADDR=localhost:6379
REDIS_PASSWORD=
REDIS_DB=0
OPENAI_API_KEY=sk-6700f2daa39b4f63bee5b615df14ea89
OPENAI_BASE_URL=https://dashscope.aliyuncs.com/compatible-mode/v1
EOF

echo "✅ 环境配置文件已创建"

# 设置数据库
echo "🗄️  设置数据库..."
if mysql -u root -e "SELECT 1" &> /dev/null; then
    mysql -u root < setup_database.sql
    echo "✅ 数据库设置完成"
else
    echo "⚠️  请手动运行数据库设置: mysql -u root -p < setup_database.sql"
fi

# 安装Go依赖
echo "📦 安装Go依赖..."
go mod tidy

echo "🎉 设置完成！"
echo ""
echo "下一步："
echo "1. 如需AI功能，在 .env 文件中设置你的 OpenAI API Key"
echo "2. 运行API服务: go run cmd/api/main.go"
echo "3. 测试功能: make test-basic"
echo "4. 访问: http://localhost:8080/health"
echo ""
echo "API文档："
echo "- POST /api/v1/chats - 创建聊天"
echo "- POST /api/v1/chats/{sessionId}/messages - 发送消息"
echo "- GET /api/v1/chats/{sessionId} - 获取历史"
