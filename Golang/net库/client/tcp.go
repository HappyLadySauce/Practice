package client

import (
	"fmt"
	"net"
)

func NewTcpClient() error {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		return err
	}
	defer conn.Close()

	for i := range 10 {
		_, err := conn.Write([]byte(fmt.Sprintf("Hello, this is message %d\n", i)))
		if err != nil {
			return err
		}
	}
	
	return nil
}