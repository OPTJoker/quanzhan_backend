package database

import (
	"context"
	"time"

	"ai-development-learning/internal/config"
	"ai-development-learning/internal/models"

	"github.com/go-redis/redis/v8"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Database struct {
	DB    *gorm.DB
	Redis *redis.Client
}

func NewDatabase(cfg *config.Config) (*Database, error) {
	// 初始化 MySQL
	db, err := gorm.Open(mysql.Open(cfg.MySQLDSN), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// 数据库表字段 自动同步
	if err := db.AutoMigrate(
		&models.Chat{},
		&models.Message{},
		&models.Document{},
		&models.Collection{},
		&models.Agent{},
	); err != nil {
		return nil, err
	}

	// 初始化 Redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})

	// 测试 Redis 连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return &Database{
		DB:    db,
		Redis: rdb,
	}, nil
}

func (d *Database) Close() error {
	sqlDB, err := d.DB.DB()
	if err != nil {
		return err
	}

	if err := sqlDB.Close(); err != nil {
		return err
	}

	return d.Redis.Close()
}
