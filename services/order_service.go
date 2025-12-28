package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/savanrajyaguru/ecommerce-go-order-service/internal/client"
	"github.com/savanrajyaguru/ecommerce-go-order-service/internal/database"
	"github.com/savanrajyaguru/ecommerce-go-order-service/internal/event"
	"github.com/savanrajyaguru/ecommerce-go-order-service/models"
	"github.com/savanrajyaguru/ecommerce-go-order-service/pkg/logger"
	"github.com/savanrajyaguru/ecommerce-go-order-service/repository"
	"go.uber.org/zap"
)

type OrderService struct {
	repo            repository.OrderRepository
	userClient      *client.UserClient
	inventoryClient *client.InventoryClient
	productClient   *client.ProductClient
	paymentClient   *client.PaymentClient
}

func NewOrderService(repo repository.OrderRepository) (*OrderService, error) {
	uc, err := client.NewUserClient()
	if err != nil {
		return nil, err
	}
	ic, err := client.NewInventoryClient()
	if err != nil {
		return nil, err
	}
	pc, err := client.NewProductClient()
	if err != nil {
		return nil, err
	}
	payc, err := client.NewPaymentClient()
	if err != nil {
		return nil, err
	}

	return &OrderService{
		repo:            repo,
		userClient:      uc,
		inventoryClient: ic,
		productClient:   pc,
		paymentClient:   payc,
	}, nil
}

type CreateOrderRequest struct {
	UserID string      `json:"user_id"`
	Items  []OrderItem `json:"items"`
	Token  string      `json:"-"` // Passed from content
}

type OrderItem struct {
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
}

func (s *OrderService) CreateOrder(ctx context.Context, req CreateOrderRequest) (*models.Order, error) {
	logger.Log.Info("Starting CreateOrder", zap.String("userID", req.UserID))

	// 1. Validate User
	if _, err := s.userClient.ValidateUser(ctx, req.Token); err != nil {
		logger.Log.Error("User validation failed", zap.Error(err))
		return nil, errors.New("user validation failed")
	}

	// 2. Process Items & Calculate Total
	var totalAmount float64
	var modelItems []models.OrderItem

	for _, item := range req.Items {
		// Validate Product & Price
		price, err := s.productClient.ValidateProduct(ctx, item.ProductID)
		if err != nil {
			return nil, fmt.Errorf("invalid product %s: %v", item.ProductID, err)
		}

		totalAmount += price * float64(item.Quantity)

		modelItems = append(modelItems, models.OrderItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			UnitPrice: price,
		})

		// 3. Reserve Stock (One by one for now)
		success, err := s.inventoryClient.ReserveStock(ctx, item.ProductID, item.Quantity)
		if err != nil || !success {
			// TODO: If we fail here, previous concurrent reservations might hang?
			// We should ideall release them. But for now return error.
			return nil, fmt.Errorf("failed to reserve stock for %s", item.ProductID)
		}
	}

	// 4. Create Order in DB
	order := &models.Order{
		UserID:      req.UserID,
		TotalAmount: totalAmount,
		Status:      models.OrderStatusCreated, // OR PAYMENT_PENDING? Prompt says "CREATED" then "PAYMENT_PENDING"
		Items:       modelItems,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Wrapped in transaction not strictly needed if single insert, but safer
	tx := database.DB.Begin()
	repoTx := s.repo.WithTrx(tx)

	if err := repoTx.CreateOrder(ctx, order); err != nil {
		tx.Rollback()
		return nil, err
	}

	// Create Event
	if event.Producer != nil {
		event.Producer.PublishEvent(ctx, "order.created", order, order.ID)
	}

	// 5. Initiate Payment
	// Update status to PAYMENT_PENDING before call?
	order.Status = models.OrderStatusPaymentPending
	repoTx.UpdateStatus(ctx, order.ID, models.OrderStatusPaymentPending)
	tx.Commit() // Commit Creation before payment call to simple avoiding locks

	// Payment Call
	txID, err := s.paymentClient.InitiatePayment(ctx, order.ID, totalAmount, req.UserID)
	if err != nil {
		// Payment Failed
		logger.Log.Error("Payment failed", zap.Error(err))
		// Cancel Order
		s.CancelOrder(ctx, order.ID, req.UserID) // This will handle status update and event
		return nil, errors.New("payment failed")
	}

	// Payment Success
	order.TransactionID = txID

	// Update to CONFIRMED
	if err := s.repo.UpdateStatus(ctx, order.ID, models.OrderStatusConfirmed); err != nil {
		logger.Log.Error("Failed to update status to confirmed", zap.Error(err))
		// Critical error (money taken but order not confirmed) -> Reconciliation needed
	}

	if event.Producer != nil {
		event.Producer.PublishEvent(ctx, "order.confirmed", order, order.ID)
	}

	return order, nil
}

func (s *OrderService) GetOrder(ctx context.Context, id string) (*models.Order, error) {
	return s.repo.GetOrder(ctx, id)
}

func (s *OrderService) GetUserOrders(ctx context.Context, userID string) ([]models.Order, error) {
	return s.repo.GetOrdersByUser(ctx, userID)
}

func (s *OrderService) GetAllOrders(ctx context.Context) ([]models.Order, error) {
	return s.repo.GetAllOrders(ctx)
}

func (s *OrderService) CancelOrder(ctx context.Context, id string, userID string) error {
	order, err := s.repo.GetOrder(ctx, id)
	if err != nil {
		return err
	}

	if order.UserID != userID {
		// Admin might bypass this check in Handler, but method needs awareness
		// If Service is generic, maybe Handler checks permissions?
		// We'll assume UserID passed is expected owner.
		// If UserID is empty (Admin), we skip check?
		if userID != "" && userID != "admin" {
			return errors.New("unauthorized")
		}
	}

	if order.Status == models.OrderStatusConfirmed {
		return errors.New("cannot cancel confirmed order")
	}

	if order.Status == models.OrderStatusCancelled {
		return nil
	}

	if err := s.repo.UpdateStatus(ctx, id, models.OrderStatusCancelled); err != nil {
		return err
	}

	if event.Producer != nil {
		event.Producer.PublishEvent(ctx, "order.cancelled", order, order.ID)
	}

	// Stock release is handled by Inventory Service consuming 'order.cancelled'.

	return nil
}
