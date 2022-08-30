package common

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var Config Server

func InitViper() *viper.Viper {
	var config string
	config = "configs/config.yaml"
	v := viper.New()
	v.SetConfigFile(config) // 指定配置文件路径
	err := v.ReadInConfig() // 读取配置文件
	if err != nil {         // 读取配置信息失败
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	// 监控配置文件变化
	v.WatchConfig()

	v.OnConfigChange(func(e fsnotify.Event) { // 配置文件发生变更之后会调用的回调函数
		fmt.Println("config file changed:", e.Name)
		if err := v.Unmarshal(&Config); err != nil { // 配置文件发生变化后要同步到结构体
			fmt.Println(err)
		}
	})
	//反序列化
	if err := v.Unmarshal(&Config); err != nil {
		fmt.Println(err)
	}
	return v
}
