package rabbitmq

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

func StartConsumer(amqpURL, queueName string) error {
	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		return fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to open a channel: %w", err)
	}

	msgs, err := ch.Consume(
		queueName,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to register a consumer: %w", err)
	}

	go func() {
		log.Println("RabbitMQ consumer started...")
		for d := range msgs {
			log.Printf("Received message: %s", d.Body)
		}
	}()

	return nil
}
