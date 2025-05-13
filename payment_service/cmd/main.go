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
	redisRepo "muchway/payment_service/repository/redis"
	"muchway/payment_service/usecase"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	db, err := sql.Open("postgres", "host=localhost port=5433 user=postgres password=3052 dbname=muchwaybet sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rabbitConn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatal("Failed to connect to RabbitMQ:", err)
	}
	defer rabbitConn.Close()

	redisRepo.InitRedisClient("localhost:6379", "", 0)
	log.Println("Redis client initialized")

	postgresRepo := postgres.NewPostgresPaymentRepository(db)

	repo := redisRepo.NewRedisPaymentRepository(postgresRepo)

	uc := usecase.NewPaymentUsecase(repo)

	go rabbitmq.StartConsumer(rabbitConn, uc, "order_events")

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
