package models

import "gorm.io/gorm"

// UserFriend @Title 关系表
type UserFriends struct {
	gorm.Model
	UserId     int64     `json:"userId"`
	User       UserBasic `json:"user"`
	FriendId   int64     `json:"friendId"`
	FriendUser UserBasic `json:"friendUser"  gorm:"foreignKey:friend_id;references:id"`
	Remark     string    `json:"remark"`
}
