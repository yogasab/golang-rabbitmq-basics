package main

import (
	"log"

	"github.com/pkg/errors"
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

	// Create queue to reserve the sent data
	queue, err := ch.QueueDeclare("hello", false, false, false, false, nil)
	if err != nil {
		panic(errors.Wrap(err, "Failed to declare queue"))
	}

	// Listen data from sender
	msgs, err := ch.Consume(queue.Name, "", true, false, false, false, nil)
	if err != nil {
		panic(errors.Wrap(err, "Failed to register consumer"))
	}

	var forever chan struct{}
	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
