package control

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"server-go/internal/app/models/request"
	"server-go/internal/app/models/response"
	"server-go/internal/app/service"
)

var GroupControl *groupControl

type groupControl struct {
}

// Add @Title 新增聊天群
func (g *groupControl) Add(c *gin.Context) {
	var input request.GroupCreateReqInput
	if err := c.ShouldBind(&input); err != nil {
		log.Error().Err(err)
		response.ErrorResp(c).WriteJsonExit()
		return
	}
	if input.Avatar == "" {
		input.Avatar = "http://" + c.Request.Host + "/resource/public/file/default/member.jpg" // 默认头像
	}
	err := service.GroupService.CreateGroupMember(input)
	if err != nil {
		log.Err(err)
		response.ErrorResp(c).WriteJsonExit()
		return
	}
	response.SuccessResp(c).WriteJsonExit()
}

// List @Title 返回用户群聊
func (g *groupControl) List(c *gin.Context) {
	userId, _ := c.Get("id")
	list, err := service.GroupService.List(userId.(int64))

	if err != nil {
		log.Error().Err(err)
		response.ErrorResp(c).WriteJsonExit()
		return
	}
	response.SuccessResp(c).SetData(list).WriteJsonExit()

}

// FindGroupName @Title  查询 群聊
func (g *groupControl) FindGroupName(c *gin.Context) {
	name := c.Param("name")
	list, err := service.GroupService.FindGroupName(c, name)
	if err != nil {
		response.ErrorResp(c).WriteJsonExit()
		return
	}
	response.SuccessResp(c).SetData(list).WriteJsonExit()
}

// JoinGroup @Title 加入群聊
func (g *groupControl) JoinGroup(c *gin.Context) {

}
