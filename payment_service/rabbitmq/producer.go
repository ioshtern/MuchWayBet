package rabbitmq

import (
	"github.com/streadway/amqp"
)

func NewConnection(url string) (*amqp.Connection, error) {
	return amqp.Dial(url)
}

func NewChannel(conn *amqp.Connection) (*amqp.Channel, error) {
	return conn.Channel()
}
