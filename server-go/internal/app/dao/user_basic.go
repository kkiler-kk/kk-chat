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
	err = db.Instance().Omit("phone").Create(&user).Error
	return user.ID, err
}

// Update @Title 修改
func (u *userBasic) Update(user *models.UserBasic) error {
	return db.Instance().Model(&user).Updates(user).Error
}

// FindUniqueUserBasic @Title 通过唯一键找到User 比如 id int64  identity email  phone
func (u *userBasic) FindUniqueUserBasic(id int64, args []string) (user *models.UserBasic, err error) {
	query := db.Instance()
	if id != 0 {
		query = query.Where("id = ?", id)
	}
	if args[0] != "" { // identity
		query = query.Or("identity = ?", args[0])
	}
	if args[1] != "" { // email
		query = query.Or("email = ?", args[1])
	}
	if args[2] != "" { // phone
		query = query.Or("phone = ?", args[2])
	}
	result := query.First(&user)
	return user, result.Error
}

// FindUniqueUserBasicLikeList @Title 通过唯一键找到User 比如 id int64  identity email  phone 和上一个区别就是 这个是模糊搜索 返回列表
func (u *userBasic) FindUniqueUserBasicLikeList(id int64, args []string) (user []*models.UserBasic, err error) {
	query := db.Instance().Debug()
	if id != 0 {
		query = query.Where("id = ?", id)
	}
	if args[0] != "" { // identity
		query = query.Or("identity like ?", args[0]+"%")
	}
	if args[1] != "" { // email
		query = query.Or("email like ?", args[1]+"%")
	}
	if args[2] != "" { // phone
		query = query.Or("phone like  ?", args[1]+"%")
	}
	result := query.Find(&user)
	return user, result.Error
}

// Login @Title 登录
func (u *userBasic) Login(id int64, password string) (user *models.UserBasic, err error) {
	// is_forbidden = 2  1 == 禁用 账号已被封禁
	err = db.Instance().Debug().Where("id = ?", id).Where("password = ?", password).Where("is_forbidden = 2").First(&user).Error
	return
}
