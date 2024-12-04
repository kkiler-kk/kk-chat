package global

import (
	"runtime"
)

var (
	// RootPtah 运行根路径
	RootPtah string
	// SysType 操作系统类型  windows | linux
	SysType = runtime.GOOS
	// Blacklists 黑名单列表
	Blacklists map[string]struct{}
	// 程序启动的打印字符串(banner)
	StartStr string
)
