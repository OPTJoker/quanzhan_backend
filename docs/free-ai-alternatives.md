# 免费AI服务替代方案

OpenAI API是付费服务，这里提供一些免费的替代方案供学习和测试使用：

## 🆓 方案1: 使用内置模拟AI（推荐初学者）

**当前已集成**: 项目已经内置了智能模拟AI服务

**特点**:
- ✅ 完全免费，无需任何API密钥
- ✅ 支持基本对话功能
- ✅ 包含智能关键词回复
- ✅ 适合学习项目架构和前端开发

**使用方法**: 
无需任何配置，保持 `.env` 中的 `OPENAI_API_KEY=your_openai_api_key_here` 即可自动使用模拟AI。

## 🆓 方案2: 免费AI平台

### 2.1 Hugging Face (推荐)
- **网址**: https://huggingface.co/
- **免费额度**: 每月1000次免费调用
- **模型**: 支持多种开源模型
- **配置方法**:
```bash
# .env 文件配置
OPENAI_BASE_URL=https://api-inference.huggingface.co/models/microsoft/DialoGPT-medium
OPENAI_API_KEY=hf_your_huggingface_token
```

### 2.2 Cohere AI
- **网址**: https://cohere.ai/
- **免费额度**: 每月5M tokens
- **特点**: 专注于文本生成和理解

### 2.3 Anthropic Claude (有限免费)
- **网址**: https://www.anthropic.com/
- **免费额度**: 有限的免费体验

## 🆓 方案3: 本地AI模型

### 3.1 Ollama (强烈推荐)
**安装步骤**:
```bash
# macOS
brew install ollama

# 启动Ollama服务
ollama serve

# 下载并运行模型（约4GB）
ollama pull llama2:7b-chat
```

**配置项目**:
```bash
# .env 文件
OPENAI_BASE_URL=http://localhost:11434/v1
OPENAI_API_KEY=ollama  # 任意值
```

### 3.2 LocalAI
- 支持多种模型格式
- 兼容OpenAI API格式
- 完全本地运行

## 🆓 方案4: 国产AI平台

### 4.1 百度文心一言
- **网址**: https://yiyan.baidu.com/
- **免费额度**: 新用户有免费调用次数
- **API文档**: https://cloud.baidu.com/

### 4.2 阿里通义千问
- **网址**: https://tongyi.aliyun.com/
- **免费额度**: 每月一定免费调用量

### 4.3 腾讯混元大模型
- **网址**: https://cloud.tencent.com/product/hunyuan
- **免费额度**: 新用户免费体验

## 🛠 快速切换配置

我已经为你准备了智能切换机制：

1. **默认模拟AI**: 无需配置，直接体验
2. **真实AI**: 配置任意真实的API密钥即可自动切换

## 📝 配置示例

### 使用模拟AI（默认）
```bash
# .env
OPENAI_API_KEY=your_openai_api_key_here  # 保持默认值
```

### 使用Ollama本地模型
```bash
# .env
OPENAI_BASE_URL=http://localhost:11434/v1
OPENAI_API_KEY=ollama
```

### 使用Hugging Face
```bash
# .env  
OPENAI_BASE_URL=https://api-inference.huggingface.co/models/microsoft/DialoGPT-medium
OPENAI_API_KEY=hf_your_token_here
```

## 🎯 推荐学习路径

1. **阶段1**: 使用内置模拟AI学习项目架构 ✅ 
2. **阶段2**: 尝试Ollama本地模型体验真实AI
3. **阶段3**: 申请免费在线AI服务
4. **阶段4**: 考虑付费OpenAI API（生产环境）

## ⚡ 立即体验

现在你已经可以直接使用AI聊天功能了！

1. 启动后端服务: `make run`
2. 启动前端服务: `cd frontend && npm run dev`  
3. 访问 http://localhost:3000
4. 创建聊天并开始对话！

**系统会自动使用模拟AI，完全免费，功能完整！**
