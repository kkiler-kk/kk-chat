package http

import (
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"log/slog"

	"server-go/internal/adapter/http/dto"
	"server-go/internal/adapter/http/handler"
	"server-go/internal/adapter/http/middleware"
	"server-go/internal/adapter/ws"
	"server-go/internal/config"
	"server-go/internal/platform/token"
	"server-go/internal/usecase/port"
)

type ServerDeps struct {
	Cfg      *config.Config
	Logger   *slog.Logger
	TokenMgr *token.Manager
	Tokens   port.TokenStore
	Auth     port.AuthUseCase
	User     port.UserUseCase
	Friend   port.FriendUseCase
	Group    port.GroupUseCase
	Chat     port.ChatUseCase
	Presence port.PresenceUseCase
	Hub      *ws.Hub
}

// NewServer 装配 Echo 实例与全部路由（spec §3.1）。
func NewServer(d ServerDeps) *echo.Echo {
	e := echo.New()
	e.HideBanner = true
	e.Debug = d.Cfg.Server.Mode == "debug"

	// 中间件链（顺序敏感）：
	// RequestID(最外，先生成 id) → RequestLogger(next 返回后记录最终状态码)
	// → ErrorHandler(捕获内层所有错误并渲染) → Recover(panic→apperror 传回外层) → CORS
	e.Use(echoMiddleware.RequestID())
	e.Use(middleware.RequestLogger(d.Logger))
	e.Use(middleware.ErrorHandler(d.Logger))
	e.Use(middleware.Recover(d.Logger))
	e.Use(echoMiddleware.CORSWithConfig(echoMiddleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "PATCH", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Authorization", "Content-Type"},
	}))

	// 静态资源（上传的图片）
	e.Static("/"+d.Cfg.Server.StaticPath, "./"+d.Cfg.Server.StaticPath)

	authH := handler.NewAuth(d.Auth)
	userH := handler.NewUser(d.User)
	friendH := handler.NewFriend(d.Friend)
	groupH := handler.NewGroup(d.Group)
	chatH := handler.NewChat(d.Chat)
	fileH := handler.NewFile(d.Cfg.Server.StaticPath, d.Logger)
	jwtMW := middleware.JWTAuth(d.TokenMgr, d.Tokens)

	v1 := e.Group("/api/v1")
	// 认证
	v1.POST("/auth/register", authH.Register)
	v1.POST("/auth/login", authH.Login)
	v1.POST("/auth/captcha", authH.Captcha)
	v1.POST("/auth/email-code", authH.EmailCode)
	v1.POST("/auth/logout", authH.Logout, jwtMW)

	// 用户
	v1.GET("/users/me", userH.Me, jwtMW)
	v1.PATCH("/users/me", userH.UpdateMe, jwtMW)
	v1.GET("/users/:id", userH.Detail, middleware.OptionalJWT(d.TokenMgr, d.Tokens))
	v1.GET("/users", userH.Search, jwtMW)

	// 好友
	v1.POST("/friends", friendH.Add, jwtMW)
	v1.GET("/friends", friendH.List, jwtMW)

	// 群组（GET /groups?search=xx 走搜索，见 handler）
	v1.POST("/groups", groupH.Create, jwtMW)
	v1.POST("/groups/:id/members", groupH.Join, jwtMW)
	v1.GET("/groups", groupH.List, jwtMW)

	// 聊天
	v1.POST("/messages", chatH.SendMessage, jwtMW)
	v1.GET("/conversations", chatH.Conversations, jwtMW)
	v1.GET("/conversations/:id/messages", chatH.History, jwtMW)

	// 文件
	v1.POST("/files/images", fileH.UploadImage, jwtMW)

	// WebSocket（升级入口自鉴权）
	v1.GET("/ws", ws.NewHandler(ws.HandlerDeps{
		Hub: d.Hub, TokenMgr: d.TokenMgr, Tokens: d.Tokens,
		Presence: d.Presence, Logger: d.Logger,
	}))

	// 健康检查
	v1.GET("/healthz", func(c echo.Context) error { return dto.OK(c, map[string]string{"status": "up"}) })
	return e
}
