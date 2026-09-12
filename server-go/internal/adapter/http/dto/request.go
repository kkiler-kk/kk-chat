package dto

type RegisterReq struct {
	Identity  string `json:"identity" validate:"required,min=3,max=32"`
	Name      string `json:"name" validate:"required,min=1,max=64"`
	Password  string `json:"password" validate:"required,min=6,max=64"`
	Email     string `json:"email" validate:"required,email,max=128"`
	EmailCode string `json:"email_code" validate:"required,len=6"`
}

type LoginReq struct {
	Account       string `json:"account" validate:"required"`
	Password      string `json:"password" validate:"required"`
	CaptchaID     string `json:"captcha_id" validate:"required"`
	CaptchaAnswer string `json:"captcha_answer" validate:"required"`
}

type EmailCodeReq struct {
	Email string `json:"email" validate:"required,email"`
}

type UpdateUserReq struct {
	Name      *string `json:"name" validate:"omitempty,min=1,max=64"`
	Phone     *string `json:"phone" validate:"omitempty,max=32"`
	Email     *string `json:"email" validate:"omitempty,email"`
	EmailCode string  `json:"email_code" validate:"omitempty,len=6"`
	Avatar    *string `json:"avatar" validate:"omitempty,max=255"`
	Signature *string `json:"signature" validate:"omitempty,max=255"`
	BirthDate *string `json:"birth_date"` // "2006-01-02"，handler 解析为 time.Time
}

type AddFriendReq struct {
	UserID int64 `json:"user_id" validate:"required"`
}

type CreateGroupReq struct {
	Name      string  `json:"name" validate:"required,min=1,max=64"`
	MemberIDs []int64 `json:"member_ids"`
}

type InviteReq struct {
	UserIDs []int64 `json:"user_ids" validate:"required,min=1"`
}

type SendMessageReq struct {
	ConversationID string `json:"conversation_id" validate:"required"`
	Content        string `json:"content" validate:"required,min=1,max=4096"`
	ContentType    string `json:"content_type" validate:"required,oneof=text image"`
}

type SearchQuery struct {
	Keyword string `query:"search" validate:"omitempty,max=64"`
}

type HistoryQuery struct {
	Cursor string `query:"cursor"` // RFC3339，可空
	Limit  int    `query:"limit" validate:"omitempty,min=1,max=100"`
}
