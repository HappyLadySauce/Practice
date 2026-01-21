package main

import (
    "fmt"
    "time"

    "github.com/gorilla/websocket"
)

const DeadlineTime = 5 * time.Second

func request() {
    conn, resp, err := websocket.DefaultDialer.Dial("ws://127.0.0.1:8082/", nil)
    if err != nil {
        panic(err)
    }
    conn.SetReadDeadline(time.Now().Add(DeadlineTime))
    defer conn.Close()
    defaultPing := conn.PingHandler()
    conn.SetPingHandler(func(appData string) error {
        conn.SetReadDeadline(time.Now().Add(DeadlineTime))
        return defaultPing(appData)
    })

    conn.SetPongHandler(func(appData string) error {
        fmt.Printf("Pong: %s\n", appData)
        conn.SetReadDeadline(time.Now().Add(DeadlineTime))
        return nil
    })

	conn.SetCloseHandler(func(code int, text string) error {
		fmt.Printf("close: %d, %s\n", code, text)
		conn.Close()
		return nil
	})

	go func() {
		for {
			mtype, message, err := conn.ReadMessage()
			if err != nil {
				fmt.Printf("read error: %v\n", err)
				return
			}
			if mtype == websocket.TextMessage {
				fmt.Printf("read %s\n", string(message))
			}
		}
	}()

	fmt.Println("http status:", resp.Status)
	fmt.Println("http headers:")
	for key, values := range resp.Header {
		fmt.Printf("%s: %s\n", key, values[0])
	}

	for {
		err := conn.WriteMessage(websocket.TextMessage, []byte("你好"))
		if err != nil {
			fmt.Printf("write error: %v\n", err)
			break
		}
		time.Sleep(2 * time.Second)
	}
}

func main() {
	request()
}
