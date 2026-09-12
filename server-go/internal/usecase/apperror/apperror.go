package apperror

import (
	"errors"
	"fmt"
)

// 业务错误码（与 spec §3.2 一致）
const (
	CodeOK                  = 0
	CodeInvalidParam        = 10001
	CodeInternal            = 10002
	CodeCaptchaInvalid      = 20001
	CodeBadCredentials      = 20002
	CodeTokenInvalid        = 20003
	CodeUserBanned          = 20004
	CodeEmailTaken          = 30001
	CodeIdentityTaken       = 30002
	CodeUserNotFound        = 30003
	CodeEmailCodeInvalid    = 30004
	CodeAlreadyFriend       = 40001
	CodeNotFriendLimit      = 40002
	CodeGroupNotFound       = 50001
	CodeNotGroupMember      = 50002
	CodeConversationInvalid = 60001
)

// Error 是应用层统一错误类型：Code 业务码，Message 用户可读消息，Err 内部错误链。
type Error struct {
	Code    int
	Message string
	Err     error
}

func (e *Error) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%d] %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("[%d] %s", e.Code, e.Message)
}

func (e *Error) Unwrap() error { return e.Err }

func New(code int, message string) *Error {
	return &Error{Code: code, Message: message}
}

func Wrap(code int, message string, err error) *Error {
	return &Error{Code: code, Message: message, Err: err}
}

// From 提取错误链中的 *Error；不是则返回 nil,false。
func From(err error) (*Error, bool) {
	var ae *Error
	if errors.As(err, &ae) {
		return ae, true
	}
	return nil, false
}
