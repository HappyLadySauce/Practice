package zap_test

import (
	Myzap "happyladysauce"
	"testing"

	"go.uber.org/zap"
)

func TestLogger(* testing.T) {
	// logger := Myzap.InitNewLogger("./log/test")
	// logger := Myzap.InitBuildLogger("./log/test")
	logger := Myzap.InitBuildLogger("./log/test", "info")
	logger.Info("this is info log.")

	// zap 一般不使用格式化输出, 而是在使用的时候就指定变量格式, 可以提升大约 50% 的性能
	logger.Info("this is flag test.", zap.Int("age", 18))
	// zap 还支持添加 Namesapce 的方式分割日志
	logger.Error("this is flag test.", zap.Namespace("testNamesapce"), zap.Int("age", 18))

	// zap 也支持使用格式化输出模式, 但是此时性能会降低 50% 左右
	Suger := logger.Sugar()
	value := 10
	Suger.Infof("this is value: %d", value)
}