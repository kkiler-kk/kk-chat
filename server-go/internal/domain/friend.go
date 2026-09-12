package domain

import "time"

type Friendship struct {
	UserID    int64
	FriendID  int64
	CreatedAt time.Time
}
