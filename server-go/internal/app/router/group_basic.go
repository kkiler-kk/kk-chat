package router

import (
	"github.com/gin-gonic/gin"
	"server-go/internal/app/control"
	"server-go/internal/app/middleware"
)

func GroupRouter(r *gin.RouterGroup) {
	groupRouter := r.Group("group")
	groupRouter.Use(middleware.AuthMiddlewareUpdate)
	groupRouter.POST("/add", control.GroupControl.Add) // 新增聊天群
}
