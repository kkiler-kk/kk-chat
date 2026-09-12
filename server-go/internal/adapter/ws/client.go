package ws

import (
	"context"
	"encoding/json"
	"log/slog"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 150 * time.Second // > 2 个心跳周期（60s*2）
	pingPeriod     = 60 * time.Second
	maxMessageSize = 8192
)

type Client struct {
	UserID int64
	conn   *websocket.Conn
	send   chan []byte
	hub    *Hub
	logger *slog.Logger
	// onPing 每次收到客户端 ping 时回调（presence 心跳续期）；onClose 连接终止时回调一次。
	onPing      func(ctx context.Context, userID int64)
	onClose     func(ctx context.Context, userID int64)
	closeOnce   sync.Once // send 通道关闭幂等
	disconnOnce sync.Once // onClose 回调幂等
}

func (h *Hub) newClient(userID int64, conn *websocket.Conn,
	onPing, onClose func(context.Context, int64)) *Client {
	return &Client{
		UserID: userID, conn: conn, send: make(chan []byte, 64),
		hub: h, logger: h.logger, onPing: onPing, onClose: onClose,
	}
}

// closeSend 关闭发送通道（writePump 收到通道关闭后发 CloseMessage 退出）。
func (c *Client) closeSend() {
	c.closeOnce.Do(func() { close(c.send) })
}

// fireClose 触发一次 onClose（presence 下线）。
func (c *Client) fireClose(ctx context.Context) {
	c.disconnOnce.Do(func() {
		if c.onClose != nil {
			c.onClose(ctx, c.UserID)
		}
	})
}

// SendToClient 序列化后入发送通道（缓冲满丢弃并告警）。
func (c *Client) SendToClient(env Envelope) {
	b, err := json.Marshal(env)
	if err != nil {
		return
	}
	select {
	case c.send <- b:
	default:
		c.logger.Warn("ws 客户端缓冲满", "user_id", c.UserID)
	}
}

// readPump 读循环：收到应用层 ping → onPing + 回 pong；协议级 pong 刷新读超时。
func (c *Client) readPump(ctx context.Context) {
	defer func() {
		c.hub.unregister(c)
		c.closeSend()
		c.conn.Close()
		c.fireClose(ctx)
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
	for {
		_, data, err := c.conn.ReadMessage()
		if err != nil {
			return
		}
		var in Incoming
		if err := json.Unmarshal(data, &in); err != nil {
			c.logger.Warn("ws 消息解析失败", "user_id", c.UserID, "err", err)
			continue
		}
		if in.Event == EventPing {
			if c.onPing != nil {
				c.onPing(ctx, c.UserID)
			}
			c.SendToClient(Envelope{Event: EventPong, Ts: time.Now().Unix()})
		}
	}
}

// writePump 写循环：从 send 通道取消息写出；pingPeriod 定时发协议级 ping。
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() { ticker.Stop(); c.conn.Close() }()
	for {
		select {
		case msg, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok { // send 已关闭
				c.conn.WriteMessage(websocket.CloseMessage,
					websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
				return
			}
			if err := c.conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
