package request

// SendEmailCodeReqInput @Title 发送邮箱验证码
type SendEmailCodeReqInput struct {
	Email string `form:"email"`
}
