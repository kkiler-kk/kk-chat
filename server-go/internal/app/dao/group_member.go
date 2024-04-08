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

func (u *groupMember) List() {

}
