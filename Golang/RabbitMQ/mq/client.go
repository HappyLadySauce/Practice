package mq

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

// first define some constants for queue and exchange names
func (m *MQ) DeclareUserQueue(userID string) (error) {
	// Open a channel on the existing connection.
	ch, err := m.Channel()
	if err != nil {
		return fmt.Errorf("failed to get channel: %v", err)
	}

	// Declare a durable queue for the user with a message length limit of 1000 bytes.
	q, err := ch.QueueDeclare(
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
		ch.Close()
		return fmt.Errorf("failed to declare queue: %v", err)
	}
	
	// Bind the queue to the "chat.direct" exchange with the routing key equal to the user ID.
	err = ch.QueueBind(
		q.Name,
		userID,
		ChatExchangeName,
		false,
		nil,
	)
	if err != nil {
		ch.Close()
		return fmt.Errorf("failed to bind queue: %v", err)
	}

	return nil
}