package service

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
	"server-go/internal/app/core/config"
	"server-go/internal/app/models"
	"server-go/internal/app/models/request"
	"server-go/internal/app/models/response"
	"server-go/internal/common/utility/mongoUtils"
	"server-go/internal/common/utility/redisUtil"
	"server-go/internal/consts"
	"sort"
	"strconv"
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
	RoutersFun["messageList"] = MessageService.MessageList       // 返回消息列表
}

// RecentChatList @Title 返回最近的聊天列表
func (m *messageService) RecentChatList(c *gin.Context, client *models.Client, req *models.WRequest) {
	str, _ := redisUtil.HGet(c, consts.RecentChat, client.ID).Result()
	var redisChatList []*response.RecentChatListRes
	_ = json.Unmarshal([]byte(str), &redisChatList)
	for _, chat := range redisChatList { // 查看是否在线
		_, chat.IsLogout = models.GetClientManager().IdInClient(models.GetUserKey(chat.ID)) // 查看用户是否在线
	}
	sort.SliceStable(redisChatList, func(i, j int) bool {
		return redisChatList[i].CreatedAt.Unix() > redisChatList[j].CreatedAt.Unix()
	})
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
	redisCountKey := fmt.Sprintf("%s%s->%s", consts.MsgCount, models.GetUserKey(message.UserId), models.GetUserKey(message.TargetId)) // 比如 1->2 获取1给2发消息的总数
	countStr, _ := redisUtil.Get(c, redisCountKey)
	var wResponse = &models.WResponse{ // 封装消息返回给发送消息的客户端
		Event: "error",
		Code:  500,
	}
	message.CreateTime = time.Now()
	//m.RecentChatList(c, client, nil) // 刷新一下自己聊天记录
	// 检擦是不是好友关系
	if !UserFriendService.IsFriend(message.UserId, message.TargetId, 1) { // TODO 默认为私聊 暂时不做群聊处理
		count, _ := strconv.ParseInt(countStr, 10, 64) // 转化为 int64
		if count >= 3 {                                // 如果不是好友关系并且发送消息超过大于等于三条 直接return
			wResponse.ErrorMsg = "非好友关系 发送信息超过三条 请添加好友再进行聊天"
			client.SendMsg(wResponse)
			return
		}
	}
	wResponse.Event = "sendMessage"
	wResponse.Code = 200
	wResponse.Data = message
	targetClient, ok := models.GetClientManager().IdInClient(models.GetUserKey(message.TargetId))

	if ok { // 说明客户端在线
		targetClient.SendMsg(wResponse) // 直接发送消息给客户端
	}
	client.SendMsg(wResponse)                                                                           // 给客户端发送
	_ = m.settingRecentChat(c, message)                                                                 //设置最近聊天列表
	id := fmt.Sprintf("%s->%s", models.GetUserKey(message.UserId), models.GetUserKey(message.TargetId)) // 1 -> 2 发消息
	b, _ := json.Marshal(message)                                                                       // 将struct 转化为string 存储到mongodb
	err := mongoUtils.Insert(c, config.Instance().Mongo.Database, id, string(b), int64(3*consts.Month)) // 插入到mongodb 默认三个月过期
	if err != nil {
		log.Error().Msgf("消息存储mongodb错误 %s", id)
	}

	_ = redisUtil.Incr(c, redisCountKey) // 聊天总数 + 1
}

