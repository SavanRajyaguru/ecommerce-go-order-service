package repository

import (
	"context"

	"github.com/savanrajyaguru/ecommerce-go-order-service/models"
	"gorm.io/gorm"
)

type OrderRepository interface {
	CreateOrder(ctx context.Context, order *models.Order) error
	GetOrder(ctx context.Context, id string) (*models.Order, error)
	GetOrdersByUser(ctx context.Context, userID string) ([]models.Order, error)
	UpdateStatus(ctx context.Context, id string, status models.OrderStatus) error
	GetAllOrders(ctx context.Context) ([]models.Order, error)
	WithTrx(tx *gorm.DB) OrderRepository
}

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) WithTrx(tx *gorm.DB) OrderRepository {
	return &orderRepository{db: tx}
}

func (r *orderRepository) CreateOrder(ctx context.Context, order *models.Order) error {
	return r.db.WithContext(ctx).Create(order).Error
}

func (r *orderRepository) GetOrder(ctx context.Context, id string) (*models.Order, error) {
	var order models.Order
	// Preload items for full order details
	err := r.db.WithContext(ctx).Preload("Items").First(&order, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *orderRepository) GetOrdersByUser(ctx context.Context, userID string) ([]models.Order, error) {
	var orders []models.Order
	err := r.db.WithContext(ctx).Preload("Items").Where("user_id = ?", userID).Order("created_at desc").Find(&orders).Error
	return orders, err
}

func (r *orderRepository) UpdateStatus(ctx context.Context, id string, status models.OrderStatus) error {
	return r.db.WithContext(ctx).Model(&models.Order{}).Where("id = ?", id).Update("status", status).Error
}

func (r *orderRepository) GetAllOrders(ctx context.Context) ([]models.Order, error) {
	var orders []models.Order
	err := r.db.WithContext(ctx).Preload("Items").Order("created_at desc").Find(&orders).Error
	return orders, err
}
