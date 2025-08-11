# AI应用开发实战学习项目

这是一个专注于AI应用开发的实战项目，帮助Go开发者快速掌握现代AI应用开发技能。

## 技术栈
- **后端**: Go (Gin框架)
- **AI框架**: LangChain (Python), OpenAI API
- **数据库**: MySQL (本地安装)
- **缓存**: Redis (本地安装)
- **前端**: React (可选)

## 学习路径

### 阶段1: AI应用开发基础
1. **Chat应用开发** - 基础聊天机器人
2. **RAG (检索增强生成)** - 知识库问答系统
3. **Agent开发** - 智能助手和工具调用

### 阶段2: 企业级AI应用
4. **多模态应用** - 文本、图像、语音处理
5. **AI工作流** - 复杂业务流程自动化
6. **向量数据库集成** - 高效相似性搜索

### 阶段3: 生产级部署
7. **模型管理** - 版本控制和A/B测试
8. **监控和日志** - AI应用性能监控
9. **微服务架构** - 可扩展的AI服务

## 项目结构
```
├── cmd/               # 应用入口点
├── internal/          # 内部包
│   ├── handlers/      # HTTP处理器
│   ├── services/      # 业务逻辑
│   ├── models/        # 数据模型
│   └── config/        # 配置管理
├── pkg/               # 公共包
├── scripts/           # Python AI脚本
├── docker/            # Docker配置
├── docs/              # 文档
└── examples/          # 示例代码
```

## 快速开始

### 1. 环境准备
- 确保已安装MySQL和Redis
- Node.js 18+ (用于前端)
- Go 1.21+ (用于后端)

### 2. 配置OpenAI API
⚠️ **重要**: 需要配置OpenAI API密钥才能使用AI功能

编辑 `.env` 文件，替换API密钥：
```bash
OPENAI_API_KEY=sk-your-actual-openai-api-key-here
```

详细配置说明请参考: [OpenAI配置指南](docs/openai-setup.md)

### 3. 启动服务

**启动后端服务:**
```bash
make run
```

**启动前端应用:**
```bash
cd frontend
npm install
npm run dev
```

### 4. 访问应用
- 前端界面: http://localhost:3000
- API接口: http://localhost:8080
- 健康检查: http://localhost:8080/health

### 5. 开始使用
1. 在前端界面点击"新建聊天"
2. 输入消息与AI助手对话
3. 体验完整的聊天功能
