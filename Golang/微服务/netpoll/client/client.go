package client

import (
	"time"

	"github.com/cloudwego/netpoll"
)

func RunClient() {
	dialer := netpoll.NewDialer()
	conn, err := dialer.DialConnection("tcp", "127.0.0.1:8083", 3*time.Second)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	ReadStringFromServer(conn)
	WriteStringToServer(conn)
}

func ReadStringFromServer(conn netpoll.Connection) {
	var reader = conn.Reader()
	var pkgSize = 1024

	for {
		pkg, _ := reader.Slice(pkgSize)
		go func() {
			pkg.Release()
		}()
	}
}

func WriteStringToServer(conn netpoll.Connection) {
	var write_datas <-chan netpoll.Writer
	var writer = conn.Writer()

	for {
		select {
		case pkg := <-write_datas:
			writer.Append(pkg)
		default:
			if writer.MallocLen() > 0 {
				writer.Flush()
			}
		}
	}
}
