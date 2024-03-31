package global

import (
	"fmt"
	"server-go/internal/app/core/config"
	"server-go/internal/app/core/db"
	"server-go/internal/app/core/log"
	"server-go/internal/app/core/redis"
	"server-go/internal/router"
)

func Init() {

	var startStr = " _   _     _   _                          _                   _   \n( ) ( )   ( ) ( )                        ( )                 ( )_ \n| |/'/'   | |/'/'    ______       ___    | |__        _ _    | ,_)\n| , <     | , <     (______)    /'___)   |  _ `\\    /'_` )   | |  \n| |\\`\\    | |\\`\\               ( (___    | | | |   ( (_| |   | |_ \n(_) (_)   (_) (_)              `\\____)   (_) (_)   `\\__,_)   `\\__)\n                                                                  "
	fmt.Println(startStr)
	// 初始化配置文件
	config.Init()

	// 初始化日志
	log.InitLog()

	// 初始化数据库
	db.InitMySQL()

	// 初始化redis
	redis.InitRedis()
	// 初始化路由
	router.InitRouter()
}
