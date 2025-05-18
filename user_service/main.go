package main

import (
	"context"
	"database/sql"
	"log"
	"net"

	"muchway/user_service/email"
	grpcServer "muchway/user_service/grpc"
	pb "muchway/user_service/proto/userpb"
	"muchway/user_service/rabbitmq"
	"muchway/user_service/repository/postgres"
	"muchway/user_service/usecase"

	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"github.com/streadway/amqp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	db, err := sql.Open("postgres", "host=localhost port=5433 user=postgres password=3052 dbname=muchwaybet sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatal("Failed to connect to RabbitMQ:", err)
	}
	defer conn.Close()

	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	if _, err := redisClient.Ping(context.Background()).Result(); err != nil {
		log.Fatal("Failed to connect to Redis:", err)
	}

	userRepo := postgres.NewPostgresUserRepository(db)

	publisher, err := rabbitmq.NewPublisher(conn)
	if err != nil {
		log.Fatal("Failed to create RabbitMQ publisher:", err)
	}

	// Initialize email service with Gmail credentials
	emailConfig := email.Config{
		SMTPHost:     "smtp.gmail.com",
		SMTPPort:     "587",
		SenderEmail:  "nbekzat251@gmail.com",
		SenderPasswd: "flza vhbo uwlj oizn",
	}
	emailService := email.NewEmailService(emailConfig)

	userUsecase := usecase.NewUserUsecase(userRepo, publisher, redisClient, emailService)

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
