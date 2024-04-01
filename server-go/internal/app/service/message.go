package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"server-go/internal/app/models"
	"server-go/internal/common/utility/redisUtil"
	"server-go/internal/consts"
	"strconv"
	"time"
)

var MessageService = &messageService{
	WebsocketFunMap: map[string]models.EventHandler{"recentChatList": RecentChatList, "sendMessage": SendMessage},
}

type messageService struct {
	WebsocketFunMap map[string]models.EventHandler
}

// RecentChatList @Title 返回最近的聊天列表
func RecentChatList(c *gin.Context, client *models.Client, req *models.WRequest) {
	fmt.Println("recentChatList: Hello ")
	redisChatList := redisUtil.HMGet(c, strconv.FormatInt(client.User.ID, 10), consts.RecentChat)
	client.SendMsg(&models.WResponse{ // 发送消息给客户端
		Event:     "recentChatList",
		Data:      redisChatList,
		Code:      200,
		ErrorMsg:  "",
		Timestamp: time.Now().Unix(),
	})
}

// SendMessage @Title 发送消息
func SendMessage(c *gin.Context, client *models.Client, req *models.WRequest) {
	fmt.Printf("%+v", req)
}
