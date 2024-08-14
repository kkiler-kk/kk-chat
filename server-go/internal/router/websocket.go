package router

import (
	"github.com/gin-gonic/gin"
	"server-go/internal/app/middleware"
	"server-go/internal/app/service"
	"server-go/internal/consts"
	"server-go/internal/websocket"
)

// WebSocket ws路由配置
func WebSocket(group *gin.RouterGroup) {
	// 需要鉴权
	wsGroup := group.Group(consts.AppWebSocket)
	// ws
	wsGroup.GET("", middleware.AuthMiddlewareUpdate, websocket.WsPage)
	// 启动websocket鉴听
	websocket.Start()
	// 注册ws消息路由
	websocket.RegisterMsg(websocket.EventHandlers{
		"quit":           service.MessageService.Quit,           // 退出组
		"recentChatList": service.MessageService.RecentChatList, // 返回最近聊天消息列表
		"sendMessage":    service.MessageService.SendMessage,    // 接收到客户端发送过来的消息
		"ping":           service.MessageService.Ping,           // 心跳
		"messageList":    service.MessageService.MessageList,    // 返回消息列表
	})
}
