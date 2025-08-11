#!/bin/bash

# AI应用开发 - 测试脚本

echo "🧪 AI应用开发 - API测试"
echo "=========================="

API_BASE="http://localhost:8080/api/v1"
PYTHON_API_BASE="http://localhost:8000/api/v1"

# 检查服务是否运行
echo "🔍 检查服务状态..."

# 检查Go API
if curl -s "$API_BASE/../health" > /dev/null; then
    echo "✅ Go API服务正常运行"
else
    echo "❌ Go API服务未运行，请先启动: go run cmd/api/main.go"
    exit 1
fi

# 检查Python API
if curl -s "$PYTHON_API_BASE/../health" > /dev/null 2>&1; then
    echo "✅ Python AI服务正常运行"
else
    echo "⚠️  Python AI服务未运行（可选）"
fi

echo ""
echo "🚀 开始API测试..."

# 测试1: 创建聊天会话
echo "📝 测试1: 创建聊天会话"
CHAT_RESPONSE=$(curl -s -X POST "$API_BASE/chats" \
    -H "Content-Type: application/json" \
    -d '{"title": "API测试对话"}')

echo "响应: $CHAT_RESPONSE"

# 从响应中提取session_id（需要jq工具）
if command -v jq &> /dev/null; then
    SESSION_ID=$(echo $CHAT_RESPONSE | jq -r '.data.session_id')
    echo "会话ID: $SESSION_ID"
else
    echo "⚠️  需要安装jq工具来解析JSON响应"
    SESSION_ID="test-session-123"
    echo "使用默认会话ID: $SESSION_ID"
fi

echo ""

# 测试2: 发送消息
echo "💬 测试2: 发送消息"
MESSAGE_RESPONSE=$(curl -s -X POST "$API_BASE/chats/$SESSION_ID/messages" \
    -H "Content-Type: application/json" \
    -d '{"message": "你好，请介绍一下你自己"}')

echo "响应: $MESSAGE_RESPONSE"
echo ""

# 测试3: 获取聊天历史
echo "📚 测试3: 获取聊天历史"
HISTORY_RESPONSE=$(curl -s -X GET "$API_BASE/chats/$SESSION_ID")
echo "响应: $HISTORY_RESPONSE"
echo ""

# 测试4: Python RAG查询（如果服务可用）
if curl -s "$PYTHON_API_BASE/../health" > /dev/null 2>&1; then
    echo "🔍 测试4: RAG知识库查询"
    RAG_RESPONSE=$(curl -s -X POST "$PYTHON_API_BASE/rag/query" \
        -H "Content-Type: application/json" \
        -d '{"question": "什么是人工智能？"}')
    echo "响应: $RAG_RESPONSE"
    echo ""

    echo "🤖 测试5: Agent任务执行"
    AGENT_RESPONSE=$(curl -s -X POST "$PYTHON_API_BASE/agent/execute" \
        -H "Content-Type: application/json" \
        -d '{"task": "计算 123 + 456", "session_id": "test-agent-001"}')
    echo "响应: $AGENT_RESPONSE"
    echo ""
fi

# 性能测试
echo "⚡ 性能测试: 并发请求"
echo "发送10个并发请求..."

for i in {1..10}; do
    curl -s -X POST "$API_BASE/chats/$SESSION_ID/messages" \
        -H "Content-Type: application/json" \
        -d "{\"message\": \"测试消息 $i\"}" &
done

wait
echo "✅ 并发请求完成"

echo ""
echo "🎉 测试完成！"
echo ""
echo "💡 提示："
echo "- 查看API文档: docs/api.md"
echo "- 学习指南: docs/learning-guide.md"
echo "- 如需配置OpenAI API Key，编辑 .env 文件"
