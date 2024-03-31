package request

import "server-go/internal/app/models"

// UserBasicCreateReqInput @Title 添加时候添加
type UserBasicCreateReqInput struct {
	UserName         string `form:"username"`
	OldPassword      string `form:"oldPassword"`
	NewPassword      string `form:"newPassword"`
	Identity         string `form:"identity"`
	Email            string `form:"email"`
	VerificationCode string `form:"verificationCode"`
}

// UserBasicLoginReqInput @Title 登录请求参数
type UserBasicLoginReqInput struct {
	UserName string             `from:"username"`
	Password string             `from:"password"`
	Remember bool               `from:"remember"`
	Code     models.CaptchaData `from:"code"`
}
