package router

import (
	"github.com/gin-gonic/gin"
	"server-go/internal/app/control"
	"server-go/internal/app/middleware"
)

func UserFriendRouter(r *gin.RouterGroup) {
	friendRouter := r.Group("friend")
	friendRouter.POST("add", middleware.AuthMiddlewareUpdate, control.UserFriendControl.AddFriend) // 添加好友
	friendRouter.POST("list", middleware.AuthMiddlewareUpdate, control.UserFriendControl.List)     // 查询当前登录用户好友
}
