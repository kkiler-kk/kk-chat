package models

import "github.com/jinzhu/gorm"

type UserBasic struct {
	gorm.Model
	Id            int64
	Name          string
	PassWord      string
	Phone         string
	Email         string
	Identity      string
	ClientIp      string
	ClientPort    uint
	LoginTime     uint64
	HeartbeatTime uint64
	LogOutTime    uint64
	IsLogout      bool
	IsAdmin       bool
	DeviceInfo    string
}

func (u *UserBasic) TableName() string {
	return ""
}
