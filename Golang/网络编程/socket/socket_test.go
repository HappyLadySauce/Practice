package socket_test

import (
	"network/socket/client"
	"network/socket/server"
	"testing"
	"time"
)

func TestTcpSocket(t *testing.T) {
	go server.NewTcpServer()
	time.Sleep(3 * time.Second)
	conn := client.ConnectTcpServer("127.0.0.1:8082")
	client.SendTcpServer(conn)
	conn.Close()
}

func TestUdpSocket(t *testing.T) {
	go server.NewUdpServer()
	time.Sleep(3 * time.Second)
	conn := client.ConnectUdpServer("127.0.0.1:8083")
	client.SendUdpServer(conn)
	conn.Close()
}