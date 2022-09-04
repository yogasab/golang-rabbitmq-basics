package main

import (
	"bytes"
	"log"
	"time"

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

	queue, err := ch.QueueDeclare("task_queue", true, false, false, false, nil)
	if err != nil {
		panic(errors.Wrap(err, "Failed to declare a queue"))
	}

	err = ch.Qos(1, 0, false)
	if err != nil {
		panic(errors.Wrap(err, "failed to set QOS"))
	}

	msgs, err := ch.Consume(queue.Name, "", false, false, false, false, nil)
	forever := make(chan struct{})

	go func() {
		for d := range msgs {
			log.Printf("Received message %s", d.Body)
			dotCount := bytes.Count(d.Body, []byte("."))
			t := time.Duration(dotCount)
			time.Sleep(t * time.Second)
			log.Println("Done !")
			d.Ack(false)
		}
	}()

	log.Printf("[*] Waiting for messages. To Exit press CTRL+C\n")
	<-forever
}
