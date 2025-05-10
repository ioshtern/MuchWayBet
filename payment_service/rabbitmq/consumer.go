package rabbitmq

import (
	"encoding/json"
	"log"
	"muchway/payment_service/usecase"

	"github.com/streadway/amqp"
)

type PaymentEvent struct {
	OrderID     string  `json:"order_id"`
	UserID      string  `json:"user_id"`
	Amount      float64 `json:"amount"`
	PaymentType string  `json:"payment_type"` // "deposit" or "withdraw"
}

func StartConsumer(conn *amqp.Connection, uc *usecase.PaymentUsecase, queue string) {
	ch, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
	}

	_, err = ch.QueueDeclare(
		queue,
		true,
		false,
		false,
		false,
		nil,
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
		var ev PaymentEvent
		err := json.Unmarshal(d.Body, &ev)
		if err != nil {
			log.Println("Failed to unmarshal payment event:", err)
			continue
		}

		if ev.PaymentType == "" {
			ev.PaymentType = "deposit"
		}

		if ev.PaymentType != "deposit" && ev.PaymentType != "withdraw" {
			log.Printf("Invalid payment type: %s. Must be 'deposit' or 'withdraw'", ev.PaymentType)
			continue
		}

		log.Printf("Processing %s of %.2f for user %s (Order: %s)",
			ev.PaymentType, ev.Amount, ev.UserID, ev.OrderID)

		payment, err := uc.ProcessPayment(ev.UserID, ev.Amount, ev.PaymentType)
		if err != nil {
			log.Printf("Failed to process payment: %v", err)
			continue
		}

		updatedPayment, err := uc.GetPaymentByID(payment.ID)
		if err != nil {
			log.Printf("Error retrieving payment after processing: %v", err)
		} else {
			log.Printf("Payment status after processing: %s", updatedPayment.Status)
		}

		log.Printf("Successfully processed %s payment (ID: %s) for order: %s",
			ev.PaymentType, payment.ID, ev.OrderID)
	}
}
