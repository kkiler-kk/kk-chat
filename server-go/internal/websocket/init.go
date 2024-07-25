package websocket

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
	"net/http"
	"server-go/common/utility/redisUtil"
	"server-go/internal/app/core/config"
	"server-go/internal/app/models"
	"server-go/internal/app/models/response"
	"time"
)

var (
	clientManager = NewClientManager()            // 客户端管理
	routers       = make(map[string]EventHandler) // 消息路由
)

// Start @Title 启动
func Start() {
	go clientManager.start()
	go clientManager.ping()
	log.Debug().Msg("start websocket...")
}

// Stop 关闭
func Stop() {
	clientManager.closeSignal <- struct{}{}
}

// WsPage ws入口
func WsPage(c *gin.Context) {
	token, _ := c.Get("token")
	var user *models.UserBasic
	jsonStr, _ := redisUtil.Get(c, config.Instance().Token.CacheKey+token.(string))
	if jsonStr == "" {
		response.ErrorResp(c).SetMsg("用户状态已过期，请重新登陆")
		return
	}
	_ = json.Unmarshal([]byte(jsonStr), &user)
	upGrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	conn, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {

		return
	}
	currentTime := uint64(time.Now().Unix())
	client := NewClient(c, conn, currentTime, user)
	clientManager.Register <- client // 新用户上线
	go client.read(c)
	go client.write(c)
	// 用户连接事件
	clientManager.Register <- client
}
