# OpenAI API 配置说明

## 获取 OpenAI API 密钥

要使用AI聊天功能，您需要配置OpenAI API密钥。

### 步骤1：获取API密钥

1. 访问 [OpenAI Platform](https://platform.openai.com/)
2. 注册或登录您的账户
3. 进入 [API Keys](https://platform.openai.com/api-keys) 页面
4. 点击 "Create new secret key" 创建新的API密钥
5. 复制生成的API密钥（以 `sk-` 开头）

### 步骤2：配置环境变量

编辑项目根目录的 `.env` 文件：

```bash
# 将 your_openai_api_key_here 替换为您的真实API密钥
OPENAI_API_KEY=sk-your-actual-api-key-here

# 其他配置保持不变
APP_PORT=8080
MYSQL_DSN=ai_user:ai_pass@tcp(localhost:3306)/ai_development?charset=utf8mb4&parseTime=True&loc=Local
REDIS_ADDR=localhost:6379
REDIS_PASSWORD=
REDIS_DB=0
OPENAI_BASE_URL=https://api.openai.com/v1
```

### 步骤3：重启服务

配置完成后，重启后端服务：

```bash
cd /Users/sharon/me/code/ai
make run
```

## 替代方案

如果您没有OpenAI API密钥，可以考虑以下替代方案：

### 1. 使用其他API服务

修改 `.env` 文件中的 `OPENAI_BASE_URL`：

```bash
# 例如使用国内的API服务
OPENAI_BASE_URL=https://your-alternative-api.com/v1
OPENAI_API_KEY=your-alternative-api-key
```

### 2. 模拟AI回复

在测试阶段，您可以临时修改代码使用模拟回复。

### 3. 本地大模型

可以配置本地运行的大模型服务（如Ollama、LocalAI等）。

## 注意事项

- ⚠️ API密钥具有付费额度，请妥善保管
- ⚠️ 不要将API密钥提交到版本控制系统
- ⚠️ 定期检查API使用量，避免意外产生费用
- ⚠️ 生产环境建议使用环境变量而不是.env文件

## 常见问题

### Q: 提示"OpenAI API key not configured"

A: 检查.env文件中的OPENAI_API_KEY是否设置为真实的API密钥。

### Q: API调用失败

A: 
1. 检查API密钥是否正确
2. 检查网络连接
3. 检查API配额是否充足
4. 查看后端日志获取详细错误信息

### Q: 想要更换AI模型

A: 在 `internal/services/chat_service.go` 中修改 `Model` 字段：

```go
resp, err := s.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
    Model:       openai.GPT4,  // 改为 GPT-4
    Messages:    messages,
    Temperature: 0.7,
    MaxTokens:   1000,
})
```

支持的模型包括：
- `openai.GPT3Dot5Turbo` - GPT-3.5-turbo（默认，性价比高）
- `openai.GPT4` - GPT-4（更强但更贵）
- `openai.GPT4Turbo` - GPT-4-turbo（最新版本）
