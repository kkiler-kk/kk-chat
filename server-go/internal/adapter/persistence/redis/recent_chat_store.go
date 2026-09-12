package redisx

import (
	"context"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"server-go/internal/domain"
	"server-go/internal/usecase/port"
)

// RecentChatStore
//
//	recent:{uid}  → ZSET，member=会话键字符串，score=最后消息 Unix 秒
//	conv:{会话键} → HASH：type / last_message_content / last_content_type /
//	               last_sender_id / last_sender_name / last_time(RFC3339)
type RecentChatStore struct{ c *redis.Client }

func NewRecentChatStore(c *redis.Client) port.RecentChatStore { return &RecentChatStore{c: c} }

func recentKey(uid int64) string { return "recent:" + strconv.FormatInt(uid, 10) }

func convKey(c domain.ConversationID) string { return "conv:" + c.Value }

func (s *RecentChatStore) Touch(ctx context.Context, uid int64, convID domain.ConversationID, ts time.Time) error {
	return s.c.ZAdd(ctx, recentKey(uid), redis.Z{Score: float64(ts.Unix()), Member: convID.Value}).Err()
}

func (s *RecentChatStore) List(ctx context.Context, uid int64, limit int) ([]domain.ConversationID, error) {
	vals, err := s.c.ZRevRange(ctx, recentKey(uid), 0, int64(limit-1)).Result()
	if err != nil {
		return nil, err
	}
	out := make([]domain.ConversationID, 0, len(vals))
	for _, v := range vals {
		if cid, err := domain.ParseConversationID(v); err == nil {
			out = append(out, cid)
		}
	}
	return out, nil
}

func (s *RecentChatStore) SaveSummary(ctx context.Context, convID domain.ConversationID, sum port.ConversationSummary) error {
	return s.c.HSet(ctx, convKey(convID), map[string]any{
		"type":                 convID.Type.String(),
		"last_message_content": sum.LastContent,
		"last_content_type":    string(sum.LastContentType),
		"last_sender_id":       sum.LastSenderID,
		"last_sender_name":     sum.LastSenderName,
		"last_time":            sum.LastTime.Format(time.RFC3339),
	}).Err()
}

func (s *RecentChatStore) Summary(ctx context.Context, convID domain.ConversationID) (*port.ConversationSummary, error) {
	m, err := s.c.HGetAll(ctx, convKey(convID)).Result()
	if err != nil || len(m) == 0 {
		return nil, err
	}
	t, _ := time.Parse(time.RFC3339, m["last_time"])
	ct, _ := domain.ParseConversationType(m["type"])
	senderID, _ := strconv.ParseInt(m["last_sender_id"], 10, 64)
	return &port.ConversationSummary{
		Type:            ct,
		LastContent:     m["last_message_content"],
		LastContentType: domain.ContentType(m["last_content_type"]),
		LastSenderID:    senderID,
		LastSenderName:  m["last_sender_name"],
		LastTime:        t,
	}, nil
}
