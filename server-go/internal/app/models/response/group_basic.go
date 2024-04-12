package response

// UserGroupListModelRes @Title 返回好友群聊
type UserGroupListModelRes struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Avatar      string `json:"avatar"`
	Identity    string `json:"identity"`
	CountMember int    `json:"countMember"` // 总共多少人数
}
