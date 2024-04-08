package dao

import (
	"server-go/internal/app/core/db"
	"server-go/internal/app/models"
)

var GroupBasic *groupBasic

type groupBasic struct {
}

// Insert @Title 添加 insert
func (u *groupBasic) Insert(group *models.GroupBasic) (ID int64, err error) {
	err = db.Instance().Create(&group).Error
	return group.ID, err
}

// Update @Title 修改
func (u *groupBasic) Update(group *models.GroupBasic) error {
	return db.Instance().Model(&group).Updates(group).Error
}

// List @Title 返回群聊消息
func (u *groupBasic) List(group *models.GroupBasic) (list []*models.GroupBasic, err error) {
	query := db.Instance()
	var groupIds []int64
	if group.Name != "" {
		query = query.Where("name  like ?", group.Name+"%")
	}
	if group.OwnerId != 0 {
		var groupMember []*models.GroupMember
		err = db.Instance().Where("user_id = ?", group.OwnerId).Find(&groupMember).Error
		if err != nil {
			return
		}
		for _, member := range groupMember {
			groupIds = append(groupIds, member.GroupID)
		}
		query = query.Where("owner_id = ?", group.OwnerId)
	}
	err = query.Where("id in(?)", groupIds).Find(&list).Error
	return
}
