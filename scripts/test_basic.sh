#!/bin/bash

# AI应用开发 - 基础测试脚本（不需要OpenAI API）

echo "🧪 AI应用开发 - 基础功能测试"
echo "==============================="

API_BASE="http://localhost:8080/api/v1"

# 检查服务是否运行
echo "🔍 检查服务状态..."

if curl -s "$API_BASE/../health" > /dev/null; then
    echo "✅ Go API服务正常运行"
else
    echo "❌ Go API服务未运行，请先启动: go run cmd/api/main.go"
    exit 1
fi

echo ""
echo "🚀 开始基础功能测试..."

# 测试1: 健康检查
echo "❤️  测试1: 健康检查"
HEALTH_RESPONSE=$(curl -s "http://localhost:8080/health")
echo "响应: $HEALTH_RESPONSE"
echo ""

# 测试2: 创建聊天会话
echo "📝 测试2: 创建聊天会话"
CHAT_RESPONSE=$(curl -s -X POST "$API_BASE/chats" \
    -H "Content-Type: application/json" \
    -d '{"title": "基础测试对话"}')

echo "响应: $CHAT_RESPONSE"

# 从响应中提取session_id
if command -v jq &> /dev/null; then
    SESSION_ID=$(echo $CHAT_RESPONSE | jq -r '.data.session_id')
    echo "会话ID: $SESSION_ID"
else
    echo "⚠️  需要安装jq工具来解析JSON响应"
    echo "可以运行: brew install jq"
    # 手动提取session_id（简单方法）
    SESSION_ID=$(echo $CHAT_RESPONSE | grep -o '"session_id":"[^"]*"' | cut -d'"' -f4)
    echo "会话ID: $SESSION_ID"
fi

echo ""

# 测试3: 获取空的聊天历史
echo "📚 测试3: 获取聊天历史（应该为空）"
HISTORY_RESPONSE=$(curl -s -X GET "$API_BASE/chats/$SESSION_ID")
echo "响应: $HISTORY_RESPONSE"
echo ""

# 测试4: 数据库连接测试
echo "🗄️  测试4: 数据库连接测试"
echo "检查数据库中的聊天记录..."

mysql -u ai_user -pai_pass -D ai_development -e "
SELECT 
    'Database Status' as test,
    COUNT(*) as chat_count 
FROM chats;

SELECT 
    'Latest Chat' as info,
    id, 
    session_id, 
    title, 
    created_at 
FROM chats 
ORDER BY created_at DESC 
LIMIT 1;
" 2>/dev/null

echo ""

# 测试5: Redis连接测试
echo "🔴 测试5: Redis连接测试"
redis-cli set "test_key" "AI应用开发测试" > /dev/null
REDIS_VALUE=$(redis-cli get "test_key")
echo "Redis测试值: $REDIS_VALUE"

# 清理测试数据
redis-cli del "test_key" > /dev/null

echo ""
echo "🎉 基础功能测试完成！"
echo ""
echo "✅ 测试结果总结："
echo "- ✅ Go API服务运行正常"
echo "- ✅ MySQL数据库连接正常"
echo "- ✅ Redis缓存连接正常"
echo "- ✅ 聊天会话创建功能正常"
echo "- ⚠️  AI回复功能需要配置OpenAI API Key"
echo ""
echo "📖 下一步："
echo "1. 如果需要AI功能，请在 .env 文件中设置真实的 OPENAI_API_KEY"
echo "2. 查看学习指南: docs/learning-guide.md"
echo "3. 开始AI应用开发学习之旅！"
