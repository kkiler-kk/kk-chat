package redisx

import (
	"context"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"server-go/internal/domain"
	"server-go/internal/usecase/port"
)

// MsgLimitStore 键：msg:limit:{会话键}:{uid}。
// INCR 后若结果为 1（首次），EXPIRE 设 TTL；usecase 判断 count > limit 则拒绝。
type MsgLimitStore struct{ c *redis.Client }

func NewMsgLimitStore(c *redis.Client) port.MsgLimitStore { return &MsgLimitStore{c: c} }

func (s *MsgLimitStore) Incr(ctx context.Context, convID domain.ConversationID, uid int64, ttl time.Duration) (int64, error) {
	key := "msg:limit:" + convID.Value + ":" + strconv.FormatInt(uid, 10)
	n, err := s.c.Incr(ctx, key).Result()
	if err != nil {
		return 0, err
	}
	if n == 1 {
		s.c.Expire(ctx, key, ttl)
	}
	return n, nil
}
