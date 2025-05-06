package rabbitmq

import (
	"bet_service/domain"
	"context"
	"encoding/json"
	"log"

	"github.com/streadway/amqp"
)

type Consumer struct {
	channel *amqp.Channel
	queue   string
}

func NewConsumer(ch *amqp.Channel, queueName string) *Consumer {
	return &Consumer{
		channel: ch,
		queue:   queueName,
	}
}

func (c *Consumer) Start(ctx context.Context) error {
	msgs, err := c.channel.Consume(
		c.queue,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case d := <-msgs:
				var bet domain.Bet
				err := json.Unmarshal(d.Body, &bet)
				if err != nil {
					log.Printf("Error decoding message: %v", err)
					continue
				}
				log.Printf("Received new bet: %+v", bet)

			}
		}
	}()

	return nil
}
