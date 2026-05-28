package main

import (
	"rabbitmq/mq"
)

const userID = "1234567890"

func main() {
	MQ, err := mq.Init("amqp://rabbitmq:rabbitmq@100.100.100.2:5672/")
	if err != nil {
		panic(err)
	}
	defer MQ.Close()

	// Declare a queue for the user and bind it to the exchange.
	if err := MQ.DeclareUserQueue(userID); err != nil {
		panic(err)
	}
}
