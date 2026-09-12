package usecase

import (
	"server-go/internal/domain"
	"server-go/internal/usecase/port"
)

// ToUserInfo 完整映射（仅用于本人视角/登录响应）。
func ToUserInfo(u *domain.User) *port.UserInfo {
	return &port.UserInfo{
		ID: u.ID, Identity: u.Identity, Name: u.Name, Avatar: u.Avatar,
		Email: u.Email, Phone: u.Phone, Signature: u.Signature,
		BirthDate: u.BirthDate, CreatedAt: u.CreatedAt,
	}
}

// ToPublicUserInfo 脱敏映射：Email/Phone 不输出。
func ToPublicUserInfo(u *domain.User) *port.UserInfo {
	info := ToUserInfo(u)
	info.Email = ""
	info.Phone = ""
	return info
}
