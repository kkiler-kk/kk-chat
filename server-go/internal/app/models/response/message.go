package response

import "time"

type RecentChatListRes struct { // 最近消息聊天列表 (用户)
	ID        int64     `json:"id"`        // id
	Name      string    `json:"name"`      // 姓名
	Content   string    `json:"content"`   // 最近聊天内容
	IsLogout  bool      `json:"isLogout"`  // 是否在线
	Avatar    string    `json:"avatar"`    // 头像
	Type      int       `json:"type"`      // 类型 群聊 or 私聊
	CreatedAt time.Time `json:"createdAt"` // 创建时间
}
