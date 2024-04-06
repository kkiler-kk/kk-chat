package control

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
	"net/http"
	"runtime/debug"
	"server-go/internal/app/core/config"
	"server-go/internal/app/models"
	"server-go/internal/app/models/response"
	"server-go/internal/app/service"
	_ "server-go/internal/app/service"
	"server-go/internal/common/utility/redisUtil"
	"time"
)

var MessageControl *messageControl

type messageControl struct {
}

var upGrade = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (u *messageControl) Handler(c *gin.Context) {
	token, _ := c.Get("token")
	var user *models.UserBasic
	jsonStr, _ := redisUtil.Get(c, config.Instance().Token.CacheKey+token.(string))
	if jsonStr == "" {
		response.ErrorResp(c).SetMsg("用户状态已过期 请重新登录...")
		return
	}
	_ = json.Unmarshal([]byte(jsonStr), &user)
	ws, err := upGrade.Upgrade(c.Writer, c.Request, nil) // 升级成 ws协议
	if err != nil {
		response.ErrorResp(c).SetMsg("系统错误:" + err.Error()).WriteJsonExit()
		return
	}
	currentTime := uint64(time.Now().Unix())
	client := models.NewClient(c, ws, currentTime)
	client.ID = models.GetUserKey(user.ID) // 设置id
	client.User = user
	models.GetClientManager().Register <- client
	service.UserBasicService.UpdateLoginTime(client.User.ID) // 更新 用户登录了数据库
	go u.read(c, client)
	go u.write(c, client)
}

// read @Title 读取客户端数据
func (u *messageControl) read(c *gin.Context, client *models.Client) {
	// 读取客户端数据
	defer func() {
		if r := recover(); r != nil {
			log.Error().Msg("client read err: %+v, stack:%+v, user:%+v")
		}
	}()

	defer func() {
		models.GetClientManager().Unregister <- client
		client.Close()
	}()

	for {
		_, message, err := client.Socket.ReadMessage()
		if err != nil {
			return
		}
		// 处理消息
		handlerMsg(c, client, message)
	}

}

// write @Title @ 写入数据向客户端
func (u *messageControl) write(c *gin.Context, client *models.Client) {
	// 向客户端写数据
	defer func() {
		if r := recover(); r != nil {
			log.Error().Msg(fmt.Sprintf("client write err: %+v, stack:%+v, user:%+v", r, string(debug.Stack()), client.User))
		}
	}()
	defer func() {
		models.GetClientManager().Unregister <- client
		client.Socket.Close()
	}()
	for {
		select {
		case <-client.CloseSignal:
			service.UserBasicService.UpdateLogout(client.User.ID) // 用户断开websocket链接 退出 更新一下数据库
			log.Info().Msg(fmt.Sprintf("websocket client quit, user:%+v", client.User.Identity))
			return
		case message, ok := <-client.Send:
			if !ok {
				// 发送数据错误 关闭连接
				log.Warn().Msg(fmt.Sprintf("client write message, user:%+v", client.User))
				return
			}
			client.Socket.WriteJSON(message)
		}
	}
}

// handlerMsg 处理消息
func handlerMsg(c *gin.Context, client *models.Client, message []byte) {
	defer func() {
		if r := recover(); r != nil {
			log.Error().Msg(fmt.Sprintf("handlerMsg recover, err:%+v, stack:%+v", r, string(debug.Stack())))
		}
	}()

	var request *models.WRequest
	if err := json.Unmarshal(message, &request); err != nil {
		log.Error().Msg(fmt.Sprintf("handlerMsg 数据解析失败,err:%+v, message:%+v", err, string(message)))
		return
	}

	if request.Event == "" {
		log.Error().Msg("handlerMsg request.Event is null")
		return
	}
	fmt.Println(fmt.Sprintf("%+v", request))
	fun, ok := service.RoutersFun[request.Event]
	if !ok {
		log.Error().Msg(fmt.Sprintf("handlerMsg function id %v: not registered", request.Event))
		return
	}
	fun(c, client, request)
}
