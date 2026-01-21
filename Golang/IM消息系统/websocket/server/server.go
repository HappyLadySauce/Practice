package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const (
	DaedLineTime	= 5 * time.Second
	HeartBeatTime	= 3 * time.Second
)

// upgrader 是一个全局的 Upgrader 实例，用于升级 HTTP 连接为 WebSocket 连接。
var (
	upgrader = websocket.Upgrader{
		HandshakeTimeout: 1 * time.Second,
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func handler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Printf("upgrade error: %s\n", err)
		return
	}

	// // 设置 WebSocket 连接的 Pong 处理函数，用于处理客户端发送的 Pong 消息。
	// conn.SetPongHandler(func(appData string) error {
	// 	fmt.Printf("recv pong: %s\n", appData)
	// 	return nil
	// })

	// // 设置 WebSocket 连接的 Ping 处理函数，用于处理客户端发送的 Ping 消息。
	// conn.SetPingHandler(func(appData string) error {
	// 	fmt.Printf("recv ping: %s\n", appData)
	// 	return nil
	// })

	// 设置 WebSocket 连接的 Close 处理函数，用于处理客户端关闭连接的情况。
	conn.SetCloseHandler(func(code int, text string) error {
		fmt.Printf("close: %d, %s\n", code, text)
		conn.Close()
		return nil
	})

	conn.SetReadDeadline(time.Now().Add(DaedLineTime))
	go heartBeat(conn)

	for {
		mtype, message, err := conn.ReadMessage()
		if err != nil {
			fmt.Printf("read error: %v\n", err)
			break
		}
		if mtype == websocket.TextMessage {
			fmt.Printf("read %s\n", string(message))
		} else if mtype == websocket.PongMessage {
			fmt.Printf("receive pong\n")
		}
	}
	conn.WriteMessage(websocket.CloseMessage, nil)
}

func heartBeat(conn *websocket.Conn) {
	// 设置心跳保持 Pong 接收逻辑
	conn.SetPongHandler(func(appData string) error {
		log.Println("receive pong")
		// 设置连接超时时间
		deadline := time.Now().Add(DaedLineTime)
		conn.SetReadDeadline(deadline)
		log.Printf("must read before %s", deadline.Format("2006-01-02 15:04:05"))
		return nil
	})

	// // 向客户端发送 Ping 消息
    // conn.SetWriteDeadline(time.Now().Add(DaedLineTime))
    // if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
    //     log.Printf("write ping error: %v\n", err)
    //     conn.WriteMessage(websocket.CloseMessage, nil)
    // }
    
    ticker := time.NewTicker(HeartBeatTime)
    defer ticker.Stop()
LOOP:
    for {
        select {
        case <- ticker.C:
            conn.SetWriteDeadline(time.Now().Add(DaedLineTime))
            if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
                log.Printf("write ping message error: %v\n", err)
                conn.WriteMessage(websocket.CloseMessage, nil)
                break LOOP
            }
            log.Println("write ping")
        }
    }
}

func main() {
	http.HandleFunc("/", handler)
	if err := http.ListenAndServe("127.0.0.1:8082", nil); err != nil {
		fmt.Printf("start http service error: %s\n", err)
	}
}