package request

import (
	"server-go/internal/app/models"
	"time"
)

// UserBasicCreateReqInput @Title 添加时候添加
type UserBasicCreateReqInput struct {
	UserName         string `form:"username" binding:"required"`
	OldPassword      string `form:"oldPassword" binding:"required"`
	NewPassword      string `form:"newPassword" binding:"required"`
	Identity         string `form:"identity" binding:"required"`
	Email            string `form:"email" binding:"required,email"`
	VerificationCode string `form:"verificationCode" binding:"required"`
}

// UserBasicLoginReqInput @Title 登录请求参数
type UserBasicLoginReqInput struct {
	UserName string             `from:"username" binding:"required"`
	Password string             `from:"password" binding:"required"`
	Remember bool               `from:"remember" `
	Code     models.CaptchaData `from:"code" binding:"required"`
}

type UserBasicUpdateReqInput struct {
	ID                   int64     `from:"id" binding:"required"`
	Avatar               string    `from:"avatar" binding:"required"`
	UserName             string    `form:"username" binding:"required"`
	Identity             string    `form:"identity" binding:"required"`
	Email                string    `form:"email" binding:"required,email"`
	VerificationCode     string    `form:"verificationCode" binding:"required"`
	PersonalitySignature string    `form:"personalitySignature"`
	BirthDate            time.Time `from:"birthDate"`
	Phone                string    `form:"phone"`
}

type UserSearchReqInput struct {
	Search string `from:"search" binding:"required" label:"搜索内容"`
}
