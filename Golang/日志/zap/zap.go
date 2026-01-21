package zap

import (
	"fmt"
	"time"
	"strings"

	"github.com/lestrrat-go/file-rotatelogs"
	// "github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func InitNewLogger(logFile string) *zap.Logger {
	// logger := zap.NewExample()	// 测试阶段
	// logger, _ := zap.NewDevelopment()	// 开发环境
	// logger, _ := zap.NewProduction()	// 生产环境(设置日志最低级别为Info. Error及以上级别日志会打印调用堆栈,会上报位置)

	// 自定义 Zap 日志记录器

	// 打开日志文件

	// file, err := os.OpenFile(logFile, os.O_CREATE | os.O_APPEND | os.O_WRONLY, os.ModePerm)	// 使用追加默认写入日志文件
	// if err != nil {
	// 	panic(err)
	// }

	// 按照日志大小进行分割
	// lumberjackLogger := &lumberjack.Logger{
	// 	Filename: logFile,	
	// 	MaxSize: 10,	// 单位为M, 文件大小超过这么多就会切分
	// 	MaxBackups: 5,	// 保留文件大的最大个数
	// 	MaxAge: 30,		// 保留文件的最大天数
	// 	Compress: false,	// 是否压缩/归档旧文件
	// }

	// 按照时间进行日志分割
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

	// zap 日志核心配置
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05.000")	// 指定时间格式
	encoderConfig.TimeKey = "time"	// 默认时间键为"ts", 改为"time"
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder	// 指定level的显示样式(大写，无颜色)

	// 封装 zap 日志核心
	core := zapcore.NewCore(
		// zapcore.Encoder 日志格式
		// zapcore.NewJSONEncoder(encoderConfig),	// 默认为 JSON 格式
		zapcore.NewConsoleEncoder(encoderConfig),		// 设置为 Console 格式
		// zapcore.AddSync(file),
		// zapcore.AddSync(lumberjackLogger),
		zapcore.AddSync(fout),


		zapcore.InfoLevel,	// 设置最低级别
	)
	logger := zap.New(
		core,
		// 增加调用堆栈, 当日志级别大于等于Error时触发
		zap.AddStacktrace(zapcore.ErrorLevel),
		// 增加函数调用者
		zap.AddCaller(),
		// zap 钩子, 支持添加多个钩子
		zap.Hooks(func(e zapcore.Entry) error {
			if e.Level >= zapcore.ErrorLevel {
				fmt.Println(e.Message)
			}
			return nil
		}),
	)

	// 添加默认参数
	logger = logger.With(
		// 添加默认的Namesapce
		zap.Namespace("test"),
		// 添加默认的Key:Value
		zap.String("ssddffaa", "666"),
	)


	return logger
}


func InitBuildLogger(logFile, level string) *zap.Logger {
	// 按照时间进行日志分割
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

	Config := zap.Config{
		Encoding: "console",
		OutputPaths: []string{"stdout"},	// 输出到标准输出和按时间滚动的日志文件
		ErrorOutputPaths: []string{"stderr"},		// 错误输出到标准错误
		InitialFields: map[string]any{"ssddffaaa": "666"},	// 公共的Field参数
		EncoderConfig: zapcore.EncoderConfig{
			// 指定日志Key名称
			MessageKey: "msg",
			LevelKey: "level",
			// 指定大写Level
			EncodeLevel: zapcore.CapitalLevelEncoder,
			// 指定时间格式
			TimeKey: "time",
			EncodeTime: zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05.000"),
		},
	}
	switch strings.ToLower(level) {
	case "debug":
		Config.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	case "info":
		Config.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	case "warn":
		Config.Level = zap.NewAtomicLevelAt(zap.WarnLevel)
	case "error":
		Config.Level = zap.NewAtomicLevelAt(zap.ErrorLevel)
	case "fatal":
		Config.Level = zap.NewAtomicLevelAt(zap.FatalLevel)
	}
		
	logger, err := Config.Build(
		zap.WrapCore(func(c zapcore.Core) zapcore.Core {
			return zapcore.NewCore(
				zapcore.NewConsoleEncoder(Config.EncoderConfig),
				zapcore.AddSync(fout),
				Config.Level,
			)
		}),
	)
	if err != nil {
		panic(err)
	}

	// 添加默认参数
	logger = logger.With(
		// 添加默认的Namesapce
		zap.Namespace("test"),
		// 添加默认的Key:Value
		zap.String("ssddffaa", "666"),

	)

	return logger
}