package mq

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	UserQueuePrefix = "chat.user."

	ChatExchangeName = "chat.direct"
)

type MQ struct {
	conn *amqp.Connection
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

	// Declare a durable direct exchange named "chat.direct".
	if err := ch.ExchangeDeclare(
		ChatExchangeName, // name
		"direct",      // type
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

	return &MQ{conn: conn}, nil
}

func (m *MQ) Channel() (*amqp.Channel, error) {
	// Check if the connection is still open before trying to create a channel.
	if m.conn == nil || m.conn.IsClosed() {
		return nil, fmt.Errorf("connection is closed")
	}

	// Open a channel on the existing connection.
	ch, err := m.conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open a channel: %v", err)
	}
	return ch, nil
}