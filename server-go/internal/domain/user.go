package domain

import "time"

type UserStatus int

const (
	UserStatusActive UserStatus = 1
	UserStatusBanned UserStatus = 2
)

type User struct {
	ID        int64
	Identity  string
	Name      string
	Password  string // bcrypt hash，任何出参不得携带
	Email     string
	Phone     string
	Avatar    string
	Signature string
	BirthDate *time.Time
	IsAdmin   bool
	Status    UserStatus
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (u *User) Banned() bool { return u.Status == UserStatusBanned }
