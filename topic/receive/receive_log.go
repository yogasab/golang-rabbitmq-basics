package main

import (
	"log"
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

	q, err := ch.QueueDeclare("", false, false, true, false, nil)
	if err != nil {
		panic(errors.Wrap(err, "failed to declare queue"))
	}

	err = ch.ExchangeDeclare("logs_topic", amqp091.ExchangeTopic, true, false, false, false, nil)
	if err != nil {
		panic(errors.Wrap(err, "failed to declare exchange"))
	}

	for _, v := range os.Args[1:] {
		err = ch.QueueBind(q.Name, v, "logs_topic", false, nil)
		if err != nil {
			panic(errors.Wrap(err, "failed to bind queue"))
		}
	}

	msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		panic(errors.Wrap(err, "failed to consume queue"))
	}

	forever := make(chan struct{})

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
		}
	}()

	log.Printf("[*] Waiting for messages. To Exit press CTRL+C\n")
	<-forever
}
