package main

import (
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	// 连接到 RabbitMQ
	conn, err := amqp.Dial("amqp://test:test@100.100.100.5:5672")
	if err != nil {
		log.Fatalf("无法连接到 RabbitMQ: %v", err)
	}
	defer conn.Close()

	// 创建 Channel
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("无法创建通道: %v", err)
	}
	defer ch.Close()

	// 声明队列 - 使用与生产者一致的参数
	_, err = ch.QueueDeclare(
		"hello",		// 队列名称
		true,       // 持久化
		true,       // 自动删除
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("无法声明队列: %v", err)
	}

	go receive("hello", ch, 1)
	go receive("hello", ch, 2)
	go receive("hello", ch, 3)

	// 使用time.Sleep替代死锁的channel操作
	log.Println("消费者已启动，等待消息...")
	for {
		time.Sleep(1 * time.Second)
	}
}

type Consumer <- chan amqp.Delivery

func receive(queueName string, ch *amqp.Channel, flag int) {
	// 消费消息 - 简化为单个goroutine，避免过多嵌套
	msgs, err := ch.Consume(
		queueName,
		"consumer-"+string(rune('0'+flag)), // 消费者标签
		false,     // 关闭自动Ack，手动确认
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Printf("消费者 %d 无法消费消息: %v", flag, err)
		return
	}

	log.Printf("消费者 %d 已准备好接收消息", flag)
	for delivery := range msgs {
		log.Printf("ID: %d 收到消息: [%s]", flag, delivery.Body)
		// 手动确认消息已处理
		if err := delivery.Ack(false); err != nil {
			log.Printf("消费者 %d 确认消息失败: %v", flag, err)
		}
	}
}