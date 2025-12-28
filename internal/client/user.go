package client

import (
	"context"
	"errors"

	"github.com/savanrajyaguru/ecommerce-go-order-service/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type UserClient struct {
	conn *grpc.ClientConn
}

func NewUserClient() (*UserClient, error) {
	// In a real scenario, we use the generated stub.
	// conn, err := grpc.NewClient(config.AppConfig.Services.UserService, grpc.WithTransportCredentials(insecure.NewCredentials()))
	// For now, we stub the connection since we lack the generated code.

	// We still dial to show intent, but can't create the client.
	conn, err := grpc.NewClient(config.AppConfig.Services.UserService, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &UserClient{conn: conn}, nil
}

func (c *UserClient) ValidateUser(ctx context.Context, token string) (string, error) {
	// STUBBED IMPLEMENTATION
	// Real implementation would be:
	// client := userpb.NewUserServiceClient(c.conn)
	// resp, err := client.ValidateUser(ctx, &userpb.ValidateUserRequest{Token: token})
	// return resp.UserId, err

	if token == "" {
		return "", errors.New("invalid token")
	}
	return "user-123", nil
}
