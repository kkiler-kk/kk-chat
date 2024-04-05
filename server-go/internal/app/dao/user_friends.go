package dao

import (
	"server-go/internal/app/core/db"
	"server-go/internal/app/models"
)

var UserFriends *userFriends

type userFriends struct {
}

func (u *userFriends) Insert(user *models.UserFriends) (ID uint, err error) {
	err = db.Instance().Create(&user).Error
	return user.ID, err
}

func (u *userFriends) Update(user *models.UserFriends) error {
	return db.Instance().Model(&user).Updates(user).Error
}

// FindUserFriends @Title 查询是否为好友关系
func (u *userFriends) FindUserFriends(user *models.UserFriends) (result []*models.UserFriends, err error) {
	query := db.Instance()
	if user.UserId != 0 {
		query = query.Where("user_id = ?", user.UserId)
	}
	if user.FriendId != 0 {
		query = query.Where("friend_id = ?", user.FriendId)
	}
	isError := query.Preload("FriendUser").Find(&result)
	return result, isError.Error
}
