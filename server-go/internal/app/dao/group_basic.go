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

// List @Title 返回好友分组
func (u *groupBasic) List(userId int64) (list []*models.GroupBasic, err error) {
	query := db.Instance().Debug()
	var groupIds []int64
	var groupMember []*models.GroupMember
	err = db.Instance().Debug().Where("user_id = ?", userId).Find(&groupMember).Error
	if err != nil {
		return
	}
	for _, member := range groupMember {
		groupIds = append(groupIds, member.GroupID)
	}
	query = query.Where("owner_id = ?", userId)
	err = query.Or("id in(?)", groupIds).Find(&list).Error
	return
}

// FindByName @Title 通过name 查询群聊
func (u *groupBasic) FindByName(group *models.GroupBasic) (list []*models.GroupBasic, err error) {
	query := db.Instance().Debug()
	if group.Name != "" {
		query = query.Or("name like ?", group.Name+"%")
	}
	if group.Identity != "" {
		query = query.Or("identity like ?", group.Identity+"%")
	}
	err = query.Find(&list).Error
	return
}

// IsGroupAdmin @Title  检查是否是群主
func (u *groupBasic) IsGroupAdmin(groupID, userID int64) bool {
	query := db.Instance().Debug()
	query = query.Where("id = ?", groupID)
	query = query.Where("owner_id = ?", userID)
	err := query.First(&models.GroupBasic{}).Error
	return err == nil // 没报错说明查询到了
}
