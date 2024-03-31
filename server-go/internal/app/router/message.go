package router

import (
	"github.com/gin-gonic/gin"
	"server-go/internal/app/control"
	"server-go/internal/app/models"
)

func MessageRouter(r *gin.RouterGroup) {
	message := r.Group("message") //, middleware.AuthMiddlewareUpdate
	// 开启websocket manage
	models.Start()
	message.GET("socket", control.MessageControl.Handler) // 用户建立websocket 连接登录
}
