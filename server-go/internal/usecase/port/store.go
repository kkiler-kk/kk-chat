package port

import (
	"context"
	"time"

	"server-go/internal/domain"
)

type TokenStore interface {
	Save(ctx context.Context, jti string, userID int64, ttl time.Duration) error
	Exists(ctx context.Context, jti string) (bool, error)
	Delete(ctx context.Context, jti string) error
}

// TokenIssuer 由 platform/token.Manager 实现；usecase 登录时签发 token。
type TokenIssuer interface {
	Issue(ctx context.Context, userID int64) (tokenStr string, jti string, err error)
}

type CaptchaStore interface {
	Save(ctx context.Context, id, answer string, ttl time.Duration) error
	VerifyAndDelete(ctx context.Context, id, answer string) (bool, error)
}

type EmailCodeStore interface {
	Save(ctx context.Context, email, code string, ttl time.Duration) error
	VerifyAndDelete(ctx context.Context, email, code string) (bool, error)
}

type PresenceStore interface {
	Online(ctx context.Context, userID int64, ttl time.Duration) error
	Offline(ctx context.Context, userID int64) error
	IsOnline(ctx context.Context, userID int64) (bool, error)
}

type ConversationSummary struct {
	Type            domain.ConversationType `json:"type"`
	LastContent     string                  `json:"last_message_content"`
	LastContentType domain.ContentType      `json:"last_content_type"`
	LastSenderID    int64                   `json:"last_sender_id"`
	LastSenderName  string                  `json:"last_sender_name"`
	LastTime        time.Time               `json:"last_time"`
}

type RecentChatStore interface {
	Touch(ctx context.Context, userID int64, convID domain.ConversationID, ts time.Time) error
	List(ctx context.Context, userID int64, limit int) ([]domain.ConversationID, error) // 按时间倒序
	SaveSummary(ctx context.Context, convID domain.ConversationID, s ConversationSummary) error
	Summary(ctx context.Context, convID domain.ConversationID) (*ConversationSummary, error)
}

type MsgLimitStore interface {
	Incr(ctx context.Context, convID domain.ConversationID, userID int64, ttl time.Duration) (int64, error)
}

type EmailSender interface {
	Send(ctx context.Context, to, subject, body string) error
}

type CaptchaGenerator interface {
	Generate(ctx context.Context) (id string, answer string, b64Image string, err error)
}

type Clock interface {
	Now() time.Time
}
