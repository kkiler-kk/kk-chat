package models

import (
	"gorm.io/gorm"
	"time"
)

type UserBasic struct {
	gorm.Model
	ID            int64
	Name          string
	Password      string
	Phone         string
	Email         string
	Salt          string
	Identity      string
	IsForbidden   int // 是否禁止登录(no) 1 yes 2 no
	IsLogout      int
	IsAdmin       int
	ClientIp      string
	ClientPort    uint
	LoginTime     time.Time
	HeartbeatTime time.Time
	LoginOutTime  time.Time
	DeviceInfo    string
}

func (u *UserBasic) TableName() string {
	return "user_basic"
}
