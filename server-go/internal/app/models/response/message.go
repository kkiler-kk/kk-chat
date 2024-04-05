package response

import "time"

type RecentChatListRes struct { // 最近消息聊天列表 (用户)
	ID        int64     // id
	Name      string    // 姓名
	Content   string    // 最近聊天内容
	Avatar    string    // 头像
	Type      int       // 类型 群聊 or 私聊
	CreatedAt time.Time // 创建时间
}
