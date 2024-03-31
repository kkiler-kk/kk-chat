package dao

import (
	"server-go/internal/app/core/db"
	"server-go/internal/app/models"
)

var UserBasic *userBasic

type userBasic struct {
}

// Insert @Title 添加 insert
func (u *userBasic) Insert(user *models.UserBasic) (ID int64, err error) {
	err = db.Instance().Create(&user).Error
	return user.ID, err
}

// Update @Title 修改
func (u *userBasic) Update(user *models.UserBasic) error {
	return db.Instance().Model(&user).Updates(user).Error
}

// FindUniqueUserBasic @Title 通过唯一键找到User 比如 id int64  identity email  phone
func (u *userBasic) FindUniqueUserBasic(id int64, args []string) (user *models.UserBasic, err error) {
	query := db.Instance().Or("id = ?", id)
	query = query.Or("identity = ?", args[0])
	query = query.Or("email = ?", args[1])
	query = query.Or("phone = ?", args[2])
	result := query.First(&user)
	return user, result.Error
}

// Login @Title 登录
func (u *userBasic) Login(id int64, password string) (user *models.UserBasic, err error) {
	err = db.Instance().Where("id = ?", id).Where("password = ?", password).First(&user).Error
	return
}
