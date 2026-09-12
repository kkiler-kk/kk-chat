package port

import "time"

type RegisterInput struct {
	Identity  string
	Name      string
	Password  string
	Email     string
	EmailCode string
}

type LoginInput struct {
	Account       string // identity 或 email
	Password      string
	CaptchaID     string
	CaptchaAnswer string
}

type UserInfo struct {
	ID        int64      `json:"id"`
	Identity  string     `json:"identity"`
	Name      string     `json:"name"`
	Avatar    string     `json:"avatar"`
	Email     string     `json:"email"`
	Phone     string     `json:"phone"`
	Signature string     `json:"signature"`
	BirthDate *time.Time `json:"birth_date,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
}

type LoginOutput struct {
	Token string    `json:"token"`
	User  *UserInfo `json:"user"`
}

type UpdateUserInput struct {
	Name      *string
	Phone     *string
	Email     *string // 改邮箱必须同时给 EmailCode
	EmailCode string
	Avatar    *string
	Signature *string
	BirthDate *time.Time
}

type UserDetailOutput struct {
	*UserInfo
	IsFriend bool `json:"is_friend"`
	IsSelf   bool `json:"is_self"`
}

type UserSearchItem struct {
	ID       int64  `json:"id"`
	Identity string `json:"identity"`
	Name     string `json:"name"`
	Avatar   string `json:"avatar"`
	IsFriend bool   `json:"is_friend"`
}

type FriendItem struct {
	ID     int64  `json:"id"`
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
	Online bool   `json:"online"`
}

type GroupItem struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Avatar      string `json:"avatar"`
	OwnerID     int64  `json:"owner_id"`
	MemberCount int    `json:"member_count"`
}

type CreateGroupInput struct {
	Name      string
	MemberIDs []int64
}

type SendMessageInput struct {
	ConversationID string // 原始字符串，usecase 内 Parse
	Content        string
	ContentType    string // "text" | "image"
}

type OutgoingMessage struct {
	ID             string    `json:"id"`
	ConversationID string    `json:"conversation_id"`
	SenderID       int64     `json:"sender_id"`
	SenderName     string    `json:"sender_name"`
	SenderAvatar   string    `json:"sender_avatar"`
	Content        string    `json:"content"`
	ContentType    string    `json:"content_type"`
	CreatedAt      time.Time `json:"created_at"`
}

type RecentConversation struct {
	ConversationID string    `json:"conversation_id"`
	Type           string    `json:"type"` // private | group
	PeerID         int64     `json:"peer_id"`
	PeerName       string    `json:"peer_name"`
	PeerAvatar     string    `json:"peer_avatar"`
	Online         bool      `json:"online"`
	LastContent    string    `json:"last_message_content"`
	LastSenderName string    `json:"last_sender_name"`
	LastTime       time.Time `json:"last_time"`
}
