package websocket

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"runtime/debug"
)

// handlerMsg 处理消息
func handlerMsg(c *gin.Context, client *Client, message []byte) {
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
	fun, ok := routers[request.Event]
	if !ok {
		log.Error().Msg(fmt.Sprintf("handlerMsg function id %v: not registered", request.Event))
		return
	}
	fun(c, client, request)
}

// RegisterMsg 注册消息
func RegisterMsg(handlers EventHandlers) {
	for id, f := range handlers {
		if _, ok := routers[id]; ok {
			log.Error().Msg(fmt.Sprintf("RegisterMsg function id %v: already registered", id))
			return
		}
		routers[id] = f
	}
}
