package email

import (
	"gopkg.in/gomail.v2"
	"server-go/internal/app/core/config"
)

var email *gomail.Dialer

func InitEmail() {
	email = gomail.NewDialer(config.Instance().Mail.Host, config.Instance().Mail.Port, config.Instance().Mail.UserName, config.Instance().Mail.Password)
}
func Instance() *gomail.Dialer {
	if email == nil {
		InitEmail()
	}
	return email
}
