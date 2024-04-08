package response

import "time"

type UserFriendListModelRes struct {
	ID           int64     `json:"id"`
	Name         string    `json:"name"`
	Avatar       string    `json:"avatar"`
	Identity     string    `json:"identity"`
	IsLogout     bool      `json:"isLogout"` // 是否退出 1 yes 2 no 2 no 说明在线
	LoginOutTime time.Time `json:"loginOutTime"`
}
