package catchaGenerUtil

import (
	"fmt"
	b64 "github.com/mojocn/base64Captcha"
	"server-go/internal/app/models"
)

var store = b64.DefaultMemStore

func CaptchaGenerate() (models.CaptchaData, error) {
	driver := b64.NewDriverDigit(40, 100, 4, 0.2, 80)
	captcha := b64.NewCaptcha(driver, store)
	id, b64s, _, err := captcha.Generate()
	if err != nil {
		fmt.Println(err.Error())
	}
	var result = models.CaptchaData{
		CaptchaId: id,
		Data:      b64s,
	}
	return result, err
}
func CaptchaVerify(data models.CaptchaData) bool {
	return store.Verify(data.CaptchaId, data.Answer, true)
}
