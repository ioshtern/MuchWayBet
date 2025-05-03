package rabbitmq

import (
	"encoding/json"

	"github.com/streadway/amqp"
)

type Publisher interface {
	Publish(exchange, key string, msg interface{}) error
}

type amqpPublisher struct{ ch *amqp.Channel }

func NewPublisher(conn *amqp.Connection) (Publisher, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	return &amqpPublisher{ch: ch}, nil
}

func (p *amqpPublisher) Publish(exchange, key string, msg interface{}) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	return p.ch.Publish(exchange, key, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        data,
	})
}
