package rabbitmq

import (
	"encoding/json"
	"log"

	"github.com/streadway/amqp"
)

type Publisher struct {
	channel *amqp.Channel
}

func NewPublisher(conn *amqp.Connection) (*Publisher, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	return &Publisher{channel: ch}, nil
}

func (p *Publisher) Publish(eventName string, data interface{}) error {
	body, err := json.Marshal(data)
	if err != nil {
		return err
	}

	_, err = p.channel.QueueDeclare(
		eventName,
		true,  // durable
		false, // auto-delete
		false, // exclusive
		false, // no-wait
		nil,
	)
	if err != nil {
		return err
	}

	err = p.channel.Publish(
		"",
		eventName,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		log.Printf("Failed to publish %s: %v", eventName, err)
		return err
	}

	log.Printf("Published %s: %s", eventName, string(body))
	return nil
}
