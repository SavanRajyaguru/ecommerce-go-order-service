package client

import (
	"context"
	"errors"
	"fmt"

	inventorypb "github.com/savanrajyaguru/ecommerce-go-inventory-service/proto"
	"github.com/savanrajyaguru/ecommerce-go-order-service/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type InventoryClient struct {
	client inventorypb.InventoryServiceClient
}

func NewInventoryClient() (*InventoryClient, error) {
	conn, err := grpc.NewClient(config.AppConfig.Services.InventoryService, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	c := inventorypb.NewInventoryServiceClient(conn)
	return &InventoryClient{client: c}, nil
}

func (c *InventoryClient) ReserveStock(ctx context.Context, productID string, quantity int) (bool, error) {
	// Proto uses uint32 for ID and int32 for quantity.
	// We need to convert string ID to uint32?
	// The models/order.go uses string for ProductID.
	// The inventory proto uses uint32. This is a mismatch.
	// We should probably parse it if possible, or assume the string is numeric.
	// For "Production Ready", we handle errors.

	// Assuming ProductID is parsable or we mocked it in previous prompt.
	// Let's assume we send 1 for now if parsing fails? Or try to scan.
	// But `fmt.Sscanf` is standard.

	var pid uint32
	_, err := fmt.Sscanf(productID, "%d", &pid)
	if err != nil {
		// Log warning, fallback to 0 or error?
		// If ID is a UUID string, we can't send it to a uint32 field.
		// The Inventory Service proto defined `uint32 product_id`.
		// If our Order Service uses UUIDs, we have a problem.
		// But in a microservice demo, maybe IDs are integers?
		// User didn't specify. I will assume integer strings.
		return false, errors.New("invalid product id format for inventory service")
	}

	req := &inventorypb.ReserveStockRequest{
		ProductId: pid,
		Quantity:  int32(quantity),
	}

	resp, err := c.client.ReserveStock(ctx, req)
	if err != nil {
		return false, err
	}

	return resp.Success, nil
}
