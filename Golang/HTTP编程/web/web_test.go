package web_test

import (
	"happyladysauce/web/client"
	"happyladysauce/web/server"
	"testing"
	"time"
)

func TestXxx(t *testing.T) {
	go server.NewWebServer()
	time.Sleep(time.Second)
	client.NewHttpClient()
}