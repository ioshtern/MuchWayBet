package rabbitmq

import (
	"bet_service/domain"
	"encoding/json"
	"log"

	"github.com/rabbitmq/amqp091-go"
)

type Publisher struct {
	channel *amqp091.Channel
	queue   amqp091.Queue
}

func NewPublisher(conn *amqp091.Connection) (*Publisher, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	q, err := ch.QueueDeclare(
		"bet.created",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return &Publisher{channel: ch, queue: q}, nil
}

func (p *Publisher) PublishBetCreated(bet *domain.Bet) error {
	body, err := json.Marshal(bet)
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
	if err != nil {
		log.Printf("Failed to publish message: %v", err)
		return err
	}

	log.Println("Published bet.created event:", string(body))
	return nil
}
func (p *Publisher) PublishBetUpdated(bet *domain.Bet) error {
	return p.publish("bet.updated", bet)
}

func (p *Publisher) PublishBetDeleted(bet *domain.Bet) error {
	return p.publish("bet.deleted", bet)
}

func (p *Publisher) publish(queueName string, bet *domain.Bet) error {
	body, err := json.Marshal(bet)
	if err != nil {
		return err
	}

	_, err = p.channel.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	err = p.channel.Publish(
		"",
		queueName,
		false,
		false,
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		log.Printf("Failed to publish to %s: %v", queueName, err)
		return err
	}

	log.Printf(" Published %s event: %s", queueName, string(body))
	return nil
}
