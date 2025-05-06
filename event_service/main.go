package main

import (
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

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/streadway/amqp"
	grpc "google.golang.org/grpc"
)

func main() {
	// Postgres
	db, err := sql.Open("postgres", "host=localhost port=5432 user=postgres password=1234 dbname=user_service sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	repo := repository.NewPostgresEventRepository(db)

	// RabbitMQ
	conn, err := amqp.Dial("amqp://user:1234@localhost:5672/")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	pub, _ := rabbitmq.NewPublisher(conn)
	cons, _ := rabbitmq.NewConsumer(conn)

	uc := usecase.NewEventUseCase(repo, pub)

	// go cons.Consume("events_queue", func(b []byte) {
	// 	log.Printf("Consumed (events_queue): %s", b)
	// })

	err = cons.Consume("bet.created", func(b []byte) {
		var bet domain.Bet
		err := json.Unmarshal(b, &bet)
		if err != nil {
			log.Printf(" Failed to decode bet.created: %v", err)
			return
		}

		log.Printf("Received bet.created event: %+v", bet)
	})
	if err != nil {
		log.Fatalf(" Consumer bet.created error: %v", err)
	}
	err = cons.Consume("bet.updated", func(b []byte) {
		var bet domain.Bet
		if err := json.Unmarshal(b, &bet); err != nil {
			log.Printf(" Failed to decode bet.updated: %v", err)
			return
		}
		log.Printf(" Received bet.updated event: %+v", bet)
	})
	if err != nil {
		log.Fatalf(" Consumer bet.updated error: %v", err)
	}

	err = cons.Consume("bet.deleted", func(b []byte) {
		var bet domain.Bet
		if err := json.Unmarshal(b, &bet); err != nil {
			log.Printf(" Failed to decode bet.deleted: %v", err)
			return
		}
		log.Printf(" Received bet.deleted event: %+v", bet)
	})
	if err != nil {
		log.Fatalf(" Consumer bet.deleted error: %v", err)
	}

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

	go func() {
		log.Println("REST on :8080")
		log.Fatal(http.ListenAndServe(":8080", r))
	}()

	// gRPC server
	lis, err := net.Listen("tcp", ":50053")
	if err != nil {
		log.Fatalf("grpc listen: %v", err)
	}
	grpcSrv := grpc.NewServer()
	proto.RegisterEventServiceServer(grpcSrv, eventsvc.NewGRPCServer(uc))
	log.Println("gRPC on :50053")
	log.Fatal(grpcSrv.Serve(lis))
}
