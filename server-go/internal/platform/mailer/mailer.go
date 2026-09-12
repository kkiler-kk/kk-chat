package mailer

import (
	"context"
	"fmt"

	"gopkg.in/gomail.v2"
	"server-go/internal/usecase/port"
)

type Mailer struct {
	host       string
	port       int
	user, pass string
	from       string
}

func New(host string, port int, user, pass, from string) *Mailer {
	return &Mailer{host: host, port: port, user: user, pass: pass, from: from}
}

func (m *Mailer) Send(_ context.Context, to, subject, body string) error {
	msg := gomail.NewMessage()
	msg.SetHeader("From", m.from)
	msg.SetHeader("To", to)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/html", body)
	d := gomail.NewDialer(m.host, m.port, m.user, m.pass)
	if err := d.DialAndSend(msg); err != nil {
		return fmt.Errorf("发送邮件到 %s 失败: %w", to, err)
	}
	return nil
}

var _ port.EmailSender = (*Mailer)(nil)
