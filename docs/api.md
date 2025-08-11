# API文档

## 基础信息

- **基础URL**: `http://localhost:8080/api/v1`
- **认证方式**: 暂无（开发环境）
- **响应格式**: JSON

## 通用响应格式

```json
{
  "success": true,
  "data": {},
  "error": null
}
```

## 聊天API

### 创建聊天会话

**POST** `/chats`

```json
{
  "title": "我的第一次对话"
}
```

**响应**:
```json
{
  "success": true,
  "data": {
    "id": 1,
    "session_id": "uuid-string",
    "title": "我的第一次对话",
    "created_at": "2024-01-01T10:00:00Z"
  }
}
```

### 发送消息

**POST** `/chats/{sessionId}/messages`

```json
{
  "message": "你好，你是谁？"
}
```

**响应**:
```json
{
  "success": true,
  "data": {
    "id": 2,
    "chat_id": 1,
    "role": "assistant",
    "content": "你好！我是AI助手，很高兴为您服务。",
    "created_at": "2024-01-01T10:00:00Z"
  }
}
```

### 获取聊天历史

**GET** `/chats/{sessionId}`

**响应**:
```json
{
  "success": true,
  "data": {
    "id": 1,
    "session_id": "uuid-string",
    "title": "我的第一次对话",
    "messages": [
      {
        "id": 1,
        "role": "user",
        "content": "你好，你是谁？",
        "created_at": "2024-01-01T10:00:00Z"
      },
      {
        "id": 2,
        "role": "assistant", 
        "content": "你好！我是AI助手，很高兴为您服务。",
        "created_at": "2024-01-01T10:00:00Z"
      }
    ]
  }
}
```

## Python AI服务API

### RAG查询

**POST** `http://localhost:8000/api/v1/rag/query`

```json
{
  "question": "什么是人工智能？",
  "collection_id": 1
}
```

**响应**:
```json
{
  "success": true,
  "data": {
    "answer": "人工智能是计算机科学的一个分支...",
    "sources": ["文档1片段", "文档2片段"],
    "confidence": 0.85
  }
}
```

### 上传文档

**POST** `http://localhost:8000/api/v1/rag/upload`

Form Data:
- `file`: 文档文件
- `title`: 文档标题
- `collection_id`: 集合ID（可选）

### Agent执行

**POST** `http://localhost:8000/api/v1/agent/execute`

```json
{
  "task": "帮我计算今天的销售额",
  "session_id": "agent-session-1",
  "tools": ["calculator", "database_query"]
}
```

**响应**:
```json
{
  "success": true,
  "data": {
    "response": "根据数据库查询，今天的销售额是15,680元",
    "session_id": "agent-session-1",
    "tools_used": ["database_query", "calculator"],
    "execution_time": 2.5
  }
}
```

## 错误码

| 错误码 | 说明 |
|--------|------|
| 400 | 请求参数错误 |
| 401 | 认证失败 |
| 404 | 资源不存在 |
| 429 | 请求频率限制 |
| 500 | 服务器内部错误 |

## 使用示例

### cURL示例

```bash
# 创建聊天
curl -X POST http://localhost:8080/api/v1/chats \
  -H "Content-Type: application/json" \
  -d '{"title": "测试对话"}'

# 发送消息
curl -X POST http://localhost:8080/api/v1/chats/session-id/messages \
  -H "Content-Type: application/json" \
  -d '{"message": "你好"}'
```

### JavaScript示例

```javascript
// 创建聊天会话
const createChat = async (title) => {
  const response = await fetch('http://localhost:8080/api/v1/chats', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ title }),
  });
  return response.json();
};

// 发送消息
const sendMessage = async (sessionId, message) => {
  const response = await fetch(`http://localhost:8080/api/v1/chats/${sessionId}/messages`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ message }),
  });
  return response.json();
};
```

### Go示例

```go
package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
)

type ChatRequest struct {
    Title string `json:"title"`
}

func createChat(title string) error {
    req := ChatRequest{Title: title}
    jsonData, _ := json.Marshal(req)
    
    resp, err := http.Post(
        "http://localhost:8080/api/v1/chats",
        "application/json",
        bytes.NewBuffer(jsonData),
    )
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    
    fmt.Println("Chat created successfully")
    return nil
}
```
