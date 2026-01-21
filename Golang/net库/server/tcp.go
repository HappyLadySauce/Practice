package server

import (
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
)

func NewTcpServer() error {
	listener, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		return err
	}
	defer listener.Close()

	var wg sync.WaitGroup

	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}

		// 为每个连接启动一个 goroutine 处理，并确保关闭连接
		wg.Add(1)
		go func(c net.Conn) {
			defer wg.Done()
			defer c.Close()

			buf := make([]byte, 4096)
			for {
				n, err := c.Read(buf)
				if errors.Is(err, io.EOF) {
					break
				} else if err != nil {
					fmt.Println("读取数据失败:", err)
					return
				}

				data := string(buf[:n])
				fmt.Println(data)
			}
		}(conn)
	}

	wg.Wait()
	return nil
}
