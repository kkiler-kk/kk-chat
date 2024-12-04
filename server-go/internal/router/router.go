package router

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"server-go/internal/app/core/config"
	"server-go/internal/app/middleware"
	"server-go/internal/app/router"
)

func InitRouter() {
	r := gin.Default()
	r.Static(config.Instance().Server.StaticPath, "./"+config.Instance().Server.StaticPath)
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(middleware.Cors())
	log.Info().Msg("start gin router~")
	group := r.Group("api")
	router.InitAppRouter(group)
	// 初始化websocket路由
	WebSocket(group)
	err := r.Run(fmt.Sprintf(":%d", config.Instance().Server.Port))
	if err != nil {
		log.Error().Msgf("端口号: %d 启动失败 err: %v", config.Instance().Server.Port, err)
		return
	}
	return
}
