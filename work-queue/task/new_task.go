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
