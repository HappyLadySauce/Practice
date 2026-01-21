package slog_test

import (
	"context"
	"testing"
	
	MySlog "happyladysauce"
)

func TestSlog(t *testing.T) {
	logger := MySlog.InitLogger("./log/test.log")
	// logger.Println("this is test log.")	// 原生log库调用
	logger.Info("this is test log.", "a", 3, "b", 7, "sum", 3+7)

	// 使用slog的context功能添加属性
	ctx1 := context.WithValue(context.Background(), "user", "HappyLadySauce")
	ctx2 := context.WithValue(ctx1, "age", 18)
	
	// 创建带有context属性的日志记录
	logger.InfoContext(ctx2, "Welcome", "user", "HappyLadySauce", "age", 18)
}