package mq

import (
	"log"
	amqp "github.com/rabbitmq/amqp091-go"
)

var conn *amqp.Connection

func Init(url string) {
	// create TCP connction for rabbitmq
	conn, err := amqp.Dial(url)
	if err != nil {
		log.Fatalf("amqp connction error:", err)
	}

	// create channel to declare exchange
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("amqp Channel error:", err)
	}

	// create "chat.direct" exchange
	if err := ch.ExchangeDeclare(
		"chat.direct", "direct", true, false, false, false, nil,
	); err != nil {
		log.Fatalf("amqp exchangDeclare error:", err)
	}
}