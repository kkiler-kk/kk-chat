package router

import "github.com/gin-gonic/gin"

func InitAppRouter(group *gin.RouterGroup) {
	groupNew := group.Group("app")
	UserBasicRouter(groupNew)  // 用户路由
	AuthRouter(groupNew)       // 杂七杂八路由
	MessageRouter(groupNew)    // 加载消息路由 socket
	FileRouter(groupNew)       // 加载文件路由
	UserFriendRouter(groupNew) // 加载好友来路由
}
