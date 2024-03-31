package router

import (
	"github.com/gin-gonic/gin"
	"server-go/internal/app/control"
)

func UserBasicRouter(r *gin.RouterGroup) {
	user := r.Group("user")
	user.POST("add", control.UserControl.CreateUserBasic)    // 新增用户
	user.POST("update", control.UserControl.UpdateUserBasic) // 修改用户信息
	user.POST("delete", control.UserControl.DeleteUserBasic) // 删除用户信息
	user.POST("login", control.UserControl.Login)            // 登录
}
