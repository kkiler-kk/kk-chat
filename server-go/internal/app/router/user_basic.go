package router

import (
	"github.com/gin-gonic/gin"
	"server-go/internal/app/control"
	"server-go/internal/app/middleware"
)

func UserBasicRouter(r *gin.RouterGroup) {
	user := r.Group("user")
	user.POST("add", control.UserControl.CreateUserBasic)                                     // 新增用户
	user.POST("update", middleware.AuthMiddlewareUpdate, control.UserControl.UpdateUserBasic) // 修改用户信息
	user.POST("delete", control.UserControl.DeleteUserBasic)                                  // 删除用户信息
	user.POST("login", control.UserControl.Login)                                             // 登录
	user.POST("logout", middleware.AuthMiddlewareUpdate, control.UserControl.Logout)          // 退出登录
	user.POST("detail/:id", middleware.AuthMiddlewareUpdate, control.UserControl.UserDetail)  // 查看用户信息
	user.POST("search", middleware.AuthMiddlewareUpdate, control.UserControl.SearchUser)      // 通过 id or email or phone 查询用户
}
