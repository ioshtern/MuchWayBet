package rabbitmq

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

type Consumer interface {
	Consume(queue string, handler func([]byte)) error
}

type amqpConsumer struct {
	ch *amqp.Channel
}

func NewConsumer(conn *amqp.Connection) (Consumer, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	return &amqpConsumer{ch: ch}, nil
}

func (c *amqpConsumer) Consume(queue string, handler func([]byte)) error {
	if c.ch == nil {
		return fmt.Errorf("channel is nil, cannot consume from queue: %s", queue)
	}

	// ⬇️ DECLARE the queue (idempotent if already exists)
	_, err := c.ch.QueueDeclare(
		queue, // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return fmt.Errorf("queue declare %q: %w", queue, err)
	}

	log.Printf("Subscribing to queue: %s", queue)
	msgs, err := c.ch.Consume(
		queue, // queue
		"",    // consumer tag
		true,  // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		return fmt.Errorf("consume %q: %w", queue, err)
	}

	log.Printf("Successfully subscribed to queue: %s", queue)
	go func() {
		for m := range msgs {
			handler(m.Body)
		}
	}()
	return nil
}
