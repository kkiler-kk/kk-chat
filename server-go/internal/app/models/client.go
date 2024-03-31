package models

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
	uuid "github.com/satori/go.uuid"
	"runtime/debug"
	"server-go/internal/utility/location"
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

// WResponse 输出对象
type WResponse struct {
	Event     string      `json:"event"`              // 事件名称
	Data      interface{} `json:"data,omitempty"`     // 数据
	Code      int         `json:"code"`               // 状态码
	ErrorMsg  string      `json:"errorMsg,omitempty"` // 错误消息
	Timestamp int64       `json:"timestamp"`          // 服务器时间
}

// Client 客户端连接
type Client struct {
	Addr          string          // 客户端地址
	ID            string          // 连接唯一标识
	Socket        *websocket.Conn // 用户连接
	Send          chan *WResponse // 待发送的数据
	SendClose     bool            // 发送是否关闭
	closeSignal   chan struct{}   // 关闭信号
	FirstTime     uint64          // 首次连接时间
	HeartbeatTime uint64          // 用户上次心跳时间
	Tags          []string        // 标签
	User          *UserBasic      // 用户信息
	context       *gin.Context    // Custom context for internal usage purpose.
	IP            string          // 客户端IP
	UserAgent     string          // 用户代理
}

// NewClient 初始化
func NewClient(c *gin.Context, socket *websocket.Conn, firstTime uint64) (client *Client) {
	client = &Client{
		Addr:          socket.RemoteAddr().String(),
		ID:            uuid.NewV4().String(),
		Socket:        socket,
		Send:          make(chan *WResponse, 100),
		SendClose:     false,
		closeSignal:   make(chan struct{}, 1),
		FirstTime:     firstTime,
		HeartbeatTime: firstTime,
		//User:          contexts.GetUser(c),
		IP: location.GetClientIp(c),
		//UserAgent: r.UserAgent(),
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

// 读取客户端数据
func (c *Client) Read() {
	defer func() {
		if r := recover(); r != nil {
			log.Error().Msg("client read err: %+v, stack:%+v, user:%+v")
		}
	}()

	defer func() {
		clientManager.Unregister <- c
		c.close()
	}()

	for {
		_, message, err := c.Socket.ReadMessage()
		if err != nil {
			return
		}
		// 处理消息
		handlerMsg(c, message)
	}
}

// 向客户端写数据
func (c *Client) Write() {
	defer func() {
		if r := recover(); r != nil {
			log.Error().Msg(fmt.Sprintf("client write err: %+v, stack:%+v, user:%+v", r, string(debug.Stack()), c.User))
		}
	}()
	defer func() {
		clientManager.Unregister <- c
		c.Socket.Close()
	}()
	for {
		select {
		case <-c.closeSignal:
			log.Info().Msg(fmt.Sprintf("websocket client quit, user:%+v", c.User))
			return
		case message, ok := <-c.Send:
			if !ok {
				// 发送数据错误 关闭连接
				log.Warn().Msg(fmt.Sprintf("client write message, user:%+v", c.User))
				return
			}
			c.Socket.WriteJSON(message)
		}
	}
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
func (c *Client) close() {
	if c.SendClose {
		return
	}
	c.SendClose = true
	c.closeSignal <- struct{}{}
}

// Heartbeat 心跳更新
func (c *Client) Heartbeat(currentTime uint64) {
	c.HeartbeatTime = currentTime
	return
}

// Close 关闭指定客户端连接
func Close(client *Client) {
	client.close()
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
		Code:      200,
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
	if c.HeartbeatTime+heartbeatExpirationTime <= currentTime {
		timeout = true
	}
	return
}

// WRequest 输入对象
type WRequest struct {
	Event string                 `json:"event"` // 事件名称
	Data  map[string]interface{} `json:"data"`  // 数据
}

// handlerMsg 处理消息
func handlerMsg(client *Client, message []byte) {
	defer func() {
		if r := recover(); r != nil {
			log.Error().Msg(fmt.Sprintf("handlerMsg recover, err:%+v, stack:%+v", r, string(debug.Stack())))
		}
	}()

	var request *WRequest
	if err := json.Unmarshal(message, &request); err != nil {
		log.Error().Msg(fmt.Sprintf("handlerMsg 数据解析失败,err:%+v, message:%+v", err, string(message)))
		return
	}

	if request.Event == "" {
		log.Error().Msg("handlerMsg request.Event is null")
		return
	}
	fmt.Println(fmt.Sprintf("%+v", request))
	//fun, ok := routers[request.Event]
	//if !ok {
	//	g.Log().Warningf(ctxManager, "handlerMsg function id %v: not registered", request.Event)
	//	return
	//}
	//fun(client, request)
}
