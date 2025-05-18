package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"muchway/event_service/domain"
	eventsvc "muchway/event_service/grpc"
	"muchway/event_service/proto"
	"muchway/event_service/rabbitmq"
	"muchway/event_service/repository"
	"muchway/event_service/usecase"
	"net"
	"net/http"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/streadway/amqp"
	grpc "google.golang.org/grpc"
)

func main() {
	ctx := context.Background()

	// Postgres
	db, err := sql.Open("postgres", "host=localhost port=5432 user=postgres password=123654789 dbname=gin sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	repo := repository.NewPostgresEventRepository(db)

	// RabbitMQ
	conn, err := amqp.Dial("amqp://guest:guest@127.0.0.1:5672/")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	pub, _ := rabbitmq.NewPublisher(conn)
	cons, _ := rabbitmq.NewConsumer(conn)

	// Redis
	rdb := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
		DB:   0,
	})
	defer rdb.Close()

	if pong, err := rdb.Ping(ctx).Result(); err != nil {
		log.Fatalf("Redis ping failed: %v", err)
	} else {
		log.Printf("Redis connected: %s", pong)
	}

	uc := usecase.NewEventUseCase(repo, pub, rdb)

	// RabbitMQ
	for _, q := range []string{"events_queue", "bet.created", "bet.updated", "bet.deleted"} {
		queue := q
		go func() {
			if err := cons.Consume(queue, func(b []byte) {
				log.Printf("Consumed (%s): %s", queue, string(b))
			}); err != nil {
				log.Printf("Consumer %s error: %v", queue, err)
			}
		}()
	}

	// HTTP
	r := mux.NewRouter()
	r.HandleFunc("/events", func(w http.ResponseWriter, req *http.Request) {
		var e domain.Event
		if err := json.NewDecoder(req.Body).Decode(&e); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		saved, err := uc.CreateEvent(req.Context(), &e)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(saved)
	}).Methods("POST")

	r.HandleFunc("/events", func(w http.ResponseWriter, req *http.Request) {
		list, err := uc.ListEvents(req.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(list)
	}).Methods("GET")

	r.HandleFunc("/events/{id}", func(w http.ResponseWriter, req *http.Request) {
		id := mux.Vars(req)["id"]
		e, err := uc.GetEvent(req.Context(), id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		json.NewEncoder(w).Encode(e)
	}).Methods("GET")

	r.HandleFunc("/events/{id}", func(w http.ResponseWriter, req *http.Request) {
		id := mux.Vars(req)["id"]
		var e domain.Event
		if err := json.NewDecoder(req.Body).Decode(&e); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		e.ID = id
		updated, err := uc.UpdateEvent(req.Context(), &e)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(updated)
	}).Methods("PUT")

	r.HandleFunc("/events/{id}", func(w http.ResponseWriter, req *http.Request) {
		id := mux.Vars(req)["id"]
		if err := uc.DeleteEvent(req.Context(), id); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}).Methods("DELETE")

	// REST
	go func() {
		log.Println("REST on :8080")
		log.Fatal(http.ListenAndServe(":8080", r))
	}()

	// GRPC
	lis, err := net.Listen("tcp", ":50053")
	if err != nil {
		log.Fatalf("grpc listen: %v", err)
	}
	grpcSrv := grpc.NewServer()
	proto.RegisterEventServiceServer(grpcSrv, eventsvc.NewGRPCServer(uc))
	log.Println("gRPC on :50053")
	log.Fatal(grpcSrv.Serve(lis))
}
