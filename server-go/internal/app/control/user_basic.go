package control

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"server-go/internal/app/models/request"
	"server-go/internal/app/models/response"
	"server-go/internal/app/service"
	"server-go/internal/common/utility"
	"server-go/internal/common/utility/catchaGenerUtil"
	"server-go/internal/common/utility/validate"
	"strconv"
)

var UserControl *userControl

type userControl struct {
}

// CreateUserBasic @Title 创建用户
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

// Login @Title 登录
func (u *userControl) Login(c *gin.Context) {
	var input request.UserBasicLoginReqInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.ErrorResp(c).SetMsg(utility.GetVerErr(err)).WriteJsonExit()
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

// UpdateUserBasic @Title 修改用户信息
func (u *userControl) UpdateUserBasic(c *gin.Context) {

}

// DeleteUserBasic @Title 删除用户信息
func (u *userControl) DeleteUserBasic(c *gin.Context) {

}

// UserDetail @Title 查询用户详细信息
func (u *userControl) UserDetail(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		log.Error().Msg("传入错误id")
		response.ErrorResp(c).SetMsg("传入错误id").WriteJsonExit()
		return
	}
	user, err := service.UserBasicService.UserDetail(c, id)
	if err != nil {
		response.ErrorResp(c).SetMsg(err.Error()).WriteJsonExit()
		return
	}
	response.SuccessResp(c).SetData(user).WriteJsonExit()
}

// Logout @Title 退出登录
func (u *userControl) Logout(c *gin.Context) {
	err := service.UserBasicService.Logout(c)
	if err != nil {
		log.Error().Msgf("用户退出失败%v", err)
		response.ErrorResp(c).SetMsg(err.Error()).WriteJsonExit()
		return
	}
	response.SuccessResp(c).WriteJsonExit()
}

// SearchUser @Title 查询用户
func (u *userControl) SearchUser(c *gin.Context) {
	var input request.UserSearchReqInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.ErrorResp(c).SetMsg(utility.GetVerErr(err)).WriteJsonExit()
		return
	}
	result, err := service.UserBasicService.SearchUser(c, input.Search)
	if err != nil {
		response.ErrorResp(c).SetMsg(err.Error()).WriteJsonExit()
		return
	}
	response.SuccessResp(c).SetData(result).WriteJsonExit()
}
