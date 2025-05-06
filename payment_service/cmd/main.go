// cmd/main.go
package main

import (
	"database/sql"
	"log"
	"net"

	_ "github.com/lib/pq"

	"github.com/streadway/amqp"

	paymentgrpc "muchway/payment_service/grpc"
	"muchway/payment_service/pb"
	"muchway/payment_service/rabbitmq"
	"muchway/payment_service/repository/postgres"
	"muchway/payment_service/usecase"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	// Initialize PostgreSQL connection directly
	db, err := sql.Open("postgres", "host=localhost port=5432 user=postgres password=1234 dbname=user_service sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Initialize RabbitMQ connection
	rabbitConn, err := amqp.Dial("amqp://user:1234@localhost:5672/")
	if err != nil {
		log.Fatal("Failed to connect to RabbitMQ:", err)
	}
	defer rabbitConn.Close()

	// Setup repository and usecase
	repo := postgres.NewPostgresPaymentRepository(db)
	uc := usecase.NewPaymentUsecase(repo)

	// Start RabbitMQ consumer
	go rabbitmq.StartConsumer(rabbitConn, uc, "order_events")

	// Start gRPC server
	server := grpc.NewServer()
	pb.RegisterPaymentServiceServer(server, paymentgrpc.NewPaymentServer(uc))
	reflection.Register(server)

	lis, err := net.Listen("tcp", ":50054")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Println("Payment gRPC service running on port 50054")
	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
