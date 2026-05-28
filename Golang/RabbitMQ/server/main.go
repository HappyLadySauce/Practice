package main

import (
	"rabbitmq/mq"
)

func main() {
	MQ, err := mq.Init("amqp://rabbitmq:rabbitmq@100.100.100.2:5672/")
	if err != nil {
		panic(err)
	}
	defer MQ.Close()
}