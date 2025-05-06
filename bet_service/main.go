package main

import (
	"database/sql"
	"log"
	"net"

	"bet_service/muchway/bet_service/proto/betpb"
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
	db, err := sql.Open("postgres", "host=localhost port=5432 user=postgres password=1234 dbname=user_service sslmode=disable")
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

	betRepo := repo.NewPostgresBetRepository(db)
	publisher, err := rabbitmq.NewPublisher(rabbitConn)
	if err != nil {
		log.Fatal("Failed to create RabbitMQ publisher:", err)
	}
	betUsecase := usecase.NewBetUsecase(betRepo, publisher)
	betServer := betgrpc.NewBetServer(betUsecase, publisher)

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
