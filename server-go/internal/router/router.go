package router

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"server-go/internal/app/core/config"
	"server-go/internal/app/middleware"
	"server-go/internal/app/router"
)

func InitRouter() {
	r := gin.Default()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(middleware.Cors())
	group := r.Group("api")
	router.InitAppRouter(group)
	r.Run(fmt.Sprintf(":%d", config.Instance().Server.Port))
	return
}
