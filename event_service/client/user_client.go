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

// GetAllUserEmails retrieves all user emails from the user service
func (c *UserClient) GetAllUserEmails(ctx context.Context) ([]string, error) {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	resp, err := c.client.GetAllUsers(ctxWithTimeout, &userpb.GetAllUsersRequest{})
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}

	if resp == nil || len(resp.Users) == 0 {
		return []string{}, nil
	}

	var emails []string
	for _, user := range resp.Users {
		if user.Email != "" {
			emails = append(emails, user.Email)
		}
	}

	log.Printf("Retrieved %d user emails", len(emails))
	return emails, nil
}
