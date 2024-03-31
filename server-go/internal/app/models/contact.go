package models

import "gorm.io/gorm"

// Contact 联系人
type Contact struct {
	gorm.Model
	OwnerId  int64 // 谁的关系
	TargetId int64 // 对应着谁
	Type     int   // 类型
	Desc     string
}
