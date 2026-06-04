package mq

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	UserQueuePrefix = "chat.user."

	ChatExchangeName = "chat.topic"

	// InboxRoutingKeyPrefix is the topic routing-key prefix for user inbox delivery.
	// Topic 模式下用户收件箱的路由键前缀。
	InboxRoutingKeyPrefix = "chat.inbox."
)

// InboxRoutingKey builds the routing key for delivering a message to a user's inbox.
// 构建投递到指定用户收件箱的路由键，例如 chat.inbox.123。
func InboxRoutingKey(userID string) string {
	return InboxRoutingKeyPrefix + userID
}

type MQ struct {
	conn *amqp.Connection
	ch   *amqp.Channel
}

func Init(url string) (*MQ, error) {
	// Open a connection to the RabbitMQ server.
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %v", err)
	}

	// Open a channel on the connection.
	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open a channel: %v", err)
	}

	// Declare a durable topic exchange named "chat.topic".
	// 声明持久化 topic 交换机。
	if err := ch.ExchangeDeclare(
		ChatExchangeName, // name
		"topic",       // type
		true,          // durable
		false,         // auto-deleted
		false,         // internal
		false,         // no-wait
		nil,           // arguments
	); err != nil {
		ch.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to declare exchange: %v", err)
	}

	return &MQ{conn: conn, ch: ch}, nil
}

func (m *MQ) Close() {
	if m.ch != nil {
		m.ch.Close()
	}
	if m.conn != nil {
		m.conn.Close()
	}
}
