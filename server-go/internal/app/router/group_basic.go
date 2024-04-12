package router

import (
	"github.com/gin-gonic/gin"
	"server-go/internal/app/control"
	"server-go/internal/app/middleware"
)

func GroupRouter(r *gin.RouterGroup) {
	groupRouter := r.Group("group")
	groupRouter.Use(middleware.AuthMiddlewareUpdate)
	groupRouter.POST("add", control.GroupControl.Add)                   // 新增聊天群
	groupRouter.POST("join", control.GroupControl.JoinGroup)            // 加入群聊
	groupRouter.POST("/find/:name", control.GroupControl.FindGroupName) // 查询群聊
	groupRouter.POST("list", control.GroupControl.List)                 // 返回群聊
}
