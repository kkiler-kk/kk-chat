package global

import (
	"fmt"
	"path"
	"runtime"
	"server-go/internal/app/core/config"
	"server-go/internal/app/core/db"
	"server-go/internal/app/core/log"
	"server-go/internal/app/core/mongo"
	"server-go/internal/app/core/redis"
	"server-go/internal/router"
)

func Init() {

	var startStr = "      ,-.      ,-.                     ,---,                   ___     \n  ,--/ /|  ,--/ /|                   ,--.' |                 ,--.'|_   \n,--. :/ |,--. :/ |     ,---,.        |  |  :                 |  | :,'  \n:  : ' / :  : ' /    ,'  .' |        :  :  :                 :  : ' :  \n|  '  /  |  '  /   ,---.'   , ,---.  :  |  |,--.  ,--.--.  .;__,'  /   \n'  |  :  '  |  :   |   |    |/     \\ |  :  '   | /       \\ |  |   |    \n|  |   \\ |  |   \\  :   :  .'/    / ' |  |   /' :.--.  .-. |:__,'| :    \n'  : |. \\'  : |. \\ :   |.' .    ' /  '  :  | | | \\__\\/: . .  '  : |__  \n|  | ' \\ \\  | ' \\ \\`---'   '   ; :__ |  |  ' | : ,\" .--.; |  |  | '.'| \n'  : |--''  : |--'         '   | '.'||  :  :_:,'/  /  ,.  |  ;  :    ; \n;  |,'   ;  |,'            |   :    :|  | ,'   ;  :   .'   \\ |  ,   /  \n'--'     '--'               \\   \\  / `--''     |  ,     .-./  ---`-'   \n                             `----'             `--`---'              "
	fmt.Println(startStr)

	_, filename, _, _ := runtime.Caller(0)
	RootPtah = path.Dir(filename)
	fmt.Println("当前程序运行路径", RootPtah)
	fmt.Println("当前操作系统", SysType)
	// 初始化配置文件
	config.Init()

	// 初始化日志
	log.InitLog()

	// 初始化数据库
	db.InitMySQL()

	// 初始化redis
	redis.InitRedis()

	// 初始化mongo
	mongo.InitMongoDB()
	// 初始化路由
	router.InitRouter()
}
