# 🚀 快速开始指南

恭喜！你已经成功设置了AI应用开发学习环境。现在可以开始学习AI应用开发了！

## ✅ 当前状态

- ✅ MySQL数据库已配置并运行
- ✅ Redis缓存服务已启动  
- ✅ Go API服务正在运行 (http://localhost:8080)
- ✅ 基础功能测试通过
- ⚠️  OpenAI API Key需要配置（用于AI功能）

## 🎯 学习路径

### 阶段1：基础AI应用开发 (1-2周)

#### 第1天：熟悉项目结构
```bash
# 查看项目结构
tree -I 'node_modules|*.log'

# 运行基础测试
./scripts/test_basic.sh

# 查看数据库
mysql -u ai_user -pai_pass -D ai_development
```

#### 第2-3天：聊天机器人开发
1. **配置OpenAI API Key**（如果有的话）：
   ```bash
   # 编辑 .env 文件
   nano .env
   # 设置: OPENAI_API_KEY=你的真实API密钥
   ```

2. **测试聊天功能**：
   ```bash
   # 创建聊天会话
   curl -X POST http://localhost:8080/api/v1/chats \
     -H "Content-Type: application/json" \
     -d '{"title": "我的第一个AI聊天"}'

   # 发送消息（需要API Key）
   curl -X POST http://localhost:8080/api/v1/chats/SESSION_ID/messages \
     -H "Content-Type: application/json" \
     -d '{"message": "你好，请介绍一下自己"}'
   ```

3. **研究代码**：
   - `cmd/api/main.go` - 应用入口
   - `internal/services/chat_service.go` - 聊天服务逻辑
   - `internal/handlers/chat_handler.go` - HTTP处理器

#### 第4-7天：RAG知识库系统
```bash
# 安装Python依赖
pip install -r scripts/requirements.txt

# 运行RAG示例
python scripts/rag_example.py

# 启动Python AI服务
python scripts/api_server.py
```

#### 第8-14天：AI Agent开发
```bash
# 研究Agent示例
python scripts/agent_example.py

# 创建自定义工具
# 编辑 scripts/agent_example.py
```

### 阶段2：高级应用开发 (2-3周)

#### 多模态处理
- 图像识别和描述
- 语音转文字
- 文档解析和问答

#### 性能优化
- 缓存策略优化
- 并发处理
- API限流

### 阶段3：生产级部署 (1-2周)

#### Docker化部署
```bash
# 可选：如果以后需要容器化部署
# 当前使用本地MySQL和Redis，更加轻量
```

## 🛠️ 常用命令

```bash
# 项目管理
make help              # 查看所有可用命令
make quickstart        # 一键启动（适合重新开始）
make build             # 构建Go应用
make run               # 运行Go API
make test              # 运行测试
make clean             # 清理项目

# 服务管理
brew services start redis    # 启动Redis
brew services stop redis     # 停止Redis
go run cmd/api/main.go       # 启动Go API
python scripts/api_server.py # 启动Python AI服务

# 数据库操作
mysql -u ai_user -pai_pass -D ai_development  # 连接数据库
redis-cli                                     # 连接Redis
```

## 📚 学习资源

### 核心文档
- [学习指南](docs/learning-guide.md) - 详细的学习路径
- [API文档](docs/api.md) - API使用说明

### 代码示例
- `cmd/api/main.go` - Go API服务
- `scripts/rag_example.py` - RAG知识库系统
- `scripts/agent_example.py` - AI Agent开发
- `scripts/api_server.py` - Python AI服务

### 在线资源
- [LangChain文档](https://python.langchain.com/)
- [OpenAI API文档](https://platform.openai.com/docs)
- [Gin框架文档](https://gin-gonic.com/)

## 🎯 实战练习建议

### 初级练习
1. **修改聊天机器人** - 改变系统提示，让AI扮演不同角色
2. **添加新的API端点** - 创建用户管理功能
3. **数据库查询** - 实现聊天记录搜索功能

### 中级练习
1. **构建知识库** - 上传PDF文档，实现问答功能
2. **自定义Agent工具** - 创建天气查询、计算器等工具
3. **前端界面** - 用React/Vue构建聊天界面

### 高级练习
1. **多用户支持** - 添加用户认证和权限管理
2. **实时通信** - 使用WebSocket实现实时聊天
3. **微服务架构** - 将功能拆分为独立的微服务

## 🔧 故障排除

### 常见问题

1. **MySQL连接失败**
   ```bash
   # 检查MySQL服务
   brew services list | grep mysql
   
   # 重启MySQL
   brew services restart mysql
   ```

2. **Redis连接失败**
   ```bash
   # 检查Redis服务
   redis-cli ping
   
   # 重启Redis
   brew services restart redis
   ```

3. **Go依赖问题**
   ```bash
   # 清理并重新安装依赖
   go clean -modcache
   go mod download
   go mod tidy
   ```

4. **OpenAI API错误**
   - 检查API Key是否正确设置
   - 确认账户有足够的credits
   - 检查网络连接

## 🎉 开始你的AI开发之旅！

现在一切都已准备就绪，你可以：

1. **立即开始** - 运行 `./scripts/test_basic.sh` 验证环境
2. **深入学习** - 查看 `docs/learning-guide.md` 
3. **动手实践** - 修改代码，添加新功能
4. **寻求帮助** - 查看文档或在线资源

记住：**最好的学习方式就是动手实践！** 🚀

Happy Coding! 🤖
