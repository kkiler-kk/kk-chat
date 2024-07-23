package emailUtil

import (
	"fmt"
	"gopkg.in/gomail.v2"
	"server-go/common/utility/utils"
	"server-go/internal/app/core/config"
	"server-go/internal/app/core/email"
)

// SendEmailCode @Title 发送邮箱验证码
func SendEmailCode(receiveEmail string) string {
	msg := gomail.NewMessage()
	msg.SetHeader("From", config.Instance().Mail.UserName)
	msg.SetHeader("To", receiveEmail)
	msg.SetHeader("Subject", "kkChat-验证码")

	randomCode := utils.RandStr(6)
	msg.SetBody("text/html", fmt.Sprintf("<b>这是你在kkChat验证码 次验证码三十分钟内有效 请不要告诉别人</b>\"%s\"", randomCode))
	// msg.Attach("/home/User/cat.jpg")
	// Send the email
	if err := email.Instance().DialAndSend(msg); err != nil {
		panic(err)
	}
	return randomCode
}
