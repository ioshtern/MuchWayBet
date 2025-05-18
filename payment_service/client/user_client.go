package client

import (
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	userpb "muchway/user_service/proto/userpb"
)

// UserClient is a client for the user service
type UserClient struct {
	client userpb.UserServiceClient
	conn   *grpc.ClientConn
}

// NewUserClient creates a new user service client
func NewUserClient(address string) (*UserClient, error) {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to user service: %w", err)
	}

	client := userpb.NewUserServiceClient(conn)
	return &UserClient{
		client: client,
		conn:   conn,
	}, nil
}

// Close closes the connection to the user service
func (c *UserClient) Close() error {
	return c.conn.Close()
}

// GetUserEmail retrieves a user's email by their ID
func (c *UserClient) GetUserEmail(userID string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Convert string ID to int64
	var id int64
	_, err := fmt.Sscanf(userID, "%d", &id)
	if err != nil {
		return "", fmt.Errorf("invalid user ID format: %w", err)
	}

	resp, err := c.client.GetUserByID(ctx, &userpb.GetUserByIDRequest{Id: id})
	if err != nil {
		return "", fmt.Errorf("failed to get user: %w", err)
	}

	if resp == nil || resp.User == nil {
		return "", fmt.Errorf("user not found")
	}

	log.Printf("Retrieved email for user %s: %s", userID, resp.User.Email)
	return resp.User.Email, nil
}
