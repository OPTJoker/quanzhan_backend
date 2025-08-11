package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"ai-development-learning/internal/config"
	"ai-development-learning/internal/database"
	"ai-development-learning/internal/models"

	"ai-development-learning/internal/intent"

	"github.com/google/uuid"
	"github.com/sashabaranov/go-openai"
	"gorm.io/gorm"
)

type ChatService struct {
	db               *database.Database
	client           *openai.Client
	config           *config.Config
	intentRecognizer *intent.IntentRecognizer
}

func NewChatService(db *database.Database, cfg *config.Config) *ChatService {
	client := openai.NewClient(cfg.OpenAIAPIKey)
	if cfg.OpenAIBaseURL != "" {
		config := openai.DefaultConfig(cfg.OpenAIAPIKey)
		config.BaseURL = cfg.OpenAIBaseURL
		client = openai.NewClientWithConfig(config)
	}

	return &ChatService{
		db:               db,
		client:           client,
		config:           cfg,
		intentRecognizer: intent.NewIntentRecognizer(cfg),
	}
}

// GetChats 获取所有聊天列表
func (s *ChatService) GetChats() ([]models.Chat, error) {
	var chats []models.Chat
	if err := s.db.DB.Order("updated_at DESC").Find(&chats).Error; err != nil {
		return nil, fmt.Errorf("failed to get chats: %w", err)
	}
	return chats, nil
}

// CreateChat 创建新的聊天会话
func (s *ChatService) CreateChat(title string) (*models.Chat, error) {
	sessionID := uuid.New().String()

	chat := &models.Chat{
		SessionID: sessionID,
		Title:     title,
	}

	if err := s.db.DB.Create(chat).Error; err != nil {
		return nil, fmt.Errorf("failed to create chat: %w", err)
	}

	return chat, nil
}

// SendMessage 发送消息并获取AI回复
func (s *ChatService) SendMessage(sessionID, userMessage string) (*models.Message, error) {
	fmt.Printf("[ChatService] 收到用户消息: sessionID=%s, message=%s\n", sessionID, userMessage)
	// ----------- 意图识别环节 -------------
	intent := "unknown"
	if s.intentRecognizer != nil {
		i, err := s.intentRecognizer.RecognizeIntent(userMessage)
		if err == nil && i != "" {
			intent = i
		}
	}
	fmt.Printf("[ChatService] 识别到意图: %s\n", intent)
	// 检查OpenAI API配置
	if s.config.OpenAIAPIKey == "" || s.config.OpenAIAPIKey == "your_openai_api_key_here" {
		// 如果没有配置真实的OpenAI API，使用模拟AI
		return s.SendMessageWithMockAI(sessionID, userMessage)
	}

	// 获取聊天会话
	var chat models.Chat
	if err := s.db.DB.Where("session_id = ?", sessionID).
		Preload("Messages").First(&chat).Error; err != nil {
		return nil, fmt.Errorf("chat not found: %w", err)
	}

	// 先不入库，只有AI回复成功才一起入库
	userMsg := &models.Message{
		ChatID:  chat.ID,
		Role:    "user",
		Content: userMessage,
	}

	// 构建对话历史，加入意图信息
	systemPrompt := "你是一个有用的AI助手，请用中文回复。当前用户意图为：" + intent
	fmt.Printf("[ChatService] 构建系统提示: %s\n", systemPrompt)
	messages := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: systemPrompt,
		},
	}

	// 添加历史消息（最近10条）
	historyCount := len(chat.Messages)
	start := 0
	if historyCount > 10 {
		start = historyCount - 10
	}

	for i := start; i < historyCount; i++ {
		msg := chat.Messages[i]
		role := openai.ChatMessageRoleUser
		if msg.Role == "assistant" {
			role = openai.ChatMessageRoleAssistant
		}
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    role,
			Content: msg.Content,
		})
	}

	// 添加当前用户消息
	messages = append(messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: userMessage,
	})

	// 调用OpenAI API
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	fmt.Printf("[ChatService] 调用大模型API...\n")
	resp, err := s.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model:       "qwen-turbo", //openai.GPT3Dot5Turbo,
		Messages:    messages,
		Temperature: 0.7,
		MaxTokens:   1000,
	})
	if err != nil {
		fmt.Printf("[ChatService] 大模型API调用失败: %v\n", err)
	}
	if len(resp.Choices) > 0 {
		fmt.Printf("[ChatService] 大模型回复: %s\n", resp.Choices[0].Message.Content)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get AI response: %w", err)
	}
	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("no response from AI")
	}
	// AI回复成功，用户消息和AI回复一起入库
	if err := s.db.DB.Create(userMsg).Error; err != nil {
		return nil, fmt.Errorf("failed to save user message: %w", err)
	}
	aiMsg := &models.Message{
		ChatID:  chat.ID,
		Role:    "assistant",
		Content: resp.Choices[0].Message.Content,
	}
	if err := s.db.DB.Create(aiMsg).Error; err != nil {
		return nil, fmt.Errorf("failed to save AI message: %w", err)
	}
	s.cacheConversation(sessionID, userMessage, aiMsg.Content)
	return aiMsg, nil
}

// GetChatHistory 获取聊天历史
func (s *ChatService) GetChatHistory(sessionID string) (*models.Chat, error) {
	var chat models.Chat
	err := s.db.DB.Where("session_id = ?", sessionID).
		Preload("Messages", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at ASC")
		}).First(&chat).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get chat history: %w", err)
	}

	return &chat, nil
}

// cacheConversation 缓存对话到Redis
func (s *ChatService) cacheConversation(sessionID, userMsg, aiMsg string) {
	ctx := context.Background()
	key := fmt.Sprintf("chat:%s:latest", sessionID)

	data := map[string]interface{}{
		"user_message": userMsg,
		"ai_message":   aiMsg,
		"timestamp":    time.Now().Unix(),
	}

	jsonData, _ := json.Marshal(data)
	s.db.Redis.Set(ctx, key, jsonData, 24*time.Hour)
}
