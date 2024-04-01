package control

import (
	"github.com/gin-gonic/gin"
	"server-go/internal/app/models/request"
	"server-go/internal/app/models/response"
	"server-go/internal/app/service"
	"server-go/internal/common/utility/catchaGenerUtil"
	"server-go/internal/common/utility/validate"
)

var AuthControl *authControl

type authControl struct {
}

func (a *authControl) SendEmailCode(c *gin.Context) {
	var input request.SendEmailCodeReqInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.ErrorResp(c).SetMsg(validate.GetValidateError(err)).WriteJsonExit()
		return
	}
	err := service.AuthService.SendEmailCode(c, input.Email)
	if err != nil {
		response.ErrorResp(c).SetMsg(validate.GetValidateError(err)).WriteJsonExit()
		return
	}
	response.SuccessResp(c).WriteJsonExit()
}
func (a *authControl) CaptchaGenerate(c *gin.Context) {
	captcha, err := catchaGenerUtil.CaptchaGenerate()
	if err != nil {
		response.ErrorResp(c).SetMsg(validate.GetValidateError(err)).WriteJsonExit()
		return
	}
	response.SuccessResp(c).SetData(captcha).WriteJsonExit()
}
