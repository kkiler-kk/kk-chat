package response

import "time"

type UserDetailModelRes struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Phone     string    `json:"phone"`
	Email     string    `json:"email"`
	Avatar    string    `json:"avatar"`
	Identity  string    `json:"identity"`
	CreatedAt time.Time `json:"createdAt"`
}

type LoginModelRes struct {
	Token string              `json:"token"`
	User  *UserDetailModelRes `json:"user"`
}
