package dao

import (
	"server-go/internal/app/core/db"
	"server-go/internal/app/models"
)

var GroupMember *groupMember

type groupMember struct {
}

// Insert @Title 添加 insert
func (u *groupMember) Insert(groupMember []*models.GroupMember) (err error) {
	err = db.Instance().Debug().Create(&groupMember).Error
	return err
}

// Update @Title 修改
func (u *groupMember) Update(groupMember *models.GroupMember) error {
	return db.Instance().Model(&groupMember).Updates(groupMember).Error
}

// List @Title 返回群聊成员列表
func (u *groupMember) List(member *models.GroupBasic) (memberList *models.GroupMember, err error) {
	return
}

// IsGroupMember @Title 检查是否是群成员
func (u *groupMember) IsGroupMember(groupID, userID int64) bool {
	query := db.Instance()
	query = query.Where("group_id = ?", groupID)
	query = query.Where("user_id = ?", userID)
	err := query.First(&models.GroupMember{}).Error
	return err == nil // 没报错说明查询到了
}
