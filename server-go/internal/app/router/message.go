package router

import (
	"github.com/gin-gonic/gin"
	"server-go/internal/app/control"
	"server-go/internal/app/middleware"
	"server-go/internal/websocket"
)

func MessageRouter(r *gin.RouterGroup) {
	message := r.Group("message")                                                          //, middleware.AuthMiddlewareUpdate
	websocket.Start()                                                                      // 开启websocket manage
	message.GET("socket", middleware.AuthMiddlewareUpdate, control.MessageControl.Handler) // 用户建立websocket 连接登录
}
