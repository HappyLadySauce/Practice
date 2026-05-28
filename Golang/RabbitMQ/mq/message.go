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