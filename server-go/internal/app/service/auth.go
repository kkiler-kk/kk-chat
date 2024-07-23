package service

import (
	"github.com/gin-gonic/gin"
	"server-go/common/utility/emailUtil"
	"server-go/common/utility/redisUtil"
	"server-go/internal/consts"
	"time"
)

var AuthService = &authService{}

type authService struct {
}

func (a *authService) SendEmailCode(c *gin.Context, email string) error {
	code := emailUtil.SendEmailCode(email)
	// 存入redis
	err := redisUtil.Set(c, consts.CodeEmail+email, code, (1*time.Hour)/2) // 验证码半小时有效时期
	return err
}
