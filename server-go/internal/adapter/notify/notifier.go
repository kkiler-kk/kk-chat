package notify

import (
	"context"
	"time"

	"server-go/internal/adapter/ws"
	"server-go/internal/usecase/port"
)

// WSNotifier 将 port.Notifier 事件桥接到 ws.Hub。
type WSNotifier struct{ hub *ws.Hub }

func NewWSNotifier(hub *ws.Hub) *WSNotifier { return &WSNotifier{hub: hub} }

func (n *WSNotifier) ToUser(_ context.Context, userID int64, ev port.NotifierEvent) {
	n.hub.SendToUser(userID, ws.Envelope{Event: ev.Event, Data: ev.Data, Ts: time.Now().Unix()})
}

func (n *WSNotifier) ToUsers(_ context.Context, userIDs []int64, ev port.NotifierEvent) {
	n.hub.SendToUsers(userIDs, ws.Envelope{Event: ev.Event, Data: ev.Data, Ts: time.Now().Unix()})
}

var _ port.Notifier = (*WSNotifier)(nil)
