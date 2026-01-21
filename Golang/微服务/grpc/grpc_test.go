package grpc_test

import (
	"sync"
	"testing"

	"happyladysauce/client"
	"happyladysauce/server"
)

func TestGrpc(t *testing.T) {
	go server.RunServer()

	wg := sync.WaitGroup{}
	wg.Add(3)

	for i := 0; i < 3; i++ {
		go func() {
			defer wg.Done()
			client.RunClient()
		}()
	}
	wg.Wait()
}
