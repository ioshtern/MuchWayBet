package rabbitmq

import (
	"encoding/json"
	"github.com/rabbitmq/amqp091-go"
)

type Publisher struct {
	conn    *amqp091.Connection
	channel *amqp091.Channel
	queue   amqp091.Queue
}

func NewPublisher(amqpURL, queueName string) (*Publisher, error) {
	conn, err := amqp091.Dial(amqpURL)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	q, err := ch.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return &Publisher{
		conn:    conn,
		channel: ch,
		queue:   q,
	}, nil
}

func (p *Publisher) Publish(message any) error {
	body, err := json.Marshal(message)
	if err != nil {
		return err
	}

	err = p.channel.Publish(
		"",
		p.queue.Name,
		false,
		false,
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	return err
}

func (p *Publisher) Close() {
	p.channel.Close()
	p.conn.Close()
}
