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

	err = ch.Qos(	// 公平分发模式
		2,	// prefetch count 一个消费方最多能有多少条未ack的消息, 如果是 consumer 以 noAck(autoAck) 启动的则 server 会忽略该参数
							// <=0 时忽略该参数, 该值越小, 负载越均衡, 但是单个消费方的吞吐也越低(主要是网络IO次数太多造成的)
		0,	// 按照消息的字节来计算, 但 server 端攒够这么多字节才发送给消费方, <= 0 时忽略该参数
		false,		// global
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