package rabbitmq

import (
	"log"

	"github.com/rabbitmq/amqp091-go"
)

type Consumer interface {
	Consume(queue string, handler func([]byte)) error
}

type amqpConsumer struct {
	ch *amqp091.Channel
}

func NewConsumer(conn *amqp091.Connection) (Consumer, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	return &amqpConsumer{ch: ch}, nil
}

func (c *amqpConsumer) Consume(queue string, handler func([]byte)) error {
	msgs, err := c.ch.Consume(queue, "", true, false, false, false, nil)
	if err != nil {
		return err
	}

	go func() {
		for msg := range msgs {
			log.Printf("Received on %s: %s", queue, string(msg.Body))
			handler(msg.Body)
		}
	}()

	return nil
}
