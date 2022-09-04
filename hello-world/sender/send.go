package main

import (
	"fmt"
	"os"

	"github.com/pkg/errors"
	"github.com/rabbitmq/amqp091-go"
	"github.com/yogasab/golang-rabbitmq-basics/broker"
)

func main() {
	conn, ch, err := broker.RabbitMQ()
	if err != nil {
		panic(err)
	}

	defer func() {
		ch.Close()
		conn.Close()
	}()

	// Declare Queue to reserve the data we send
	queue, err := ch.QueueDeclare("hello", false, false, false, false, nil)
	if err != nil {
		panic(errors.Wrap(err, "Failed to declare queue"))
	}

	// Send data to queue
	err = ch.Publish("", queue.Name, false, false, amqp091.Publishing{
		ContentType: "text/plain",
		Body:        []byte(os.Args[1]),
	})
	if err != nil {
		panic(errors.Wrap(err, "Failed to publish message"))
	}
	// Take the second index from command
	fmt.Println("Sending Message: ...", os.Args[1])
}
