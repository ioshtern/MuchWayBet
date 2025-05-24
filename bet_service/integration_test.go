package main

import (
	"context"
	"log"
	"testing"
	"time"

	"bet_service/muchway/bet_service/proto/betpb"

	"google.golang.org/grpc"
)

var createdBetID string

func TestCreateBet_Integration(t *testing.T) {
	conn, err := grpc.Dial("localhost:50052", grpc.WithInsecure())
	if err != nil {
		t.Fatalf("failed to connect to gRPC server: %v", err)
	}
	defer conn.Close()

	client := betpb.NewBetServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &betpb.CreateBetRequest{
		Bet: &betpb.Bet{
			UserId:  "user123",
			EventId: "event456",
			Amount:  150.0,
			Odds:    2.75,
		},
	}

	resp, err := client.CreateBet(ctx, req)
	if err != nil {
		t.Fatalf("CreateBet failed: %v", err)
	}

	if resp.Bet == nil || resp.Bet.UserId != "user123" {
		t.Errorf("unexpected response: %+v", resp.Bet)
	}

	createdBetID = resp.Bet.Id
	log.Printf("CreateBet passed. Bet ID: %s", createdBetID)
}

func TestGetBetByID_Integration(t *testing.T) {
	if createdBetID == "" {
		t.Skip("CreateBet must run first")
	}

	conn, err := grpc.Dial("localhost:50052", grpc.WithInsecure())
	if err != nil {
		t.Fatalf("failed to connect to gRPC server: %v", err)
	}
	defer conn.Close()

	client := betpb.NewBetServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := client.GetBetByID(ctx, &betpb.GetBetByIDRequest{
		Id: createdBetID,
	})
	if err != nil {
		t.Fatalf("GetBetByID failed: %v", err)
	}

	if resp.Bet == nil || resp.Bet.Id != createdBetID {
		t.Errorf("unexpected bet returned: %+v", resp.Bet)
	}

	log.Printf("GetBetByID passed. Got bet with ID: %s", resp.Bet.Id)
}

func TestDeleteBet_Integration(t *testing.T) {
	if createdBetID == "" {
		t.Skip("CreateBet must run first")
	}

	conn, err := grpc.Dial("localhost:50052", grpc.WithInsecure())
	if err != nil {
		t.Fatalf("failed to connect to gRPC server: %v", err)
	}
	defer conn.Close()

	client := betpb.NewBetServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = client.DeleteBet(ctx, &betpb.DeleteBetRequest{
		Id: createdBetID,
	})
	if err != nil {
		t.Fatalf("DeleteBet failed: %v", err)
	}

	log.Printf("DeleteBet passed. Bet ID deleted: %s", createdBetID)
}
