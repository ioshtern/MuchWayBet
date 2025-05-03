package rabbitmq

import "github.com/streadway/amqp"

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
	msgs, err := c.ch.Consume(queue, "", true, false, false, false, nil)
	if err != nil {
		return err
	}
	go func() {
		for m := range msgs {
			handler(m.Body)
		}
	}()
	return nil
}
