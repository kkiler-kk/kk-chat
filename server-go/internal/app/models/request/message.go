package request

import "time"

// SendMessageInputReq websocket 请求客户端 发送给服务的数据
type SendMessageInputReq struct {
	Content    string          `json:"content"`  // 消息内容
	TargetId   int64           `json:"targetId"` // 接收者
	TargetUser UserOrGroupInfo `json:"targetUser"`
	UserId     int64           `json:"userId"` // 发送者
	User       UserOrGroupInfo `json:"user"`
	Type       int             `json:"type"`       // 类型 1 私聊 2 群聊 or 广播 未定
	CreateTime time.Time       `json:"createTime"` // 创建时间
}
type UserOrGroupInfo struct {
	ID     int64  `json:"id"`
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
}

type ListMessageInput struct {
	Type     int   `json:"type"`
	TargetId int64 `json:"targetId"`
	UserId   int64 `json:"userId"`
}
