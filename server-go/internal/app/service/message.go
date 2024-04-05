package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
	"server-go/internal/app/models"
	"server-go/internal/app/models/request"
	"server-go/internal/common/utility/redisUtil"
	"server-go/internal/consts"
	"time"
)

var MessageService = &messageService{}
var RoutersFun map[string]models.EventHandler

type messageService struct {
	WebsocketFunMap map[string]models.EventHandler
}

func init() { // 初始化 websocket 路由
	RoutersFun = make(map[string]models.EventHandler)
	RoutersFun["recentChatList"] = MessageService.RecentChatList // 返回最近聊天消息列表
	RoutersFun["sendMessage"] = MessageService.SendMessage       // 接收到客户端发送过来的消息
	RoutersFun["ping"] = MessageService.Ping                     // 心跳
}

// RecentChatList @Title 返回最近的聊天列表
func (m *messageService) RecentChatList(c *gin.Context, client *models.Client, req *models.WRequest) {
	fmt.Println("recentChatList: Hello ")
	redisChatList := redisUtil.HMGet(c, client.ID, consts.RecentChat)
	client.SendMsg(&models.WResponse{ // 发送消息给客户端
		Event:     "recentChatList",
		Data:      redisChatList,
		Code:      200,
		ErrorMsg:  "",
		Timestamp: time.Now().Unix(),
	})
}

// SendMessage @Title 发送消息
func (m *messageService) SendMessage(c *gin.Context, client *models.Client, req *models.WRequest) {
	var message *request.SendMessageInputReq
	if err := mapstructure.Decode(req.Data, &message); err != nil {
		log.Error().Str("user_id", req.Data["userId"].(string)).Msg("发送消息失败")
		return
	}
	var wResponse = &models.WResponse{ // 封装消息返回给发送消息的客户端
		Event: "sendMessage",
		Data:  message.Content,
		Code:  200,
	}
	client, ok := models.GetClientManager().IdInClient(models.GetUserKey(message.TargetId))
	if ok { // 说明客户端在线
		client.SendMsg(wResponse) // 直接发送消息给客户端
	}
	//redisUtil.HGet(c, consts.RecentChat) // 最近获取聊天
}

func (m *messageService) Ping(c *gin.Context, client *models.Client, req *models.WRequest) {
	fmt.Println("ping: ", req.Event)
	models.SendSuccess(client, req.Event)
}

// SettingRecentChat @Title 设置最近聊天列表
func (m *messageService) SettingRecentChat() {

}
