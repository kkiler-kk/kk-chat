package middleware

import (
	"errors"
	"log/slog"
	"net/http"
	"runtime/debug"

	"github.com/labstack/echo/v4"
	"server-go/internal/adapter/http/dto"
	"server-go/internal/domain"
	"server-go/internal/usecase/apperror"
)

// codeToStatus 业务码 → HTTP 状态码（spec §3.2）。
var codeToStatus = map[int]int{
	apperror.CodeInvalidParam:        http.StatusBadRequest,
	apperror.CodeInternal:            http.StatusInternalServerError,
	apperror.CodeCaptchaInvalid:      http.StatusBadRequest,
	apperror.CodeBadCredentials:      http.StatusUnauthorized,
	apperror.CodeTokenInvalid:        http.StatusUnauthorized,
	apperror.CodeUserBanned:          http.StatusForbidden,
	apperror.CodeEmailTaken:          http.StatusConflict,
	apperror.CodeIdentityTaken:       http.StatusConflict,
	apperror.CodeUserNotFound:        http.StatusNotFound,
	apperror.CodeEmailCodeInvalid:    http.StatusBadRequest,
	apperror.CodeAlreadyFriend:       http.StatusConflict,
	apperror.CodeNotFriendLimit:      http.StatusTooManyRequests,
	apperror.CodeGroupNotFound:       http.StatusNotFound,
	apperror.CodeNotGroupMember:      http.StatusForbidden,
	apperror.CodeConversationInvalid: http.StatusBadRequest,
}

// domainSentinelToCode 兜底：未包装的 domain 哨兵错误映射。
var domainSentinelToCode = map[error]int{
	domain.ErrUserNotFound:        apperror.CodeUserNotFound,
	domain.ErrEmailTaken:          apperror.CodeEmailTaken,
	domain.ErrIdentityTaken:       apperror.CodeIdentityTaken,
	domain.ErrBadCredentials:      apperror.CodeBadCredentials,
	domain.ErrUserBanned:          apperror.CodeUserBanned,
	domain.ErrCaptchaInvalid:      apperror.CodeCaptchaInvalid,
	domain.ErrEmailCodeInvalid:    apperror.CodeEmailCodeInvalid,
	domain.ErrGroupNotFound:       apperror.CodeGroupNotFound,
	domain.ErrNotGroupMember:      apperror.CodeNotGroupMember,
	domain.ErrMsgLimitExceeded:    apperror.CodeNotFriendLimit,
	domain.ErrTokenInvalid:        apperror.CodeTokenInvalid,
	domain.ErrInvalidConversation: apperror.CodeConversationInvalid,
}

// ErrorHandler 统一错误渲染：AppError→映射状态码+信封；未知错误→500+堆栈日志。
func ErrorHandler(logger *slog.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			err := next(c)
			if err == nil {
				return nil
			}
			reqID := c.Response().Header().Get(echo.HeaderXRequestID)
			code, message, status := apperror.CodeInternal, "服务器内部错误", http.StatusInternalServerError

			if ae, ok := apperror.From(err); ok {
				code, message = ae.Code, ae.Message
				if s, ok2 := codeToStatus[ae.Code]; ok2 {
					status = s
				}
				logger.Warn("业务错误", "code", code, "msg", message,
					"request_id", reqID, "path", c.Path(), "err", ae.Err)
			} else if he, ok := err.(*echo.HTTPError); ok { // 框架错误（404/405/401等）
				status = he.Code
				code = status * 100 // 例 404 → 40400，前端仅按非 0 处理
				if m, ok2 := he.Message.(string); ok2 {
					message = m
				} else {
					message = http.StatusText(status)
				}
			} else {
				mapped := false
				for sentinel, c2 := range domainSentinelToCode {
					if errors.Is(err, sentinel) {
						code, message, mapped = c2, sentinel.Error(), true
						if s, ok2 := codeToStatus[c2]; ok2 {
							status = s
						}
						break
					}
				}
				if !mapped {
					logger.Error("未处理错误", "request_id", reqID, "path", c.Path(),
						"err", err, "stack", string(debug.Stack()))
				}
			}
			if c.Response().Committed {
				return nil
			}
			return c.JSON(status, dto.Response{Code: code, Message: message, RequestID: reqID})
		}
	}
}
