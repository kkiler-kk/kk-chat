package router

import (
	"github.com/gin-gonic/gin"
	"server-go/internal/app/middleware"
	"server-go/internal/consts"
	"server-go/internal/websocket"
)

// WebSocket ws路由配置
func WebSocket(group *gin.RouterGroup) {
	// 需要鉴权
	group.Group(consts.AppWebSocket, func(c *gin.Context) {
		middleware.AuthMiddlewareUpdate(c)
		// ws
		group.GET("/", websocket.WsPage)
	})
	// 启动websocket鉴听
	websocket.Start()
	// 注册ws消息路由
	websocket.RegisterMsg(websocket.EventHandlers{})
}
