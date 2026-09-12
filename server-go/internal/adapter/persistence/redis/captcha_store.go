package redisx

import (
	"context"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"server-go/internal/usecase/port"
)

// CaptchaStore 键：captcha:{id}，TTL 由调用方传入（5min）。
type CaptchaStore struct{ c *redis.Client }

func NewCaptchaStore(c *redis.Client) port.CaptchaStore { return &CaptchaStore{c: c} }

func (s *CaptchaStore) Save(ctx context.Context, id, answer string, ttl time.Duration) error {
	return s.c.Set(ctx, "captcha:"+id, answer, ttl).Err()
}

// VerifyAndDelete 用 GETDEL（Redis ≥ 6.2）；图形验证码不区分大小写。
func (s *CaptchaStore) VerifyAndDelete(ctx context.Context, id, answer string) (bool, error) {
	v, err := s.c.GetDel(ctx, "captcha:"+id).Result()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return strings.EqualFold(v, answer), nil
}
