package mq

import (
	"log"
	amqp "github.com/rabbitmq/amqp091-go"
)

var conn *amqp.Connection

func Init(url string) {
	// Open a connection to the RabbitMQ server.
	conn, err := amqp.Dial(url)
	if err != nil {
		log.Fatal("failed to connect to RabbitMQ:", err)
	}

	// Open a channel on the connection.
	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		log.Fatal("failed to open a channel:", err)
	}

	// Declare a durable direct exchange named "chat.direct".
	if err := ch.ExchangeDeclare(
		"chat.direct", // name
		"direct",      // type
		true,          // durable
		false,         // auto-deleted
		false,         // internal
		false,         // no-wait
		nil,           // arguments
	); err != nil {
		ch.Close()
		conn.Close()
		log.Fatal("failed to declare exchange:", err)
	}
}