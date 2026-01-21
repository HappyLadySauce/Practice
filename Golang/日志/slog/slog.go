package slog

import (
	// "log"
	"log/slog"
	"os"
)

func InitLogger(logFile string) *slog.Logger {
	fout, err := os.OpenFile(logFile, os.O_CREATE | os.O_APPEND | os.O_WRONLY, os.ModePerm)
	if err != nil {
		panic(err)
	}

	// 原生 log 日志库
	// logger := log.New(fout, "[MY_LOG]", log.Ldate | log.Lmicroseconds)	// 通过flag参数定义日志的格式, 时间精确到微秒1E-6s

	// slog 扩展库
	logger := slog.New(
		// JSON 格式
		// slog.NewJSONHandler(fout, &slog.HandlerOptions{
		// 	AddSource: true,	// 添加调用者信息
		// 	Level: slog.LevelInfo,	// 设置最低日志级别
		// }),

		slog.NewTextHandler(fout,&slog.HandlerOptions{
			AddSource: true,
			Level: slog.LevelInfo,
			ReplaceAttr: func (groups []string, a slog.Attr) slog.Attr {
				if a.Key != slog.TimeKey {
					return a
				}

				t := a.Value.Time()

				a.Value = slog.StringValue(t.Format("2006-01-02 15:04:05.000"))

				return a
			},
		}),
	)

	return logger
}