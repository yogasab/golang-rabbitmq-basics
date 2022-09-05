package main

import (
	"flag"
	"log"

	"github.com/pkg/errors"
	"github.com/rabbitmq/amqp091-go"
	"github.com/yogasab/golang-rabbitmq-basics/broker"
)

func main() {
	state := flag.String("s", "Jakarta", "State name")
	country := flag.String("c", "Indonesia", "Country name")
	flag.Parse()

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

	err = ch.ExchangeDeclare("logs_header", amqp091.ExchangeHeaders, true, false, false, false, nil)
	if err != nil {
		panic(errors.Wrap(err, "failed to declare exchange"))
	}

	table := amqp091.Table{
		"negara": *country,
		"state":  *state,
	}

	err = ch.QueueBind(q.Name, "", "logs_header", false, table)
	if err != nil {
		panic(errors.Wrap(err, "failed to bind queue"))
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
