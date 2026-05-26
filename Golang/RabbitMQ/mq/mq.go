package mq

import (
	"log"
	amqp "github.com/rabbitmq/amqp091-go"
)

var conn *amqp.Connection
var ch *amqp.Channel

func Init(url string) {
	conn, err := amqp.Dial(url)
	if err != nil {
		log.Fatalf("amqp connction error:", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("amqp Channel error:", err)
	}

	
}