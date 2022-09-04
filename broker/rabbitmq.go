package broker

import (
	"github.com/pkg/errors"
	"github.com/rabbitmq/amqp091-go"
)

func RabbitMQ() (*amqp091.Connection, *amqp091.Channel, error) {
	// Create connection to RabbitMQ
	conn, err := amqp091.Dial("amqp://guest:guest@localhost:5672")
	if err != nil {
		return nil, nil, errors.Wrap(err, "Failed to connect RabbitMQ")
	}

	// Create channel to send data
	ch, err := conn.Channel()
	if err != nil {
		return nil, nil, errors.Wrap(err, "Failed to get channel")
	}

	return conn, ch, nil
}
