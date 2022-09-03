package main

import (
	"log"

	"github.com/pkg/errors"
	"github.com/rabbitmq/amqp091-go"
)

func main() {
	// Create connection to RabbitMQ
	conn, err := amqp091.Dial("amqp://guest:guest@localhost:5672")
	log.Println(err)
	if err != nil {
		panic(errors.Wrap(err, "Failed to connect RabbitMQ"))
	}
	defer conn.Close()

	// Create channel to send data
	ch, err := conn.Channel()
	if err != nil {
		panic(errors.Wrap(err, "Failed to get channel"))
	}
	defer ch.Close()

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
