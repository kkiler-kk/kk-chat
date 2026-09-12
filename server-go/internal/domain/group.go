package domain

import "time"

type GroupRole string

const (
	RoleOwner  GroupRole = "owner"
	RoleMember GroupRole = "member"
)

type Group struct {
	ID        int64
	Name      string
	Avatar    string
	OwnerID   int64
	CreatedAt time.Time
}

type GroupMember struct {
	GroupID  int64
	UserID   int64
	Role     GroupRole
	JoinedAt time.Time
}
