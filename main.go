package main

import (
	"cmp/common"
	"cmp/router"
	"github.com/gin-gonic/gin"
)

func main() {

	common.InitViper() // 初始化Viper
	common.InitZap()   // 初始化zap日志库
	common.InitDb()    // gorm连接数据库
	InitServer()

}

func InitServer() {
	r := gin.Default()
	PrivateGroup := r.Group("/api/v1/")
	{
		// 容器相关
		routers.InitContainerRouter(PrivateGroup)
	}
	err := r.Run()
	if err != nil {
		return
	}
}
