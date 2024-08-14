package service

import (
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
	"server-go/internal/app/core/db"
	"server-go/internal/app/dao"
	"server-go/internal/app/models"
	"server-go/internal/app/models/request"
	"server-go/internal/app/models/response"
)

var GroupService = &groupService{}

type groupService struct {
}

// CreateGroupMember @Title 创建群聊
func (g *groupService) CreateGroupMember(input request.GroupCreateReqInput) (err error) {
	var group = &models.GroupBasic{
		Name:     input.Name,
		OwnerId:  input.OwnerId,
		Avatar:   input.Avatar,
		Identity: uuid.NewV4().String(), // 默认分配uuid
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

// List @Title 返回群聊列表
func (g *groupService) List(userId int64) (list []*response.UserGroupListModelRes, err error) {
	groupList, err := dao.GroupBasic.List(userId)
	list = make([]*response.UserGroupListModelRes, len(groupList))
	for i, group := range groupList {
		list[i] = &response.UserGroupListModelRes{
			ID:          group.ID,
			Name:        group.Name,
			Avatar:      group.Avatar,
			Identity:    group.Identity,
			CountMember: GroupMemberService.Count(group.ID),
		}
	}
	return
}

// FindGroupName @Title 查看群聊
func (g *groupService) FindGroupName(c *gin.Context, name string) (result []*response.UserOrGroupModelRes, err error) {
	group := &models.GroupBasic{Name: name, Identity: name}
	list, err := dao.GroupBasic.FindByName(group)
	result = make([]*response.UserOrGroupModelRes, len(list))
	for i, temp := range list {
		result[i] = &response.UserOrGroupModelRes{
			Id:       temp.ID,
			Avatar:   temp.Avatar,
			Name:     temp.Name,
			Identity: temp.Identity,

			IsFriend: false,
		}
	}
	return
}

// IsGroupExist @Title 检查是否存在群成员 or 群主中
func (g *groupService) IsGroupExist(groupId, userId int64) bool {
	existGroup := dao.GroupBasic.IsGroupAdmin(groupId, userId)
	existMember := dao.GroupMember.IsGroupMember(groupId, userId)
	return existGroup || existMember
}
