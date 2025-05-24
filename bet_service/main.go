package main

import (
	"database/sql"
	"log"
	"net"

	"bet_service/muchway/bet_service/proto/betpb"
	redisrepo "bet_service/repository"
	repo "bet_service/repository/postgres"
	betgrpc "bet_service/transport/grpc"
	"bet_service/transport/rabbitmq"
	"bet_service/usecase"

	_ "github.com/lib/pq"
	"github.com/rabbitmq/amqp091-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	db, err := sql.Open("postgres", "postgres://postgres:1234@localhost:5432/user_service?sslmode=disable")

	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("Database is not reachable:", err)
	}
	log.Println(" Connected to PostgreSQL.")

	rabbitConn, err := amqp091.Dial("amqp://user:1234@localhost:5672/")
	if err != nil {
		log.Fatal("Failed to connect to RabbitMQ:", err)
	}
	defer rabbitConn.Close()
	log.Println(" Connected to RabbitMQ.")

	redisrepo.InitRedisClient("localhost:6379", "", 0)
	log.Println(" Connected to Redis.")

	betRepo := repo.NewPostgresBetRepository(db)
	publisher, err := rabbitmq.NewPublisher(rabbitConn)
	if err != nil {
		log.Fatal("Failed to create RabbitMQ publisher:", err)
	}
	betUsecase := usecase.NewBetUsecase(betRepo, publisher)
	betServer := betgrpc.NewBetServer(betUsecase, publisher)
	consumer, err := rabbitmq.NewConsumer(rabbitConn)
	if err != nil {
		log.Fatal("Failed to create RabbitMQ consumer:", err)
	}

	_ = consumer.Consume("user.created", func(b []byte) {
		log.Printf(" New user created: %s", b)
	})

	_ = consumer.Consume("user.logged_in", func(b []byte) {
		log.Printf(" User logged in: %s", b)
	})

	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	betpb.RegisterBetServiceServer(grpcServer, betServer)
	reflection.Register(grpcServer)

	log.Println(" BetService gRPC server running on port 50052")

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC server: %v", err)
	}
}
