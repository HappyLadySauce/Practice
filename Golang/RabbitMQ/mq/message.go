package mq

import (
	"encoding/json"
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Define a struct for the message format.
type Message struct {
	MsgID     string `json:"msg_id"`
	From      string `json:"from"`
	To        string `json:"to"`
	Content   string `json:"content"`
	Timestamp int64  `json:"timestamp"`
}


func (m *MQ) SendMessage(msg *Message) (error) {
	// Marshal the message struct to JSON.
	body, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %v", err)
	}

	// Publish the message to the "chat.direct" exchange with the routing key equal to the recipient user ID.
	if err := m.ch.Confirm(false); err != nil {
		return fmt.Errorf("failed to put channel in confirm mode: %v", err)
	}
	// Publish the message to the exchange with the routing key equal to the recipient user ID.
	confirms := m.ch.NotifyPublish(make(chan amqp.Confirmation, 1))

	// Publish the message to the "chat.direct" exchange with the routing key equal to the recipient user ID.
	err = m.ch.Publish(
		ChatExchangeName,
		msg.To,
		true,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			DeliveryMode: amqp.Persistent,
			Body: body,
			MessageId: msg.MsgID,
			Timestamp: time.Now(),
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish message: %v", err)
	}

	if confirmed := <-confirms; !confirmed.Ack {
		return fmt.Errorf("failed to confirm message delivery")
	}

	return nil
}

func (m *MQ) ConsumeMessages(userID string) (<-chan Message, error) {
	// Consume messages from the user's queue.
	msgs, err := m.ch.Consume(
		UserQueuePrefix + userID,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to consume messages: %v", err)
	}

	// Create a channel to send the messages to the caller.
	messageChan := make(chan Message)

	// Start a goroutine to read messages from the RabbitMQ channel and send them to the caller.
	go func() {
		for d := range msgs {
			var msg Message
			if err := json.Unmarshal(d.Body, &msg); err != nil {
				fmt.Printf("failed to unmarshal message: %v", err)
				continue
			}
			messageChan <- msg
		}
	}()

	return messageChan, nil
}