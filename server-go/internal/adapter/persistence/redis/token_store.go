package redisx

import (
	"context"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"server-go/internal/usecase/port"
)

// TokenStore 键：token:{jti}，值：user_id 十进制字符串。
type TokenStore struct {
	c   *redis.Client
	ttl time.Duration
}

func NewTokenStore(c *redis.Client, ttl time.Duration) port.TokenStore {
	return &TokenStore{c: c, ttl: ttl}
}

func (s *TokenStore) Save(ctx context.Context, jti string, userID int64, ttl time.Duration) error {
	use := s.ttl
	if ttl > 0 {
		use = ttl
	}
	return s.c.Set(ctx, "token:"+jti, strconv.FormatInt(userID, 10), use).Err()
}

func (s *TokenStore) Exists(ctx context.Context, jti string) (bool, error) {
	n, err := s.c.Exists(ctx, "token:"+jti).Result()
	return n > 0, err
}

func (s *TokenStore) Delete(ctx context.Context, jti string) error {
	return s.c.Del(ctx, "token:"+jti).Err()
}
