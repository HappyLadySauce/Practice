package mq

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

// first define some constants for queue and exchange names
func (m *MQ) DeclareUserQueue(userID string) (error) {
	// Declare a durable queue for the user with a message length limit of 1000 bytes.
	q, err := m.ch.QueueDeclare(
		UserQueuePrefix + userID,
		true,
		false,
		false,
		false,
		amqp.Table{
			"x-message-length": 1000,
		},
	)
	if err != nil {
		m.ch.Close()
		return fmt.Errorf("failed to declare queue: %v", err)
	}
	
	// Bind the queue to the topic exchange with an inbox routing key for this user.
	// 将队列绑定到 topic 交换机，路由键为 chat.inbox.{userID}。
	err = m.ch.QueueBind(
		q.Name,
		InboxRoutingKey(userID),
		ChatExchangeName,
		false,
		nil,
	)
	if err != nil {
		m.ch.Close()
		return fmt.Errorf("failed to bind queue: %v", err)
	}

	return nil
}
