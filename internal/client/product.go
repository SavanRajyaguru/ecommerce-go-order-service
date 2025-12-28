package client

import (
	"context"
	"errors"

	"github.com/savanrajyaguru/ecommerce-go-order-service/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ProductClient struct {
	conn *grpc.ClientConn
}

func NewProductClient() (*ProductClient, error) {
	conn, err := grpc.NewClient(config.AppConfig.Services.ProductService, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return &ProductClient{conn: conn}, nil
}

func (c *ProductClient) ValidateProduct(ctx context.Context, productID string) (float64, error) {
	// STUBBED IMPLEMENTATION
	// Real:
	// client := productpb.NewProductServiceClient(c.conn)
	// resp, err := client.ValidateProduct(ctx, &ValidateProductRequest{Id: productID})

	if productID == "invalid" {
		return 0, errors.New("product not found")
	}

	// Dummy price
	return 100.0, nil
}
