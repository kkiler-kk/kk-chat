package models

import (
	"gorm.io/gorm"
)

type Message struct {
	gorm.Model
	FormId   int64 // 发送者
	TargetId int64 // 接收者
	Type     int   // 发送类型 群聊 私聊 广播
	Media    int   // 消息类型 文字 图片 音频
	Content  string
	Pic      string
	Url      string
	Desc     string
	Amount   int
}

// Trainer @Title 插入到mongodb 模型
type Trainer struct {
	ID        string `json:"id"`
	Content   string `json:"content"`
	StartTime int64  `json:"startTime"`
	EndTime   int64  `json:"endTime"`
}
