package redisx

import (
	"context"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"server-go/internal/usecase/port"
)

// PresenceStore 键：presence:{uid}，值 "1"，TTL 到期即视为离线。
type PresenceStore struct {
	c   *redis.Client
	ttl time.Duration
}

func NewPresenceStore(c *redis.Client, ttl time.Duration) port.PresenceStore {
	return &PresenceStore{c: c, ttl: ttl}
}

func presenceKey(uid int64) string { return "presence:" + strconv.FormatInt(uid, 10) }

// Online 写入/续期同一操作。
func (s *PresenceStore) Online(ctx context.Context, userID int64, ttl time.Duration) error {
	use := s.ttl
	if ttl > 0 {
		use = ttl
	}
	return s.c.Set(ctx, presenceKey(userID), "1", use).Err()
}

func (s *PresenceStore) Offline(ctx context.Context, userID int64) error {
	return s.c.Del(ctx, presenceKey(userID)).Err()
}

func (s *PresenceStore) IsOnline(ctx context.Context, userID int64) (bool, error) {
	n, err := s.c.Exists(ctx, presenceKey(userID)).Result()
	return n > 0, err
}
