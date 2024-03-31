package log

import (
	"fmt"
	"github.com/rs/zerolog"
	"server-go/internal/app/core/config"
	"server-go/internal/consts"
	"server-go/internal/utility/fileUtils"
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
	// 日志存储路径
	var path = config.Instance().Server.LogPath + "/" + logName
	// 写入文件
	fileUtils.PutContentFile(path, fmt.Sprintf("{\"level\":'%s', \"time\": '%s', \"msg\":'%s'}\r\n", level.String(), time.Now().Format(consts.TimeFormat), msg))
}
