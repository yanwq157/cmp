package common

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

//模仿 NewProductionConfig().Build(options…) 相关过程,自己创建，定制化logger对象

var Log *zap.Logger

func InitZap() {
	// 1 日志输出路径
	file, _ := os.OpenFile("./logs/cmp.log", os.O_APPEND|os.O_RDWR, 0744)
	// 把文件对象做成WriteSyncer类型
	writeSyncer := zapcore.AddSync(file)
	// 2 encoder编码，就两种
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoder := zapcore.NewConsoleEncoder(encoderConfig)
	// 3 创建core对象，指定encoder编码，WriteSyncer对象和日志级别
	core := zapcore.NewCore(encoder, writeSyncer, zapcore.DebugLevel)
	// 4 创建logger对象
	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.DebugLevel))
	Log = logger
	//return Log
	//defer logger.Sync()

}

// zap预置的生成logger的方式，都是通过 NewProductionConfig() 来生成相关配置的
// 自定义一个 NewProductionConfig(), Build方法就是通过配置Config对象来生成的logger。

//func InitZap() *zap.Logger {
//	// 1 得到config对象
//	//调用了 NewProductionConfig()方法，内部初始化创建，返回了一个 Config 对象
//	conf := zap.NewProductionConfig()
//	// 2 修改config对象的属性，如编码，输出路径等
//	conf.Encoding = "json"
//	conf.OutputPaths = []string{"./logs/cmp.log"}
//	conf.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
//	//3 通过config对象得到logger对象指针
//	//Build， 内部通过Config对象的配置，利用New方法生成相应的 logger对象，并返回
//	Log, _ := conf.Build()
//	return Log
//}
