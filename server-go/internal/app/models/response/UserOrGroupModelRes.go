package response

type UserOrGroupModelRes struct {
	Id       int64  `json:"id"`
	Avatar   string `json:"avatar"`
	Name     string `json:"name"`
	Identity string `json:"identity"`
	IsFriend bool   `json:"isFriend"` // 是否是好友 or 是否加入群聊
}
