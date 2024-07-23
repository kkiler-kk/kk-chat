package log

import (
	"fmt"
	"github.com/rs/zerolog"
	"os"
	"server-go/common/utility/fileUtils"
	"server-go/internal/app/core/config"
	"server-go/internal/consts"
	"strings"
	"time"
)

type LogHook struct {
}

func (hook LogHook) Run(e *zerolog.Event, level zerolog.Level, msg string) {
	var now = time.Now().Format(consts.TimeFormatDay)
	var logName = strings.Replace(config.Instance().Server.AccessLogPattern, "{Ymd}", now, 1)
	if config.Instance().Server.ErrorLogEnabled && level.String() == "error" {
		logName = strings.Replace(config.Instance().Server.ErrorLogPattern, "{Ymd}", now, 1)
	}
	exist, err := fileUtils.PathExists(config.Instance().Server.LogPath)
	if err != nil {
		fmt.Printf("get dir error![%v]\n", err)
		return
	}
	if !exist { // 文件夹不存在 创建文件夹
		err = os.Mkdir(config.Instance().Server.LogPath, os.ModePerm)
		if err != nil {
			fmt.Printf("mkdir failed![%v]\n", err)
			return
		}
	}
	// 日志存储路径
	var path = config.Instance().Server.LogPath + "/" + logName
	// 写入文件
	fmt.Println("path: ", path)
	fileUtils.PutContentFile(path, fmt.Sprintf("{\"level\":'%s', \"time\": '%s', \"msg\":'%s'}\r\n", level.String(), time.Now().Format(consts.TimeFormat), msg))
}
