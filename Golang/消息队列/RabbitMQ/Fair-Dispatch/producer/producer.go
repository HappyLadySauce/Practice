package main

import (
	"context"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

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
	_, err = ch.QueueDeclare(
		"hello",		// 队列名称
		true,		// 持久化
		false,		// 自动删除
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("无法声明队列: %v", err)
	}

	msg := "Hello Happyladysauce"
	log.Println("开始发送消息...")
	
	// 使用同步方式发送消息，确保消息成功发送
	for i := 0; i < 3; i++ {
		if err := send(msg+" - 消息"+string(rune('0'+i)), "hello", ch); err != nil {
			log.Printf("发送消息失败: %v", err)
		}
	}

	log.Println("消息发送完成，等待2秒确保消息被处理...")
	time.Sleep(2 * time.Second)
}

func send(msg string, queueName string, ch *amqp.Channel) error {
	// 发送消息
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second) // 延长超时时间
	defer cancel()
	
	log.Printf("发送消息: %s", msg)
	return ch.PublishWithContext(
		ctx,
		"",
		queueName, // exchange 为空时, key 就是队列名称
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