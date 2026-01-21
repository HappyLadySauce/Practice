package logrus_test

import (
	"happyladysauce"
	"testing"
)

func TestLogger(t *testing.T) {
	logger := logrus.InitLogger("./log/test", "info")
	logger.Debug("this is debug log.")
	logger.Info("this is test log.")
	logger.Warn("this is warn log.")
}