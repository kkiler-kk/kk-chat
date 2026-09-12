package usecase

import (
	"context"
	"time"

	"server-go/internal/usecase/apperror"
	"server-go/internal/usecase/port"
)

type PresenceDeps struct {
	Presence port.PresenceStore
	Friends  port.FriendRepo
	Notifier port.Notifier
	TTL      time.Duration
}

type presenceUseCase struct{ d PresenceDeps }

func NewPresence(d PresenceDeps) port.PresenceUseCase { return &presenceUseCase{d: d} }

func (u *presenceUseCase) broadcast(ctx context.Context, userID int64, online bool) {
	friendUsers, err := u.d.Friends.ListFriends(ctx, userID)
	if err != nil || len(friendUsers) == 0 {
		return
	}
	ids := make([]int64, 0, len(friendUsers))
	for _, f := range friendUsers {
		ids = append(ids, f.ID)
	}
	u.d.Notifier.ToUsers(ctx, ids, port.NotifierEvent{
		Event: "presence.changed",
		Data:  map[string]any{"user_id": userID, "online": online},
	})
}

func (u *presenceUseCase) OnConnect(ctx context.Context, userID int64) error {
	if err := u.d.Presence.Online(ctx, userID, u.d.TTL); err != nil {
		return apperror.Wrap(apperror.CodeInternal, "上线失败", err)
	}
	u.broadcast(ctx, userID, true)
	return nil
}

func (u *presenceUseCase) Heartbeat(ctx context.Context, userID int64) error {
	if err := u.d.Presence.Online(ctx, userID, u.d.TTL); err != nil {
		return apperror.Wrap(apperror.CodeInternal, "心跳失败", err)
	}
	return nil
}

func (u *presenceUseCase) OnDisconnect(ctx context.Context, userID int64) error {
	if err := u.d.Presence.Offline(ctx, userID); err != nil {
		return apperror.Wrap(apperror.CodeInternal, "下线失败", err)
	}
	u.broadcast(ctx, userID, false)
	return nil
}

var _ port.PresenceUseCase = (*presenceUseCase)(nil)
