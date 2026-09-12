package usecase

import (
	"context"
	"errors"

	"server-go/internal/domain"
	"server-go/internal/usecase/apperror"
	"server-go/internal/usecase/port"
)

type FriendDeps struct {
	Users    port.UserRepo
	Friends  port.FriendRepo
	Presence port.PresenceStore
}

type friendUseCase struct{ d FriendDeps }

func NewFriend(d FriendDeps) port.FriendUseCase { return &friendUseCase{d: d} }

func (u *friendUseCase) Add(ctx context.Context, userID, targetID int64) error {
	if userID == targetID {
		return apperror.New(apperror.CodeInvalidParam, "不能添加自己为好友")
	}
	if _, err := u.d.Users.ByID(ctx, targetID); err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return apperror.New(apperror.CodeUserNotFound, domain.ErrUserNotFound.Error())
		}
		return apperror.Wrap(apperror.CodeInternal, "添加好友失败", err)
	}
	exists, err := u.d.Friends.Exists(ctx, userID, targetID)
	if err != nil {
		return apperror.Wrap(apperror.CodeInternal, "添加好友失败", err)
	}
	if exists {
		return apperror.New(apperror.CodeAlreadyFriend, "已经是好友了")
	}
	if err := u.d.Friends.Add(ctx, userID, targetID); err != nil {
		return apperror.Wrap(apperror.CodeInternal, "添加好友失败", err)
	}
	return nil
}

func (u *friendUseCase) List(ctx context.Context, userID int64) ([]port.FriendItem, error) {
	friends, err := u.d.Friends.ListFriends(ctx, userID)
	if err != nil {
		return nil, apperror.Wrap(apperror.CodeInternal, "查询好友失败", err)
	}
	out := make([]port.FriendItem, 0, len(friends))
	for _, f := range friends {
		online, err := u.d.Presence.IsOnline(ctx, f.ID)
		if err != nil {
			online = false // 在线状态查询失败降级为离线，不打断列表
		}
		out = append(out, port.FriendItem{ID: f.ID, Name: f.Name, Avatar: f.Avatar, Online: online})
	}
	return out, nil
}

var _ port.FriendUseCase = (*friendUseCase)(nil)
