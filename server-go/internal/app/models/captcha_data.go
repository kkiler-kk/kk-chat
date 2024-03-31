package models

type CaptchaData struct {
	CaptchaId string `json:"captchaId" from:"captchaId"` //验证码id
	Data      string `json:"data" from:"data"`           //验证码数据base64类型
	Answer    string `json:"answer" from:"answer"`       //验证码数字
}
