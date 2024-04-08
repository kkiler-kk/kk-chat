package models

import (
	"gorm.io/gorm"
)

// GroupBasic 群聊分组
type GroupBasic struct {
	gorm.Model
	ID      int64
	Name    string // 群名称
	OwnerId int64  // 群主id
	Icon    string // 头像
	Type    int    // 类型
	Desc    string // 描述
}

func (g GroupBasic) TableName() string {
	return "group_basic"
}
