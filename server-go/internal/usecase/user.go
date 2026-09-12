package usecase

import (
	"context"
	"errors"
	"strings"

	"server-go/internal/domain"
	"server-go/internal/usecase/apperror"
	"server-go/internal/usecase/port"
)

type UserDeps struct {
	Users   port.UserRepo
	Friends port.FriendRepo
	Codes   port.EmailCodeStore
}

type userUseCase struct{ d UserDeps }

func NewUser(d UserDeps) port.UserUseCase { return &userUseCase{d: d} }

func (u *userUseCase) Me(ctx context.Context, userID int64) (*port.UserInfo, error) {
	user, err := u.d.Users.ByID(ctx, userID)
	if err != nil {
		return nil, u.mapUserErr(err, "查询用户失败")
	}
	return ToUserInfo(user), nil
}

func (u *userUseCase) UpdateMe(ctx context.Context, userID int64, in port.UpdateUserInput) error {
	if in.Email != nil && *in.Email != "" {
		if in.EmailCode == "" {
			return apperror.New(apperror.CodeInvalidParam, "修改邮箱需要验证码")
		}
		ok, err := u.d.Codes.VerifyAndDelete(ctx, *in.Email, in.EmailCode)
		if err != nil {
			return apperror.Wrap(apperror.CodeInternal, "更新资料失败", err)
		}
		if !ok {
			return apperror.New(apperror.CodeEmailCodeInvalid, domain.ErrEmailCodeInvalid.Error())
		}
		if existing, err := u.d.Users.ByEmail(ctx, *in.Email); err == nil && existing.ID != userID {
			return apperror.New(apperror.CodeEmailTaken, domain.ErrEmailTaken.Error())
		} else if err != nil && !errors.Is(err, domain.ErrUserNotFound) {
			return apperror.Wrap(apperror.CodeInternal, "更新资料失败", err)
		}
	}
	p := port.UpdateProfile{
		Name:      in.Name,
		Phone:     in.Phone,
		Email:     in.Email,
		Avatar:    in.Avatar,
		Signature: in.Signature,
		BirthDate: in.BirthDate,
	}
	if err := u.d.Users.UpdateProfile(ctx, userID, p); err != nil {
		return u.mapUserErr(err, "更新资料失败")
	}
	return nil
}

func (u *userUseCase) Detail(ctx context.Context, viewerID, targetID int64) (*port.UserDetailOutput, error) {
	user, err := u.d.Users.ByID(ctx, targetID)
	if err != nil {
		return nil, u.mapUserErr(err, "查询用户失败")
	}
	out := &port.UserDetailOutput{UserInfo: ToPublicUserInfo(user)}
	out.IsSelf = viewerID == targetID
	if out.IsSelf {
		out.UserInfo = ToUserInfo(user) // 本人可见隐私字段
		out.IsFriend = true
		return out, nil
	}
	if viewerID > 0 {
		isFriend, err := u.d.Friends.Exists(ctx, viewerID, targetID)
		if err != nil {
			return nil, apperror.Wrap(apperror.CodeInternal, "查询用户失败", err)
		}
		out.IsFriend = isFriend
		if isFriend {
			out.UserInfo = ToUserInfo(user) // 好友可见隐私字段
		}
	}
	return out, nil
}

func (u *userUseCase) Search(ctx context.Context, viewerID int64, keyword string) ([]port.UserSearchItem, error) {
	keyword = strings.TrimSpace(keyword)
	if keyword == "" {
		return []port.UserSearchItem{}, nil
	}
	users, err := u.d.Users.Search(ctx, keyword, 20)
	if err != nil {
		return nil, apperror.Wrap(apperror.CodeInternal, "搜索用户失败", err)
	}
	out := make([]port.UserSearchItem, 0, len(users))
	for _, user := range users {
		item := port.UserSearchItem{
			ID: user.ID, Identity: user.Identity, Name: user.Name, Avatar: user.Avatar,
		}
		if viewerID > 0 && viewerID != user.ID {
			isFriend, err := u.d.Friends.Exists(ctx, viewerID, user.ID)
			if err == nil {
				item.IsFriend = isFriend
			}
		}
		out = append(out, item)
	}
	return out, nil
}

func (u *userUseCase) mapUserErr(err error, msg string) error {
	if errors.Is(err, domain.ErrUserNotFound) {
		return apperror.New(apperror.CodeUserNotFound, domain.ErrUserNotFound.Error())
	}
	if errors.Is(err, domain.ErrEmailTaken) {
		return apperror.New(apperror.CodeEmailTaken, domain.ErrEmailTaken.Error())
	}
	return apperror.Wrap(apperror.CodeInternal, msg, err)
}

var _ port.UserUseCase = (*userUseCase)(nil)
