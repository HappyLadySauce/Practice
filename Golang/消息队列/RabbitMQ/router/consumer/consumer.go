package main

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

var ExchangeName string = "HappyLadySauce" // 与生产者保持一致

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

	// 声明 Exchange，与生产者保持一致
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

	keys := []string{"info", "debug", "warn"}

	// 声明队列 - 使用与生产者一致的参数
	q, err := ch.QueueDeclare(
		"",		// 当队列名称为空时, 会随机创建一个队列
		true,       // 持久化
		true,       // 自动删除
		true,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("无法声明队列: %v", err)
	}

	for _, key := range keys {
		err = ch.QueueBind(
			q.Name,
			key,
			ExchangeName,
			false,
			nil,
		)
		if err != nil {
			log.Fatalf("无法绑定队列到交换机: %v", err)
		}
	}

	log.Printf("队列 %s 已绑定到交换机 %s，绑定的路由键: %v", q.Name, ExchangeName, keys)

	// 启动消费者
	receive(q.Name, ch, 1)
}

type Consumer <- chan amqp.Delivery

func receive(queueName string, ch *amqp.Channel, flag int) {
	log.Printf("消费者 %d 开始监听队列: %s", flag, queueName)
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