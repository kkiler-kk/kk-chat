package redisx

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

// NewClient 创建并探活 Redis 连接。
func NewClient(ctx context.Context, addr, password string, db int) (*redis.Client, error) {
	c := redis.NewClient(&redis.Options{Addr: addr, Password: password, DB: db})
	if err := c.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("连接 redis 失败: %w", err)
	}
	return c, nil
}
