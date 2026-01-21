package logrus

import (
	"fmt"
	"strings"
	"time"

	"github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
)

// 定义一个钩子结构体, 用于实现 logrus.Hook 接口
type AppHook struct {
	AppName string
}

// 适用于哪些 Level
func (h *AppHook) Levels() []logrus.Level {
	return []logrus.Level {
		logrus.ErrorLevel,
		logrus.FatalLevel,
		logrus.PanicLevel,
	}
	// return logrus.AllLevels	// 所有日志均生效
}

// 在 Fire 函数时可读取或修改 logrus.Entry
func (h *AppHook) Fire(entry *logrus.Entry) error {
	entry.Data["app"] = h.AppName	// 修改 logrus.Entry: 添加了一个 Key:value
	fmt.Println(entry.Message)	// 读取 logrus.Entry, 比如将 Error, Fatal 和 Panic 级别的错误日志发送到 logstash, kafka 等
	return nil
}

func InitLogger(logFile, level string) *logrus.Logger {
	// 新建一个 logrus 日志记录器 logger
	logger := logrus.New()

	// 设置日志级别
	switch strings.ToLower(level) {
	case "debug":
		logger.SetLevel(logrus.DebugLevel)
	case "info":
		logger.SetLevel(logrus.InfoLevel)
	case "warn":
		logger.SetLevel(logrus.WarnLevel)
	case "error":
		logger.SetLevel(logrus.ErrorLevel)
	case "panic":
		logger.SetLevel(logrus.PanicLevel)
	}

	// 设置日志文本格式
	logger.SetFormatter(&logrus.TextFormatter{
		// ForceColors: true,	// 强制显示颜色(仅在某些终端中正常工作)
		DisableColors: true,	// 强制不显示颜色
		TimestampFormat: "2006-01-02 15:04:05.000",	// 显示ms
	})

	// 设置为JSON日志格式
	// logger.SetFormatter(&logrus.JSONFormatter{
	// 	TimestampFormat: "2025-01-02 15:02:02.000",
	// })

	fout, err := rotatelogs.New(
		logFile+".%Y%m%d%H"+".log",	// 指定日志文件的路径和名称,路径不存在时会创建
		// rotatelogs.WithLinkName(logFile),	// 创建最新的日志软连接
		rotatelogs.WithRotationTime(1 * time.Hour),		// 每隔1个小时生成一份新的日志文件
		rotatelogs.WithMaxAge(7 * 24 * time.Hour),		// 只保留最近7天的日志
		// rotatelogs.WithRotationCount()		// 只保留最近的几份日志
	)
	if err != nil {
		panic(err)
	}
	logger.SetOutput(fout)	// 设置输出日志文件
	// logger.SetOutput(os.Stdout)	// 将日志输出到终端
	logger.SetReportCaller(true)	// 输出是从哪里开始调起的日志打印, 日志中会包含file和func

	// 添加钩子函数
	logger.AddHook(&AppHook{"666"})

	return logger
}
