package main

import (
	"fmt"
	"os"

	"github.com/pkg/errors"
	"github.com/rabbitmq/amqp091-go"
)

func main() {
	conn, err := amqp091.Dial("amqp://guest:guest@localhost:5672")
	if err != nil {
		panic(errors.Wrap(err, "Failed to connect to RabbitMQ"))
	}

	ch, err := conn.Channel()
	if err != nil {
		panic(errors.Wrap(err, "Failed to open a channel"))
	}
	defer ch.Close()

	queue, err := ch.QueueDeclare("task_queue", true, false, false, false, nil)
	if err != nil {
		panic(errors.Wrap(err, "Failed to declare a queue"))
	}

	err = ch.Publish("", queue.Name, false, false, amqp091.Publishing{
		ContentType: "text/plain",
		Body:        []byte(os.Args[1]),
	})
	if err != nil {
		panic(errors.Wrap(err, "failed to publish message"))
	}

	fmt.Println("Send message:", os.Args[1])
}
