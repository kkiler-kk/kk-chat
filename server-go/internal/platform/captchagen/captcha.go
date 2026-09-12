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
	// GenerateIdQuestionAnswer 只生成 (id, 题目文本, 答案)，不渲染图片；
	// 必须再调 DrawCaptcha 绘制并编码为 data URI（与 base64Captcha.Captcha.Generate 流程一致）。
	var q string
	id, q, answer = g.driver.GenerateIdQuestionAnswer()
	var item base64Captcha.Item
	if item, err = g.driver.DrawCaptcha(q); err != nil {
		return
	}
	b64Image = item.EncodeB64string()
	return
}

var _ port.CaptchaGenerator = (*Generator)(nil)
