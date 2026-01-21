package main

import (
	"context"
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

var ExchangeName string = "HappyLadySauce"

func main() {
	// 连接到 RabbitMQ Server
	conn, err := amqp.Dial("amqp://test:test@100.100.100.5:5672")
	if err != nil {
		log.Fatalf("无法连接到 RabbitMQ: %v", err)
	}
	defer conn.Close()

	// 创建 channel
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("无法创建通道: %v", err)
	}
	defer ch.Close()

	// 声明队列
	// _, err = ch.QueueDeclare(
	// 	"hello",		// 队列名称
	// 	true,		// 持久化
	// 	true,		// 自动删除
	// 	false,
	// 	false,
	// 	nil,
	// )
	// if err != nil {
	// 	log.Fatalf("无法声明队列: %v", err)
	// }

	// 声明 Exchange(交换机), 如果 Exchange 不存在会创建它; 如果 Exchange 已存在, Server 会检查声明的参数和 Exchange 的真实参数是否一致
	err = ch.ExchangeDeclare(
		ExchangeName,
		"direct",
		true,
		false,
		false,
		false,
		nil,
	) 
	if err != nil {
		log.Fatalf("无法声明交换机: %v", err)
	}


	msg := "Hello Happyladysauce"
	log.Println("开始发送消息...")
	
	// 使用同步方式发送消息，确保消息成功发送
	messages := []struct {
		key  string
		num  int
	}{
		{"info", 1},
		{"debug", 2},
		{"warn", 3},
	}

	for _, m := range messages {
		messageContent := fmt.Sprintf("%s - 消息%d", msg, m.num)
		if err := send(messageContent, ch, ExchangeName, m.key); err != nil {
			log.Printf("发送消息失败 [路由键: %s]: %v", m.key, err)
		}
	}

	log.Println("消息发送完成，等待2秒确保消息被处理...")
	time.Sleep(2 * time.Second)
}

func send(msg string, ch *amqp.Channel, exchange, key string) error {
	// 发送消息
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second) // 延长超时时间
	defer cancel()
	
	log.Printf("发送消息 [路由键: %s]: %s", key, msg)
	
	// 检查上下文是否已取消
	select {
	case <-ctx.Done():
		return fmt.Errorf("上下文已取消: %v", ctx.Err())
	default:
		// 上下文未取消，继续发送消息
	}
	
	return ch.PublishWithContext(
		ctx,
		exchange,
		key, 			// exchange 为空时, key 就是队列名称
		false,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,     // 消息持久化
			Body: []byte(msg),
			ContentType: "text/plain",
			Timestamp: time.Now(),            // 添加时间戳便于调试
		},
	)
}