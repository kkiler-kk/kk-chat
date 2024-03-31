package router

import "github.com/gin-gonic/gin"

func InitAppRouter(group *gin.RouterGroup) {
	groupNew := group.Group("app")
	UserBasicRouter(groupNew)
	AuthRouter(groupNew)
	MessageRouter(groupNew)
}
