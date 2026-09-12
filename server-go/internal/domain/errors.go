package domain

import "errors"

var (
	ErrUserNotFound        = errors.New("用户不存在")
	ErrEmailTaken          = errors.New("邮箱已被注册")
	ErrIdentityTaken       = errors.New("账号已被占用")
	ErrBadCredentials      = errors.New("账号或密码错误")
	ErrUserBanned          = errors.New("账号已被封禁")
	ErrCaptchaInvalid      = errors.New("图形验证码错误或已过期")
	ErrEmailCodeInvalid    = errors.New("邮箱验证码错误或已过期")
	ErrNotFriend           = errors.New("不是好友关系")
	ErrMsgLimitExceeded    = errors.New("非好友关系发送消息超过上限，请先添加好友")
	ErrGroupNotFound       = errors.New("群组不存在")
	ErrNotGroupMember      = errors.New("不是群成员")
	ErrTokenInvalid        = errors.New("token 无效或已过期")
	ErrInvalidConversation = errors.New("无效的会话标识")
)
