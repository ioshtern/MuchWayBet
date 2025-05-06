package main

import (
	"database/sql"
	"log"
	"net"

	grpcServer "muchway/user_service/grpc"
	pb "muchway/user_service/proto/userpb"
	"muchway/user_service/rabbitmq"
	"muchway/user_service/repository/postgres"
	"muchway/user_service/usecase"

	_ "github.com/lib/pq"
	"github.com/streadway/amqp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	// PostgreSQL подключение
	db, err := sql.Open("postgres", "host=localhost port=5432 user=postgres password=1234 dbname=user_service sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// RabbitMQ подключение
	conn, err := amqp.Dial("amqp://user:1234@localhost:5672/")
	if err != nil {
		log.Fatal("Failed to connect to RabbitMQ:", err)
	}
	defer conn.Close()

	// Инициализация
	userRepo := postgres.NewPostgresUserRepository(db)

	publisher, err := rabbitmq.NewPublisher(conn)
	if err != nil {
		log.Fatal("Failed to create RabbitMQ publisher:", err)
	}

	userUsecase := usecase.NewUserUsecase(userRepo, publisher)

	// gRPC сервер
	server := grpc.NewServer()
	pb.RegisterUserServiceServer(server, grpcServer.NewUserServer(userUsecase))
	reflection.Register(server)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Println("User gRPC service running on port 50051")
	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
