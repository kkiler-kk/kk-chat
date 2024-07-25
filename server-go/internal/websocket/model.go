package websocket

import (
	"github.com/gin-gonic/gin"
)

// WRequest 输入对象
type WRequest struct {
	Event string                 `json:"event"` // 事件名称
	Data  map[string]interface{} `json:"data"`  // 数据
}

type ClientWResponse struct {
	ID        string
	WResponse *WResponse
}
type TagWResponse struct {
	Tag       string
	WResponse *WResponse
}
type UserWResponse struct {
	UserID    int64
	WResponse *WResponse
}

// WResponse 输出对象
type WResponse struct {
	Event     string      `json:"event"`              // 事件名称
	Data      interface{} `json:"data,omitempty"`     // 数据
	Code      int         `json:"code"`               // 状态码
	ErrorMsg  string      `json:"errorMsg,omitempty"` // 错误消息
	Timestamp int64       `json:"timestamp"`          // 服务器时间
}

// EventHandler 消息处理器
type EventHandler func(c *gin.Context, client *Client, req *WRequest)

type EventHandlers map[string]EventHandler
