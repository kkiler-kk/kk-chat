package log

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
)

func InitLog() {
	// 让日志等级输出为彩色
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: "2006-01-02 15:04:05.000"})
	log.Logger = log.Hook(LogHook{}) // 添加钩子 来日志时触发 添加进文件
}
