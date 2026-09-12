package port

import (
	"context"
	"time"

	"server-go/internal/domain"
)

type AuthUseCase interface {
	Register(ctx context.Context, in RegisterInput) error
	Login(ctx context.Context, in LoginInput) (*LoginOutput, error)
	Logout(ctx context.Context, jti string, userID int64) error
	GenerateCaptcha(ctx context.Context) (id, b64Image string, err error)
	SendEmailCode(ctx context.Context, email string) error
}

type UserUseCase interface {
	Me(ctx context.Context, userID int64) (*UserInfo, error)
	UpdateMe(ctx context.Context, userID int64, in UpdateUserInput) error
	Detail(ctx context.Context, viewerID, targetID int64) (*UserDetailOutput, error)
	Search(ctx context.Context, viewerID int64, keyword string) ([]UserSearchItem, error)
}

type FriendUseCase interface {
	Add(ctx context.Context, userID, targetID int64) error
	List(ctx context.Context, userID int64) ([]FriendItem, error)
}

type GroupUseCase interface {
	Create(ctx context.Context, ownerID int64, in CreateGroupInput) (*GroupItem, error)
	Join(ctx context.Context, groupID, userID int64) error
	List(ctx context.Context, userID int64) ([]GroupItem, error)
	Search(ctx context.Context, keyword string) ([]GroupItem, error)
}

type ChatUseCase interface {
	SendMessage(ctx context.Context, senderID int64, in SendMessageInput) (*OutgoingMessage, error)
	History(ctx context.Context, userID int64, convID domain.ConversationID, cursor time.Time, limit int) ([]OutgoingMessage, error)
	RecentConversations(ctx context.Context, userID int64) ([]RecentConversation, error)
}

type PresenceUseCase interface {
	OnConnect(ctx context.Context, userID int64) error    // 写 presence + 广播上线
	Heartbeat(ctx context.Context, userID int64) error    // 续期 presence
	OnDisconnect(ctx context.Context, userID int64) error // 删 presence + 广播下线
}
