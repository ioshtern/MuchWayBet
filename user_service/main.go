package main

import (
	"database/sql"
	"log"
	grpcServer "muchway/user_service/grpc"
	pb "muchway/user_service/proto/userpb"
	"muchway/user_service/repository/postgres"
	"muchway/user_service/usecase"
	_ "muchway/user_service/usecase"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	_ "github.com/lib/pq"
)

func main() {
	db, err := sql.Open("postgres", "host=localhost port=5432 user=postgres password=1234 dbname=user_service sslmode=disable")

	if err != nil {
		log.Fatal(err)
	}

	userRepo := postgres.NewPostgresUserRepository(db)
	userUsecase := usecase.NewUserUsecase(userRepo)
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
