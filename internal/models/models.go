package models

import (
	"time"

	"gorm.io/gorm"
)

// Chat 聊天会话
type Chat struct {
	ID        uint           `json:"id" gorm:"primarykey"`
	SessionID string         `json:"session_id" gorm:"uniqueIndex;size:36"`
	Title     string         `json:"title" gorm:"size:255"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
	Messages  []Message      `json:"messages" gorm:"foreignKey:ChatID;constraint:OnDelete:CASCADE"`
}

// Message 聊天消息
type Message struct {
	ID        uint           `json:"id" gorm:"primarykey"`
	ChatID    uint           `json:"chat_id"`
	Role      string         `json:"role" gorm:"size:20"` // user, assistant, system
	Content   string         `json:"content" gorm:"type:text"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// Document 文档知识库
type Document struct {
	ID          uint           `json:"id" gorm:"primarykey"`
	Title       string         `json:"title" gorm:"size:255"`
	Content     string         `json:"content" gorm:"type:longtext"`
	Source      string         `json:"source" gorm:"size:255"`
	Embedding   string         `json:"-" gorm:"type:json"` // 向量嵌入
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
	Collections []Collection   `json:"collections" gorm:"many2many:collection_documents;"`
}

// Collection 文档集合
type Collection struct {
	ID          uint           `json:"id" gorm:"primarykey"`
	Name        string         `json:"name" gorm:"size:100;uniqueIndex"`
	Description string         `json:"description" gorm:"size:500"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
	Documents   []Document     `json:"documents" gorm:"many2many:collection_documents;"`
}

// Agent AI代理
type Agent struct {
	ID           uint           `json:"id" gorm:"primarykey"`
	Name         string         `json:"name" gorm:"size:100"`
	Description  string         `json:"description" gorm:"size:500"`
	SystemPrompt string         `json:"system_prompt" gorm:"type:text"`
	Model        string         `json:"model" gorm:"size:50"`
	Temperature  float32        `json:"temperature" gorm:"default:0.7"`
	MaxTokens    int            `json:"max_tokens" gorm:"default:1000"`
	Tools        string         `json:"tools" gorm:"type:json"` // 可用工具列表
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
}
