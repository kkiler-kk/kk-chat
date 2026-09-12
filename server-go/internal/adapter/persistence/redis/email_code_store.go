package redisx

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"server-go/internal/usecase/port"
)

// EmailCodeStore 键：verify:email:{email}，TTL 由调用方传入（10min）。
type EmailCodeStore struct{ c *redis.Client }

func NewEmailCodeStore(c *redis.Client) port.EmailCodeStore { return &EmailCodeStore{c: c} }

func (s *EmailCodeStore) Save(ctx context.Context, email, code string, ttl time.Duration) error {
	return s.c.Set(ctx, "verify:email:"+email, code, ttl).Err()
}

func (s *EmailCodeStore) VerifyAndDelete(ctx context.Context, email, code string) (bool, error) {
	v, err := s.c.GetDel(ctx, "verify:email:"+email).Result()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return v == code, nil
}
