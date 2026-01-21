package client

import (
	"fmt"
	"log"
	"net"
	"time"
)

// tcp
func ConnectTcpServer(serverAddr string) *net.TCPConn {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", serverAddr)
	if err != nil {
		fmt.Printf("client err: %v\n", err)
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		fmt.Printf("client err: %v\n", err)
	}
	return conn
}

func SendTcpServer(conn net.Conn) {
	n, err := conn.Write([]byte("Hello"))
	if err != nil {
		fmt.Printf("client err: %v\n", err)
	}
	log.Printf("write byte: %v\n", n)
}

// udp
func ConnectUdpServer(serverAddr string) net.Conn {
	conn, err := net.DialTimeout("udp", serverAddr, 3*time.Minute)
	if err != nil {
		fmt.Printf("connect server failed. err: %v\n", err)
	}
	log.Printf("establish connection to server %s myself %s\n", conn.RemoteAddr().String(), conn.LocalAddr().String())
	return conn
}

func SendUdpServer(conn net.Conn) {
	n, err := conn.Write([]byte("Hello"))
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}
	log.Printf("send bytes: %v\n", n)
}