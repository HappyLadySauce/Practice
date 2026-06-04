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

	// Publish to the topic exchange; routing key targets the recipient inbox.
	// 发布到 topic 交换机，路由键指向收件人收件箱。
	if err := m.ch.Confirm(false); err != nil {
		return fmt.Errorf("failed to put channel in confirm mode: %v", err)
	}
	confirms := m.ch.NotifyPublish(make(chan amqp.Confirmation, 1))

	err = m.ch.Publish(
		ChatExchangeName,
		InboxRoutingKey(msg.To),
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

	// Unbuffered: handoff to the caller completes before Ack.
	// 无缓冲：交给调用方后再 Ack，避免“已确认但未送达应用”。
	messageChan := make(chan Message)

	go func() {
		defer close(messageChan)
		for d := range msgs {
			var msg Message
			if err := json.Unmarshal(d.Body, &msg); err != nil {
				fmt.Printf("failed to unmarshal message: %v\n", err)
				// Poison message: discard, do not requeue.
				// 非法消息：丢弃，不重新入队。
				if err := d.Nack(false, false); err != nil {
					fmt.Printf("failed to nack invalid message: %v\n", err)
				}
				continue
			}

			messageChan <- msg

			// Ack only after the caller has received the message.
			// 调用方已从 channel 取走后，再向 broker 确认。
			if err := d.Ack(false); err != nil {
				fmt.Printf("failed to ack message %s: %v\n", msg.MsgID, err)
			}
		}
	}()

	return messageChan, nil
}