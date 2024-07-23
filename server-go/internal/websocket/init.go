package websocket

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
	"net/http"
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
	client := NewClient(c, conn, currentTime)
	go client.read(c)
	go client.write(c)
	// 用户连接事件
	clientManager.Register <- client
}
