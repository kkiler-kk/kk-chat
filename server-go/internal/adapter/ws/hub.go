package ws

import (
	"context"
	"encoding/json"
	"log/slog"
	"sync"
	"time"
)

// Hub 维护 userID → clients 映射（支持多端），广播按快照进行。
type Hub struct {
	mu      sync.RWMutex
	clients map[int64]map[*Client]struct{}
	logger  *slog.Logger
}

func NewHub(logger *slog.Logger) *Hub {
	return &Hub{
		clients: make(map[int64]map[*Client]struct{}),
		logger:  logger,
	}
}

func (h *Hub) register(c *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.clients[c.UserID] == nil {
		h.clients[c.UserID] = make(map[*Client]struct{})
	}
	h.clients[c.UserID][c] = struct{}{}
}

func (h *Hub) unregister(c *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if set, ok := h.clients[c.UserID]; ok {
		delete(set, c)
		if len(set) == 0 {
			delete(h.clients, c.UserID)
		}
	}
}

func (h *Hub) SendToUser(userID int64, env Envelope) {
	if env.Ts == 0 {
		env.Ts = time.Now().Unix()
	}
	b, err := json.Marshal(env)
	if err != nil {
		h.logger.Error("ws 信封序列化失败", "event", env.Event, "err", err)
		return
	}
	h.mu.RLock()
	defer h.mu.RUnlock()
	for c := range h.clients[userID] {
		select {
		case c.send <- b:
		default: // 发送缓冲满：丢弃并记录，避免阻塞业务
			h.logger.Warn("ws 发送缓冲满，丢弃消息", "user_id", userID, "event", env.Event)
		}
	}
}

func (h *Hub) SendToUsers(userIDs []int64, env Envelope) {
	for _, id := range userIDs {
		h.SendToUser(id, env)
	}
}

// Kick 向用户所有连接推送 system.kick 并关闭。
func (h *Hub) Kick(userID int64, reason string) {
	h.SendToUser(userID, Envelope{Event: EventSystemKick, Data: map[string]any{"reason": reason}})
	h.mu.RLock()
	targets := make([]*Client, 0)
	for c := range h.clients[userID] {
		targets = append(targets, c)
	}
	h.mu.RUnlock()
	for _, c := range targets {
		c.closeSend()
	}
}

// Close 通知所有连接关闭（优雅停机时调用）。
func (h *Hub) Close() {
	h.mu.Lock()
	defer h.mu.Unlock()
	for _, set := range h.clients {
		for c := range set {
			c.closeSend()
		}
	}
	h.clients = make(map[int64]map[*Client]struct{})
}

// Run 阻塞至 ctx 取消后 Close（供 main 以 goroutine 运行，统一生命周期）。
func (h *Hub) Run(ctx context.Context) {
	<-ctx.Done()
	h.Close()
}
