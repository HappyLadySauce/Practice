package logger

import (
    "log/slog"
    "os"
    "path"
    "time"
    "github.com/lestrrat-go/file-rotatelogs"
    "happyladysauce/utils/conf"
)

// 全局单例
var (
	Log *slog.Logger
)

// InitLogger 初始化日志记录器
func InitLogger() {
    logFile := config.Config.GetString("log.file")
    level := config.Config.GetString("log.level")
	// 解析日志级别
	logLevel := slog.LevelDebug
	switch level {
	case "info":
		logLevel = slog.LevelInfo
	case "warn":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	}

    // 确保日志目录存在
    if err := os.MkdirAll(path.Dir(logFile), 0755); err != nil {
        panic(err)
    }

	// 匹配.log则删除.log后缀
	if path.Ext(logFile) == ".log" {
		logFile = logFile[:len(logFile)-4]
	}

    // 轮转文件名模式：<logFile>.YYYYMMDDHH
    fout, err := rotatelogs.New(
        logFile + ".%Y%m%d%H" + ".log",
        rotatelogs.WithMaxAge(7*24*time.Hour),
        rotatelogs.WithRotationTime(1*time.Hour),
        // 注意：RotationCount 与 MaxAge 互斥，保留 MaxAge 即可
    )
    if err != nil {
        panic(err)
    }

	logger := slog.New(
		slog.NewTextHandler(fout, &slog.HandlerOptions{
			AddSource: true,
			Level: logLevel,
            ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
                switch a.Key {
                case slog.TimeKey:
                    // 格式化时间字段
                    t := a.Value.Time()
                    if !t.IsZero() {
                        a.Value = slog.StringValue(t.Format("2006-01-02 15:04:05.000"))
                    } else {
                        a.Value = slog.StringValue(time.Now().Format("2006-01-02 15:04:05.000"))
                    }
                case slog.SourceKey:
                    // 仅保留文件名
                    if src, ok := a.Value.Any().(*slog.Source); ok {
                        src.File = path.Base(src.File)
                        a.Value = slog.AnyValue(src)
                    }
                }
                return a
            },
		}),
	)
	Log = logger
}

