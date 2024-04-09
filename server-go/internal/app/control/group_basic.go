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
	if input.Icon == "" {
		input.Icon = "http:" + c.Request.Host + "/resource/public/file/default/member.png" // 默认头像
	}
	err := service.GroupService.CreateGroupMember(input)
	if err != nil {
		log.Err(err)
		response.ErrorResp(c).WriteJsonExit()
		return
	}
	response.SuccessResp(c).WriteJsonExit()
}

// FindGroupName @Title
func (g *groupControl) FindGroupName(c *gin.Context) {
	//name := c.Param("name")

}
