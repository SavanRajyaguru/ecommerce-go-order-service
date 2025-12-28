package client

import (
	"context"
	"errors"

	"github.com/savanrajyaguru/ecommerce-go-order-service/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type PaymentClient struct {
	conn *grpc.ClientConn
}

func NewPaymentClient() (*PaymentClient, error) {
	conn, err := grpc.NewClient(config.AppConfig.Services.PaymentService, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return &PaymentClient{conn: conn}, nil
}

func (c *PaymentClient) InitiatePayment(ctx context.Context, orderID string, amount float64, userID string) (string, error) {
	// STUBBED IMPLEMENTATION
	if amount <= 0 {
		return "", errors.New("invalid amount")
	}
	return "tx-123456", nil
}
