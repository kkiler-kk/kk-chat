package captchagen

import (
	"context"

	"github.com/mojocn/base64Captcha"
	"server-go/internal/usecase/port"
)

// Generator 基于 base64Captcha 的 4 位数字图形验证码生成器。
// 答案不落生成器本地存储，由 usecase 存入 port.CaptchaStore（Redis）。
type Generator struct {
	driver base64Captcha.Driver
}

func New() *Generator {
	return &Generator{driver: base64Captcha.NewDriverDigit(80, 240, 4, 0.7, 80)}
}

func (g *Generator) Generate(_ context.Context) (id, answer, b64Image string, err error) {
	// DriverDigit 的 GenerateIdQuestionAnswer 返回 (id, base64图片, 答案)
	id, b64Image, answer = g.driver.GenerateIdQuestionAnswer()
	return
}

var _ port.CaptchaGenerator = (*Generator)(nil)
