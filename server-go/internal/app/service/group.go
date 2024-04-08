package service

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"server-go/internal/app/core/db"
	"server-go/internal/app/dao"
	"server-go/internal/app/models"
	"server-go/internal/app/models/request"
)

var GroupService = &groupService{}

type groupService struct {
}

// CreateGroupMember @Title 创建群聊
func (g *groupService) CreateGroupMember(input request.GroupCreateReqInput) (err error) {
	var group = &models.GroupBasic{
		Name:    input.Name,
		OwnerId: input.OwnerId,
		Icon:    input.Icon,
	}
	return db.Instance().Transaction(func(tx *gorm.DB) error {
		id, err := dao.GroupBasic.Insert(group)
		if err != nil {
			return err
		}
		var memberList = make([]*models.GroupMember, len(input.MembersID))
		for i, memberId := range input.MembersID {
			memberList[i] = &models.GroupMember{
				GroupID: id,
				UserID:  memberId,
			}
		}
		err = dao.GroupMember.Insert(memberList)
		return err
	})
}

// CreateGroupMember @Title 创建群聊
func (g *groupService) FindGroupName(c *gin.Context, name string) {

}
