package infra

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/kocar/aurelia/internal/config"
)

type RedisProvider struct {
	Client *redis.Client
}

func NewRedisProvider(cfg *config.AppConfig) (*RedisProvider, error) {
	slog.Info("inicializando o provider de Redis", "url", cfg.RedisURL)

	client := redis.NewClient(&redis.Options{
		Addr: cfg.RedisURL,
	})

	// Test connection
	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("redis ping: %w", err)
	}

	return &RedisProvider{
		Client: client,
	}, nil
}

func (r *RedisProvider) SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error) {
	return r.Client.SetNX(ctx, key, value, expiration).Result()
}
