package control

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"server-go/internal/app/models/request"
	"server-go/internal/app/models/response"
	"server-go/internal/app/service"
)

var UserFriendControl *userFriendControl

type userFriendControl struct {
}

// AddFriend @Title 添加好友
func (u *userFriendControl) AddFriend(c *gin.Context) {
	var r request.UserFriendAddReqInput
	if err := c.ShouldBind(&r); err != nil {
		log.Error().Msgf("Err: %+v", err)
		response.ErrorResp(c).SetMsg(err.Error()).WriteJsonExit()
		return
	}
	err := service.UserFriendService.AddFriend(r)
	if err != nil {
		log.Error().Msgf("添加好友失败: %+v", err)
		response.ErrorResp(c).SetMsg(err.Error()).WriteJsonExit()
		return
	}
	response.SuccessResp(c).WriteJsonExit()
}

// List @Title 查询当前登录用户的好友
func (u *userFriendControl) List(c *gin.Context) {
	result, err := service.UserFriendService.List(c)
	if err != nil {
		log.Error().Err(err)
		response.ErrorResp(c).SetMsg(err.Error()).WriteJsonExit()
		return
	}
	response.SuccessResp(c).SetData(result).WriteJsonExit()
}
