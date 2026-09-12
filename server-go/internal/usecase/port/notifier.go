package port

import "context"

// NotifierEvent 传输无关的推送事件；Event 取值见 spec §3.4。
type NotifierEvent struct {
	Event string
	Data  any
}

// Notifier 语义为 fire-and-forget：推送失败只记日志，绝不阻塞业务流程。
type Notifier interface {
	ToUser(ctx context.Context, userID int64, ev NotifierEvent)
	ToUsers(ctx context.Context, userIDs []int64, ev NotifierEvent)
}
