package conf

import (
	"fmt"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"log"
	"os"
	"time"
)

const (
	outputDir = "./logs/"
	outPath   = "cmp.log"
	errPath   = "cmp.err"
)

var level zapcore.Level
var logger *zap.Logger

func InitZap() {
	_, err := os.Stat(outputDir)
	if err != nil {
		if os.IsNotExist(err) {
			err := os.Mkdir(outputDir, os.ModePerm)
			if err != nil {
				fmt.Printf("mkdir failed![%v]\n", err)
			}
		}
	}
	// 设置日志格式
	config := zapcore.EncoderConfig{
		MessageKey:  "msg",
		LevelKey:    "level",
		EncodeLevel: zapcore.CapitalLevelEncoder, //将级别转换成大写
		TimeKey:     "ts",
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02 15:04:05"))
		},
		CallerKey:    "file",
		EncodeCaller: zapcore.ShortCallerEncoder,
		EncodeDuration: func(d time.Duration, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendInt64(int64(d) / 1000000)
		},
	}
	//NewConsoleEncoder 创建一个编码器，其输出是为人类而不是机器消耗而设计的。 它以纯文本格式序列化核心日志条目数据（消息、级别、时间戳等），并将结构化上下文保留为 JSON。
	encoder := zapcore.NewConsoleEncoder(config)
	// 设置级别
	switch Config.Zap.Level {
	case "debug":
		level = zap.DebugLevel
	case "info":
		level = zap.InfoLevel
	case "warn":
		level = zap.WarnLevel
	case "error":
		level = zap.ErrorLevel
	case "panic":
		level = zap.PanicLevel
	case "fatal":
		level = zap.FatalLevel
	default:
		level = zap.InfoLevel
	}
	// 实现两个判断日志等级的interface  可以自定义级别展示
	// 在标准错误和标准输出之间拆分日志
	infoLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		//判断warn级别和level级别
		return lvl < zapcore.WarnLevel && lvl >= level
	})

	warnLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.WarnLevel && lvl >= level
	})

	// 获取 info、warn日志文件的io.Writer 抽象 getWriter() 在下方实现
	infoWriter := getWriter(outPath)
	warnWriter := getWriter(errPath)

	// 最后创建具体的Logger
	core := zapcore.NewTee(
		// 将info及以下写入logPath,  warn及以上写入errPath
		zapcore.NewCore(encoder, zapcore.AddSync(infoWriter), infoLevel),
		zapcore.NewCore(encoder, zapcore.AddSync(warnWriter), warnLevel),
		//日志都会在console中展示
		zapcore.NewCore(zapcore.NewConsoleEncoder(config),
			zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout)), level),
	)
	// 需要传入 zap.AddCaller() 才会显示打日志点的文件名和行数
	//AddStacktrace 将 Logger 配置为记录所有消息的堆栈跟踪
	logger = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.WarnLevel))

}

func getWriter(filename string) io.Writer {
	// 生成rotatelogs的Logger 实际生成的文件名 demo.log.YY-mm-dd-HH
	// demo.log是指向最新日志的链接
	hook, err := rotatelogs.New(
		filename+".%Y%m%d%H", // 没有使用go风格反人类的format格式
		rotatelogs.WithLinkName(filename),
		rotatelogs.WithMaxAge(time.Hour*24*30),    // 保存30天
		rotatelogs.WithRotationTime(time.Hour*24), //切割频率 24小时
	)
	if err != nil {
		log.Println("日志启动异常")
		panic(err)
	}
	return hook
}
