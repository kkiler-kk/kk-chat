package control

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"server-go/internal/app/core/config"
	"server-go/internal/app/models"
	"server-go/internal/app/models/response"
	"server-go/internal/utility/redisUtil"
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
	token := c.Query("token")
	var user *models.UserBasic
	jsonStr, _ := redisUtil.Get(c, config.Instance().Token.CacheKey+token)
	_ = json.Unmarshal([]byte(jsonStr), &user)
	ws, err := upGrade.Upgrade(c.Writer, c.Request, nil) // 升级成 ws协议
	if err != nil {
		response.ErrorResp(c).SetMsg("系统错误:" + err.Error()).WriteJsonExit()
		return
	}
	currentTime := uint64(time.Now().Unix())
	client := models.NewClient(c, ws, currentTime)
	client.User = user
	models.GetClientManager().Register <- client
	go client.Read()
	go client.Write()
}
