package middleware

import (
	"context"
	"strings"

	"github.com/labstack/echo/v4"
	"server-go/internal/platform/token"
	"server-go/internal/usecase/apperror"
	"server-go/internal/usecase/port"
)

type ctxKey struct{ name string }

var (
	userIDKey = &ctxKey{"uid"}
	jtiKey    = &ctxKey{"jti"}
)

// authenticate 解析 Bearer token 并校验 Redis 白名单；成功返回 (uid, jti, true)。
func authenticate(ctx context.Context, tokenMgr *token.Manager, tokens port.TokenStore, header string) (int64, string, bool) {
	raw := strings.TrimPrefix(header, "Bearer ")
	if raw == "" || raw == header {
		return 0, "", false
	}
	claims, err := tokenMgr.Parse(raw)
	if err != nil {
		return 0, "", false
	}
	ok, err := tokens.Exists(ctx, claims.JTI)
	if err != nil || !ok {
		return 0, "", false
	}
	return claims.UserID, claims.JTI, true
}

// JWTAuth 强制鉴权：失败返回 apperror 20003，成功将 uid/jti 写入请求 context。
func JWTAuth(tokenMgr *token.Manager, tokens port.TokenStore) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			uid, jti, ok := authenticate(c.Request().Context(), tokenMgr, tokens,
				c.Request().Header.Get("Authorization"))
			if !ok {
				return apperror.New(apperror.CodeTokenInvalid, "登录状态已失效，请重新登录")
			}
			ctx := context.WithValue(c.Request().Context(), userIDKey, uid)
			ctx = context.WithValue(ctx, jtiKey, jti)
			c.SetRequest(c.Request().WithContext(ctx))
			return next(c)
		}
	}
}

// OptionalJWT 宽松鉴权：token 有效则写入 ctx，无效/缺失直接放行（用于游客可访问接口）。
func OptionalJWT(tokenMgr *token.Manager, tokens port.TokenStore) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			uid, jti, ok := authenticate(c.Request().Context(), tokenMgr, tokens,
				c.Request().Header.Get("Authorization"))
			if ok {
				ctx := context.WithValue(c.Request().Context(), userIDKey, uid)
				ctx = context.WithValue(ctx, jtiKey, jti)
				c.SetRequest(c.Request().WithContext(ctx))
			}
			return next(c)
		}
	}
}

func UserIDFromContext(ctx context.Context) (int64, bool) {
	v, ok := ctx.Value(userIDKey).(int64)
	return v, ok
}

func JTIFromContext(ctx context.Context) (string, bool) {
	v, ok := ctx.Value(jtiKey).(string)
	return v, ok
}
