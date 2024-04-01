package control

import (
	"github.com/gin-gonic/gin"
	"server-go/internal/app/models/request"
	"server-go/internal/app/models/response"
	"server-go/internal/app/service"
	"server-go/internal/common/utility/catchaGenerUtil"
	"server-go/internal/common/utility/validate"
)

var UserControl *userControl

type userControl struct {
}

func (u *userControl) CreateUserBasic(c *gin.Context) {
	var input request.UserBasicCreateReqInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.ErrorResp(c).SetMsg(validate.GetValidateError(err)).WriteJsonExit()
		return
	}
	if input.NewPassword != input.OldPassword {
		response.ErrorResp(c).SetMsg("两次密码不一致, 请重新输入").WriteJsonExit()
		return
	}
	err := service.UserBasicService.CreateUserBasic(c, input)
	if err != nil {
		response.ErrorResp(c).SetMsg(err.Error()).WriteJsonExit()
		return
	}
	response.SuccessResp(c).WriteJsonExit()
}

func (u *userControl) Login(c *gin.Context) {
	var input request.UserBasicLoginReqInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.ErrorResp(c).SetMsg(validate.GetValidateError(err)).WriteJsonExit()
		return
	}
	verifyCode := catchaGenerUtil.CaptchaVerify(input.Code)
	if !verifyCode { // 验证码错误
		response.ErrorResp(c).SetMsg("验证码错误").WriteJsonExit()
		return
	}
	token, err := service.UserBasicService.Login(c, input)
	if err != nil {
		response.ErrorResp(c).SetMsg(err.Error()).WriteJsonExit()
		return
	}
	response.SuccessResp(c).SetData(token).WriteJsonExit()
}
func (u *userControl) UpdateUserBasic(c *gin.Context) {

}

func (u *userControl) DeleteUserBasic(c *gin.Context) {

}
