package server

import (
	"fmt"
	"log"
	"net"
	"time"
)

// TCP 是面向字节流的, 一次 Read 到的数据可能包含了多个报文, 也可能只包含了半个报文, 
// 一条报文在什么地方结束需要通信双方事先约定好

// UDP 是面向报文的, 一次 Read 只读一个报文, 如果没有把一个报文读完, 后面的内容会被丢弃掉,
// 下次就读不到了, 同时取不到 remoteAddr

// 一个报文的大小取决于一次 conn.Write() 调用的大小
// UDP 调用 Read 时, 每调用一次会去取一次 conn.Write() 的内容
// TCP 调用 Read 时, 每次调用一次会读取一段字节流

// tcp server 端,一个连接对应一个 client
func NewTcpServer() {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", "127.0.0.1:8082")
	if err != nil {
		fmt.Printf("open tcp listen failed. err: %v\n", err)
	}
	listener, err := net.ListenTCP("tcp4", tcpAddr)
	if err != nil {
		fmt.Printf("open tcp listen failed. err: %v\n", err)
	}

	log.Println("waiting for client connection ...")
	conn, err := listener.Accept()
	if err != nil {
		fmt.Printf("accept client connection failed. err: %v\n", err)
		return
	}
	log.Printf("establish to client %s\n", conn.RemoteAddr().String())
	// 定义超时,分为读写超时,如果未指定则是通用超时时间
	conn.SetDeadline(time.Now().Add(5 * time.Second))
	defer conn.Close()

	request := make([]byte, 256, 256)
	n, err := conn.Read(request)
	if err != nil {
		fmt.Printf("read client message failed. err: %v\n", err)
	}
	log.Printf("receive %s\n", string(request[:n]))
	listener.Close()
}

// udp server 端,一个连接对应多个 client
func NewUdpServer() {
	udpAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:8083")
	if err != nil {
		fmt.Printf("%v\n", err)
	}
	// 这里与 tcp 不同, net.ListenUDP 直接返回了一个 conn 连接, 而不是一个连接监听器
	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		fmt.Printf("%v\n", err)
	}
	log.Println("waiting for client connection ...")
	defer conn.Close()

	request := make([]byte, 256, 256)
	n, remoteAddr, err := conn.ReadFromUDP(request)
	if err != nil {
		fmt.Printf("%v\n", err)
	}
	log.Printf("receive request %s form %s\n", string(request[:n]), remoteAddr.String())
}
