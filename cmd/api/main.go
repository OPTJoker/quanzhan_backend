package main

import (
	"log"

	"ai-development-learning/internal/config"
	"ai-development-learning/internal/database"
	"ai-development-learning/internal/handlers"
	"ai-development-learning/internal/services"

	"github.com/gin-gonic/gin"
)

func main() {
	// 加载配置
	cfg := config.Load()

	// 初始化数据库
	db, err := database.NewDatabase(cfg)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	// 初始化服务
	chatService := services.NewChatService(db, cfg)

	// 初始化处理器
	chatHandler := handlers.NewChatHandler(chatService)

	// 设置路由
	router := gin.Default()

	// 添加CORS中间件
	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// API路由
	api := router.Group("/api/v1")
	{
		// 聊天相关路由
		chats := api.Group("/chats")
		{
			chats.GET("", chatHandler.GetChats)                         // 获取所有聊天
			chats.POST("", chatHandler.CreateChat)                      // 创建聊天
			chats.POST("/:sessionId/messages", chatHandler.SendMessage) // 发送消息
			chats.GET("/:sessionId", chatHandler.GetChatHistory)        // 获取历史
		}
	}

	// 健康检查
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "AI Development Learning API is running",
		})
	})

	log.Printf("Server starting on port %s", cfg.AppPort)
	log.Fatal(router.Run(":" + cfg.AppPort))
}
