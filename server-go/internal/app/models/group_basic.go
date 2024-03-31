package models

import (
	"gorm.io/gorm"
)

// GroupBasic 群聊分组
type GroupBasic struct {
	gorm.Model
	Name    string // 群名称
	OwnerId int64
	Icon    string // 头像
	Type    int    // 类型
	Desc    string // 描述
}
