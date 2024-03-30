package models

import (
	"github.com/jinzhu/gorm"
	"time"
)

type UserBasic struct {
	gorm.Model
	ID            int64
	Name          string
	PassWord      string
	Phone         string
	Email         string
	Salt          string
	Identity      string
	ClientIp      string
	ClientPort    uint
	LoginTime     time.Time
	HeartbeatTime time.Time
	LoginOutTime  time.Time
	IsLogout      bool
	IsAdmin       bool
	DeviceInfo    string
}

func (u *UserBasic) TableName() string {
	return "user_basic"
}
