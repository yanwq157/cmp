package main

import (
	"cmp/internal/conf"
)

func main() {
	conf.InitViper() // 初始化Viper
	conf.InitZap()   // 初始化zap日志库
	conf.InitDb()    // 初始化连接数据库
}
