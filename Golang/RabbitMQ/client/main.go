package main

import (
	"log"
	"fmt"
	"rabbitmq/mq"

	"github.com/google/uuid"
)

var userID string

func main() {
	MQ, err := mq.Init("amqp://rabbitmq:rabbitmq@100.100.100.2:5672/")
	if err != nil {
		panic(err)
	}
	defer MQ.Close()

	fmt.Print("Enter user ID: ")
	fmt.Scanf("%s", &userID)
	fmt.Printf("user ID: %s\n", userID)

	// Declare a queue for the user and bind it to the exchange.
	if err := MQ.DeclareUserQueue(userID); err != nil {
		panic(err)
	}

	// Start consuming messages for the user.
	go func() {
		msgChan, err := MQ.ConsumeMessages(userID)
		if err != nil {
			log.Printf("error consuming messages for user %s: %v\n", userID, err)
		}

		// Process incoming messages.
		for msg := range msgChan {
			fmt.Printf("received message: %+v\n", msg)
		}
	}()

	// Keep the main function running to allow message processing.
	for {
		select {
		default:
			msg := recoverInfo(userID)
			if err := MQ.SendMessage(&msg); err != nil {
				log.Printf("error publishing message: %v\n", err)
			}
		}
	}
}

func recoverInfo(userID string) mq.Message {
	var to, content string
	
	fmt.Print("Enter recipient user ID: ")
	fmt.Scan(&to)
	fmt.Print("Enter message content: ")
	fmt.Scan(&content)

	return mq.Message{
		MsgID: uuid.New().String(),
		From: userID,
		To: to,
		Content: content,
	}
}