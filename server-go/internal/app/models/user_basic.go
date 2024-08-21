package models

import (
	"gorm.io/gorm"
	"time"
)

type UserBasic struct {
	gorm.Model
	ID                   int64
	Name                 string
	Password             string
	Phone                string
	Email                string
	Salt                 string
	Avatar               string
	PersonalitySignature string
	BirthDate            time.Time
	Identity             string
	IsForbidden          int // 是否禁止登录(no) 1 yes 2 no
	IsLogout             int // 是否退出 1 yes 2 no 2 no 说明在线
	IsAdmin              int // 是否是admin 管理员 1 yes 2 no
	ClientIp             string
	ClientPort           uint
	LoginTime            time.Time // 登录时间
	HeartbeatTime        time.Time // 心跳时间
	LoginOutTime         time.Time // 退出时间
	DeviceInfo           string
}

func (u *UserBasic) TableName() string {
	return "user_basic"
}

func (u *UserBasic) BeforeCreate(tx *gorm.DB) (err error) {
	return
}

func (u *UserBasic) AfterCreate(tx *gorm.DB) (err error) {
	return
}
