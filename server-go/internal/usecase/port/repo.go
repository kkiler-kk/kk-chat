package port

import (
	"context"
	"time"

	"server-go/internal/domain"
)

type UserRepo interface {
	Create(ctx context.Context, u *domain.User) error // 成功后回填 u.ID
	ByID(ctx context.Context, id int64) (*domain.User, error)
	ByIdentity(ctx context.Context, identity string) (*domain.User, error)
	ByEmail(ctx context.Context, email string) (*domain.User, error)
	ByIDs(ctx context.Context, ids []int64) ([]*domain.User, error)
	UpdateProfile(ctx context.Context, id int64, p UpdateProfile) error
	Search(ctx context.Context, keyword string, limit int) ([]*domain.User, error)
}

// UpdateProfile 局部更新：nil 字段不改。
type UpdateProfile struct {
	Name      *string
	Phone     *string
	Email     *string
	Avatar    *string
	Signature *string
	BirthDate *time.Time
}

type FriendRepo interface {
	Add(ctx context.Context, userID, friendID int64) error // 事务写双向两行
	Exists(ctx context.Context, a, b int64) (bool, error)
	ListFriends(ctx context.Context, userID int64) ([]*domain.User, error)
}

type GroupRepo interface {
	Create(ctx context.Context, g *domain.Group, memberIDs []int64) error // 事务：群+owner成员+初始成员
	ByID(ctx context.Context, id int64) (*domain.Group, error)
	Join(ctx context.Context, groupID, userID int64) error
	IsMember(ctx context.Context, groupID, userID int64) (bool, error)
	MemberIDs(ctx context.Context, groupID int64) ([]int64, error)
	MemberCount(ctx context.Context, groupID int64) (int, error)
	ListByUser(ctx context.Context, userID int64) ([]*domain.Group, error)
	SearchByName(ctx context.Context, keyword string, limit int) ([]*domain.Group, error)
}

type MessageRepo interface {
	Save(ctx context.Context, msg *domain.Message) error // 成功后回填 msg.ID
	// History 返回按 created_at 升序的一页消息：取 created_at < cursor 的最近 limit 条。
	// cursor 为零值表示从最新开始。
	History(ctx context.Context, convID domain.ConversationID, cursor time.Time, limit int) ([]domain.Message, error)
}
