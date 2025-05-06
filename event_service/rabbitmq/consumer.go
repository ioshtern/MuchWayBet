package rabbitmq

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

type Consumer interface {
	Consume(queue string, handler func([]byte)) error
}

type amqpConsumer struct{ ch *amqp.Channel }

func NewConsumer(conn *amqp.Connection) (Consumer, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	return &amqpConsumer{ch: ch}, nil
}

func (c *amqpConsumer) Consume(queue string, handler func([]byte)) error {
	if c.ch == nil {
		return fmt.Errorf(" channel is nil, cannot consume from queue: %s", queue)
	}

	log.Printf(" Subscribing to queue: %s", queue)

	msgs, err := c.ch.Consume(
		queue,
		"",
		true,  // auto-ack
		false, // not exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		log.Printf(" Failed to consume from queue %s: %v", queue, err)
		return err
	}

	log.Printf(" Successfully subscribed to queue: %s", queue)

	go func() {
		for m := range msgs {
			handler(m.Body)
		}
	}()

	return nil
}
