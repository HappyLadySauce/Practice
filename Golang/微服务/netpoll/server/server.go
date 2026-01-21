package server

import (
	"context"
	"time"

	"github.com/cloudwego/netpoll"
)

// handle 是在连接读取到数据时调用的函数
func handle(ctx context.Context, connection netpoll.Connection) error {
	return nil
}

// prepareHandler 是在连接准备好读取数据时调用的函数
func prepareHandler(connection netpoll.Connection) context.Context {
	ctx := context.Background()
	return ctx
}

func RunServer() {
	listener, err := netpoll.CreateListener("tcp", "127.0.0.1:8083")
	if err != nil {
		panic(err)
	}
	// defer listener.Close()	// netpoll 会检查连接活性关闭 listener

	eventLoop, err := netpoll.NewEventLoop(
		handle,                                // 连接读取到数据时调用的函数
		netpoll.WithOnPrepare(prepareHandler), // 连接准备好读取数据时调用的函数
		// netpoll.WithOnConnect(),	// 连接建立时调用的函数
		netpoll.WithReadTimeout(time.Second), // 读取数据超时时间
	)
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	eventLoop.Serve(listener)
	eventLoop.Shutdown(ctx)
}
