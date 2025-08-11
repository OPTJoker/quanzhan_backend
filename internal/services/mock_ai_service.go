package services

import (
	"fmt"
	"strings"
	"time"

	"ai-development-learning/internal/models"
)

// MockAI 服务 - 用于测试和演示
func (s *ChatService) SendMessageWithMockAI(sessionID, userMessage string) (*models.Message, error) {
	// 获取聊天会话
	var chat models.Chat
	if err := s.db.DB.Where("session_id = ?", sessionID).
		Preload("Messages").First(&chat).Error; err != nil {
		return nil, fmt.Errorf("chat not found: %w", err)
	}

	// 保存用户消息
	userMsg := &models.Message{
		ChatID:  chat.ID,
		Role:    "user",
		Content: userMessage,
	}

	if err := s.db.DB.Create(userMsg).Error; err != nil {
		return nil, fmt.Errorf("failed to save user message: %w", err)
	}

	// 模拟AI回复
	mockResponses := []string{
		"这是一个模拟的AI回复。您的消息是：" + userMessage,
		"我理解您说的是：" + userMessage + "。这是一个演示回复。",
		"感谢您的消息！作为AI助手，我正在为您提供模拟回复。",
		"您好！我是模拟AI助手，正在处理您的消息：" + userMessage,
		"这是一个测试回复，用于演示聊天功能。您的输入：" + userMessage,
	}

	// 根据消息长度选择回复
	responseIndex := len(userMessage) % len(mockResponses)
	aiResponse := mockResponses[responseIndex]

	// 添加一些智能回复逻辑
	if contains(userMessage, []string{"你好", "hello", "hi"}) {
		aiResponse = "你好！我是AI助手，很高兴为您服务！有什么可以帮助您的吗？"
	} else if contains(userMessage, []string{"再见", "拜拜", "bye"}) {
		aiResponse = "再见！感谢使用AI聊天服务，期待下次为您服务！"
	} else if contains(userMessage, []string{"帮助", "help"}) {
		aiResponse = "我可以帮助您进行对话交流。这是一个演示版本，支持基本的聊天功能。"
	} else if contains(userMessage, []string{"时间", "几点"}) {
		aiResponse = fmt.Sprintf("当前时间是：%s", time.Now().Format("2006-01-02 15:04:05"))
	}

	// 保存AI回复
	aiMsg := &models.Message{
		ChatID:  chat.ID,
		Role:    "assistant",
		Content: aiResponse,
	}

	if err := s.db.DB.Create(aiMsg).Error; err != nil {
		return nil, fmt.Errorf("failed to save AI message: %w", err)
	}

	// 缓存对话
	s.cacheConversation(sessionID, userMessage, aiMsg.Content)

	return aiMsg, nil
}

// 辅助函数：检查消息是否包含关键词
func contains(message string, keywords []string) bool {
	message = strings.ToLower(message)
	for _, keyword := range keywords {
		if strings.Contains(message, strings.ToLower(keyword)) {
			return true
		}
	}
	return false
}