// MessageList @Title 返回消息列表
func (m *messageService) MessageList(c *gin.Context, client *models.Client, req *models.WRequest) {
	fmt.Println("messageList", req)
	var searchMessageInput *request.ListMessageInput
	if err := mapstructure.Decode(req.Data, &searchMessageInput); err != nil {
		log.Error().Str("user_id", req.Data["userId"].(string)).Msg("发送消息失败")
		return
	}

	// 获取两个消息
	userMessageList, err := mongoUtils.ListMessAge(c, config.Instance().Mongo.Database, fmt.Sprintf("%s->%s", models.GetUserKey(searchMessageInput.UserId), models.GetUserKey(searchMessageInput.TargetId)))
	if err != nil {
		return
	}
	targetUserMessageList, err := mongoUtils.ListMessAge(c, config.Instance().Mongo.Database, fmt.Sprintf("%s->%s", models.GetUserKey(searchMessageInput.TargetId), models.GetUserKey(searchMessageInput.UserId)))
	if err != nil {
		return
	}
	mergeArr := append(userMessageList, targetUserMessageList...) // 合并数组
	messageList := make([]*request.SendMessageInputReq, len(mergeArr))
	for i, message := range mergeArr { // 把字符串消息转换一下
		var revMsg *request.SendMessageInputReq
		_ = json.Unmarshal([]byte(message.Content), &revMsg)
		messageList[i] = revMsg
	}
	sort.SliceStable(messageList, func(i, j int) bool {
		return messageList[i].CreateTime.Unix() < messageList[j].CreateTime.Unix()
	})
	client.SendMsg(&models.WResponse{ // 发送消息给客户端
		Event:     "messageList",
		Data:      messageList,
		Code:      200,
		ErrorMsg:  "",
		Timestamp: time.Now().Unix(),
	})
}

// settingRecentChat @Title 设置最近聊天列表
func (m *messageService) settingRecentChat(c *gin.Context, message *request.SendMessageInputReq) (err error) {
	var recentModel = &response.RecentChatListRes{ // 发送者先改变一次
		ID:        message.UserId,
		Name:      message.User.Name,
		Content:   message.Content,
		Avatar:    message.User.Avatar,
		Type:      message.Type,
		CreatedAt: time.Now(),
	}
	err = m.setRecentChat(c, message.TargetId, recentModel) // 希望你能理解这段代码 或许也是对未来的我 说的这段话
	recentModel.ID = message.TargetId
	recentModel.Name = message.TargetUser.Name
	recentModel.Avatar = message.TargetUser.Avatar
	err = m.setRecentChat(c, message.UserId, recentModel)
	return err
}

// setRecentChat @Title 设置最近聊天列表 因为有I and target 目标发送人 代码重复 再封装一成 重复利用
func (m *messageService) setRecentChat(c *gin.Context, id int64, userRecentModel *response.RecentChatListRes) (err error) {
	var userRecentList []interface{} // TODO *response.RecentChatListRes 这是原本类型 等我有时间再优化一下把
	str, _ := redisUtil.HGet(c, consts.RecentChat, models.GetUserKey(id)).Result()
	if str != "" { // 说明有最近聊天的人 检查一下 是否存在接收者 信息 更新一下最新消息 怎么理解呢
		err = json.Unmarshal([]byte(str), &userRecentList)
		if err != nil {
			return
		}
		var index = -1
		for i, user := range userRecentList {
			m, ok := user.(map[string]interface{})
			if ok { // 如果在最近聊天列表找到联系人 记录下标 直接break
				values := m["id"]
				if values == float64(userRecentModel.ID) { // 妈的 默认float 64 sb
					index = i
					break
				}
			}
		}
		if index != -1 { //  说明找到了
			userRecentList[index] = userRecentModel // 重置刷新一下
		} else { // 没有就 append
			userRecentList = append(userRecentList, userRecentModel)
		}
	} else { // 说明是聊天第一人咯 直接append
		userRecentList = append(userRecentList, userRecentModel)
	}
	err = redisUtil.HSet(c, consts.RecentChat, models.GetUserKey(id), userRecentList...) // 重置一下最近的聊天列表
	return
}

func (m *messageService) Ping(c *gin.Context, client *models.Client, req *models.WRequest) {
	fmt.Println("ping: ", req.Event)
	UserBasicService.UpdateHeartTime(client.User.ID) // 更新一下心跳时间
	models.SendSuccess(client, req.Event)
}
