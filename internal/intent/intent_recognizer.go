package intent

import (
	"context"
	"fmt"
	"time"

	"ai-development-learning/internal/config"

	"github.com/sashabaranov/go-openai"
)

type IntentRecognizer struct {
	client *openai.Client
	config *config.Config
}

func NewIntentRecognizer(cfg *config.Config) *IntentRecognizer {
	client := openai.NewClient(cfg.OpenAIAPIKey)
	if cfg.OpenAIBaseURL != "" {
		c := openai.DefaultConfig(cfg.OpenAIAPIKey)
		c.BaseURL = cfg.OpenAIBaseURL
		client = openai.NewClientWithConfig(c)
	}
	return &IntentRecognizer{
		client: client,
		config: cfg,
	}
}

// RecognizeIntent 通过大模型识别用户意图
func (r *IntentRecognizer) RecognizeIntent(userMessage string) (string, error) {
	fmt.Printf("[IntentRecognizer] 开始识别意图，用户输入: %s\n", userMessage)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	prompt := "请判断用户输入的意图类别，只返回类别名称。例如：问答、闲聊、任务、指令、其他。用户输入：" + userMessage
	req := openai.ChatCompletionRequest{
		Model: "qwen-turbo",
		Messages: []openai.ChatCompletionMessage{
			{
				Role: openai.ChatMessageRoleSystem,
				Content: `你是一个意图识别助手。分析用户输入并返回对应的意图类别。
				可能的意图类别：
				- 问答：询问问题
				- 闲聊：日常聊天
				- 任务：请求执行某个操作
				- 指令：明确的命令
				- 其他：无法分类的内容
				
				只返回意图类别名称，不要任何解释。`,
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: prompt,
			},
		},
		Temperature: 0,
		MaxTokens:   20,
	}
	resp, err := r.client.CreateChatCompletion(ctx, req)
	if err != nil || len(resp.Choices) == 0 {
		fmt.Printf("[IntentRecognizer] 识别失败: %v\n", err)
		return "unknown", fmt.Errorf("intent recognition failed: %w", err)
	}
	intent := resp.Choices[0].Message.Content
	fmt.Printf("[IntentRecognizer] 识别结果: %s\n", intent)
	return intent, nil
}
