package netpoll_test

import (
	"testing"
	"time"

	"happyladysauce/client"
	"happyladysauce/server"
)

func TestNetpoll(t *testing.T) {
	go server.RunServer()
	time.Sleep(1 * time.Second)
	client.RunClient()
}
