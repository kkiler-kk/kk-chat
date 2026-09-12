package ws

// 事件名常量（与 spec §3.4 严格一致）
const (
	EventPing            = "ping"
	EventPong            = "pong"
	EventChatMessage     = "chat.message"
	EventRecentUpdated   = "chat.recent_updated"
	EventPresenceChanged = "presence.changed"
	EventGroupUpdated    = "group.updated"
	EventSystemKick      = "system.kick"
)

// Envelope 是 WS 统一信封。
type Envelope struct {
	Event string `json:"event"`
	Data  any    `json:"data,omitempty"`
	Ts    int64  `json:"ts"`
}

// Incoming C→S 消息结构（当前仅 ping，预留 data）。
type Incoming struct {
	Event string         `json:"event"`
	Data  map[string]any `json:"data"`
}
