package rabbitmq

import (
	"encoding/json"
	"log"
	"muchway/payment_service/domain"
	"muchway/payment_service/usecase"

	"github.com/streadway/amqp"
)

type OrderCreatedEvent struct {
	OrderID string  `json:"order_id"`
	UserID  string  `json:"user_id"`
	Amount  float64 `json:"amount"`
}

func StartConsumer(conn *amqp.Connection, uc *usecase.PaymentUsecase, queue string) {
	ch, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
	}

	// Declare the queue (create if it doesn't exist)
	_, err = ch.QueueDeclare(
		queue, // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		log.Fatal("Failed to declare queue:", err)
	}

	msgs, err := ch.Consume(queue, "", true, false, false, false, nil)
	if err != nil {
		log.Fatal("Failed to start consuming:", err)
	}

	log.Println("Started consuming from queue:", queue)

	for d := range msgs {
		var ev OrderCreatedEvent
		err := json.Unmarshal(d.Body, &ev)
		if err != nil {
			log.Println("Failed to unmarshal order event:", err)
			continue
		}
		p := &domain.Payment{UserID: ev.UserID, Type: "deposit", Amount: ev.Amount, Status: "pending"}
		err = uc.CreatePayment(p)
		if err != nil {
			log.Println("Failed to create payment:", err)
		}
		log.Println("Processed payment for order:", ev.OrderID)
	}
}
