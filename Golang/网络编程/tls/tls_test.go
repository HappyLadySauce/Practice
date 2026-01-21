package tls_test

import (
	"testing"
	"time"
	"network/tls/client"
	"network/tls/server"
)

func TestTcpClient(t *testing.T) {
	go server.NewTcpServer("127.0.0.1:8080")
	time.Sleep(100 * time.Millisecond) // 给服务器时间启动
	client.NewTcpClient("127.0.0.1:8080")
}

