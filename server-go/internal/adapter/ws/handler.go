package ws

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"server-go/internal/platform/token"
	"server-go/internal/usecase/port"
)

type HandlerDeps struct {
	Hub      *Hub
	TokenMgr *token.Manager
	Tokens   port.TokenStore
	Presence port.PresenceUseCase
	Logger   *slog.Logger
}

// NewHandler 返回 echo.HandlerFunc：?token= 鉴权 → 升级 → 注册 → 读写循环。
func NewHandler(d HandlerDeps) echo.HandlerFunc {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true }, // 开发环境放开；生产按配置收紧
	}
	return func(c echo.Context) error {
		raw := c.QueryParam("token")
		if raw == "" {
			return echo.NewHTTPError(http.StatusUnauthorized, "missing token")
		}
		claims, err := d.TokenMgr.Parse(raw)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "invalid token")
		}
		ok, err := d.Tokens.Exists(c.Request().Context(), claims.JTI)
		if err != nil || !ok {
			return echo.NewHTTPError(http.StatusUnauthorized, "token expired")
		}
		conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
		if err != nil {
			d.Logger.Error("ws 升级失败", "err", err)
			return nil
		}
		ctx := c.Request().Context()
		onPing := func(ctx context.Context, uid int64) { _ = d.Presence.Heartbeat(ctx, uid) }
		onClose := func(ctx context.Context, uid int64) { _ = d.Presence.OnDisconnect(ctx, uid) }
		client := d.Hub.newClient(claims.UserID, conn, onPing, onClose)
		d.Hub.register(client)
		if err := d.Presence.OnConnect(ctx, claims.UserID); err != nil {
			d.Logger.Warn("presence 上线失败", "user_id", claims.UserID, "err", err)
		}
		// 读循环生命周期长于请求 ctx（Go 1.21+ WithoutCancel）
		go client.readPump(context.WithoutCancel(ctx))
		go client.writePump()
		return nil
	}
}
