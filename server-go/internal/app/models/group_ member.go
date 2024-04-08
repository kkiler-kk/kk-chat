package models

import "gorm.io/gorm"

type GroupMember struct {
	gorm.Model
	ID      int64
	GroupID int64
	UserID  int64
	Type    int
	IsAdmin int
}

func (g GroupMember) TableName() string {
	return "group_member"
}
