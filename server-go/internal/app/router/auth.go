package router

import (
	"github.com/gin-gonic/gin"
	"server-go/internal/app/control"
)

func AuthRouter(r *gin.RouterGroup) {
	auth := r.Group("auth")
	auth.POST("send/email/code", control.AuthControl.SendEmailCode)   // 发送邮箱验证码
	auth.POST("captchaGenerate", control.AuthControl.CaptchaGenerate) // 生成登录验证码
}
