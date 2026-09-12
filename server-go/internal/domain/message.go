package domain

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type ConversationType int

const (
	ConvPrivate ConversationType = 1
	ConvGroup   ConversationType = 2
)

func (t ConversationType) String() string {
	if t == ConvGroup {
		return "group"
	}
	return "private"
}

func ParseConversationType(s string) (ConversationType, error) {
	switch s {
	case "private":
		return ConvPrivate, nil
	case "group":
		return ConvGroup, nil
	}
	return 0, ErrInvalidConversation
}

type ContentType string

const (
	ContentText  ContentType = "text"
	ContentImage ContentType = "image"
)

// ConversationID 私聊双方共享同一会话键："u:{小ID}_{大ID}"；群聊："g:{群ID}"。
type ConversationID struct {
	Type  ConversationType
	Value string
}

func NewPrivateConversation(a, b int64) ConversationID {
	lo, hi := a, b
	if a > b {
		lo, hi = b, a
	}
	return ConversationID{Type: ConvPrivate, Value: fmt.Sprintf("u:%d_%d", lo, hi)}
}

func NewGroupConversation(groupID int64) ConversationID {
	return ConversationID{Type: ConvGroup, Value: fmt.Sprintf("g:%d", groupID)}
}

func (c ConversationID) String() string { return c.Value }

func ParseConversationID(s string) (ConversationID, error) {
	switch {
	case strings.HasPrefix(s, "u:"):
		parts := strings.Split(strings.TrimPrefix(s, "u:"), "_")
		if len(parts) != 2 {
			return ConversationID{}, ErrInvalidConversation
		}
		a, err1 := strconv.ParseInt(parts[0], 10, 64)
		b, err2 := strconv.ParseInt(parts[1], 10, 64)
		if err1 != nil || err2 != nil || a == 0 || b == 0 || a > b {
			return ConversationID{}, ErrInvalidConversation
		}
		return ConversationID{Type: ConvPrivate, Value: s}, nil
	case strings.HasPrefix(s, "g:"):
		id, err := strconv.ParseInt(strings.TrimPrefix(s, "g:"), 10, 64)
		if err != nil || id <= 0 {
			return ConversationID{}, ErrInvalidConversation
		}
		return ConversationID{Type: ConvGroup, Value: s}, nil
	}
	return ConversationID{}, ErrInvalidConversation
}

// PrivateParticipants 返回私聊会话的两个参与者ID；群聊会话 ok=false。
func (c ConversationID) PrivateParticipants() (a, b int64, ok bool) {
	if c.Type != ConvPrivate {
		return 0, 0, false
	}
	parts := strings.Split(strings.TrimPrefix(c.Value, "u:"), "_")
	a, _ = strconv.ParseInt(parts[0], 10, 64)
	b, _ = strconv.ParseInt(parts[1], 10, 64)
	return a, b, true
}

// GroupID 返回群聊会话的群ID；私聊会话 ok=false。
func (c ConversationID) GroupID() (id int64, ok bool) {
	if c.Type != ConvGroup {
		return 0, false
	}
	id, err := strconv.ParseInt(strings.TrimPrefix(c.Value, "g:"), 10, 64)
	if err != nil {
		return 0, false
	}
	return id, true
}

type Message struct {
	ID             string // mongo ObjectID hex
	ConversationID ConversationID
	SenderID       int64
	Content        string
	ContentType    ContentType
	CreatedAt      time.Time
}
