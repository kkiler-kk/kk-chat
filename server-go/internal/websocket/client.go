package websocket

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
	"runtime/debug"
	"server-go/common/utility/location"
	"server-go/internal/app/models"
	"server-go/internal/consts"
	"time"
)

const (
	// 用户连接超时时间
	heartbeatExpirationTime = 5 * 60
)

// 用户登录
type login struct {
	UserId int64
	Client *Client
}

// Client 客户端连接
type Client struct {
	Addr          string            // 客户端地址
	ID            string            // 连接唯一标识
	Socket        *websocket.Conn   // 用户连接
	Send          chan *WResponse   // 待发送的数据
	SendClose     bool              // 发送是否关闭
	CloseSignal   chan struct{}     // 关闭信号
	FirstTime     uint64            // 首次连接时间
	HeartbeatTime uint64            // 用户上次心跳时间
	Tags          []int64           // 标签
	GroupIds      []int64           // 群聊id若有次id 就推送消息
	User          *models.UserBasic // 用户信息
	context       *gin.Context      // Custom context for internal usage purpose.
	IP            string            // 客户端IP
	UserAgent     string            // 用户代理
}

// NewClient 初始化
func NewClient(c *gin.Context, socket *websocket.Conn, firstTime uint64, user *models.UserBasic) (client *Client) {
	client = &Client{
		ID:            GetUserKey(user.ID),
		Addr:          socket.RemoteAddr().String(),
		Socket:        socket,
		Send:          make(chan *WResponse, 100),
		SendClose:     false,
		CloseSignal:   make(chan struct{}, 1),
		FirstTime:     firstTime,
		HeartbeatTime: firstTime,
		User:          user,
		IP:            location.GetClientIp(c),
		//UserAgent:     c.UserAgent(),
	}
	return
}

// GetKey 读取客户端数据
func (l *login) GetKey() (key string) {
	key = GetUserKey(l.UserId)
	return
}

// GetUserKey 获取用户key
func GetUserKey(userId int64) (key string) {
	key = fmt.Sprintf("%s_%d", "ws", userId)
	return
}

// SendMsg 发送数据
func (c *Client) SendMsg(msg *WResponse) {
	if c == nil || c.SendClose {
		return
	}
	defer func() {
		if r := recover(); r != nil {
			log.Error().Msg(fmt.Sprintf("SendMsg err:%+v, stack:%+v", r, string(debug.Stack())))
		}
	}()
	c.Send <- msg
}

// 关闭客户端
func (c *Client) Close() {
	if c.SendClose {
		return
	}
	c.SendClose = true
	c.CloseSignal <- struct{}{}
}

// Heartbeat 心跳更新
func (c *Client) Heartbeat(currentTime uint64) {
	c.HeartbeatTime = currentTime
	return
}

// Close 关闭指定客户端连接
func Close(client *Client) {
	client.Close()
}

// SendSuccess 发送成功消息
func SendSuccess(client *Client, event string, data ...interface{}) {
	d := interface{}(nil)
	if len(data) > 0 {
		d = data[0]
	}
	client.SendMsg(&WResponse{
		Event:     event,
		Data:      d,
		Code:      consts.StatusOK,
		Timestamp: time.Now().Unix(),
	})
	before(client)
}

// before
func before(client *Client) {
	client.Heartbeat(uint64(time.Now().Unix()))
}

// IsHeartbeatTimeout 心跳是否超时
func (c *Client) IsHeartbeatTimeout(currentTime uint64) (timeout bool) {
	if c.HeartbeatTime+consts.HeartbeatExpirationTime <= currentTime {
		timeout = true
	}
	return
}

// 关闭客户端
func (c *Client) close() {
	if c.SendClose {
		return
	}
	c.SendClose = true
	c.CloseSignal <- struct{}{}
}

// 读取客户端数据
func (c *Client) read(ctx *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			log.Error().Msg(fmt.Sprintf("client read err: %+v, stack:%+v, user:%+v", r, string(debug.Stack()), c.User))
		}
	}()

	defer c.close()

	for {
		_, message, err := c.Socket.ReadMessage()
		if err != nil {
			return
		}
		// 处理消息
		handlerMsg(ctx, c, message)
	}
}

// 向客户端写数据
func (c *Client) write(ctx *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			log.Warn().Msg(fmt.Sprintf("client write err: %+v, stack:%+v, user:%+v", r, string(debug.Stack()), c.User))
		}
	}()
	defer func() {
		clientManager.Unregister <- c
		_ = c.Socket.Close()
	}()
	for {
		select {
		case <-c.CloseSignal:
			log.Info().Msg(fmt.Sprintf("websocket client quit, user:%+v", c.User))
			return
		case message, ok := <-c.Send:
			if !ok {
				// 发送数据错误 关闭连接
				log.Warn().Msg(fmt.Sprintf("client write message, user:%+v", c.User))
				return
			}
			_ = c.Socket.WriteJSON(message)
		}
	}
}
