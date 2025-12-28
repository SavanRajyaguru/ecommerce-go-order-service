package models

import (
	"time"

	"gorm.io/gorm"
)

type OrderStatus string

const (
	OrderStatusCreated        OrderStatus = "CREATED"
	OrderStatusPaymentPending OrderStatus = "PAYMENT_PENDING"
	OrderStatusConfirmed      OrderStatus = "CONFIRMED"
	OrderStatusCancelled      OrderStatus = "CANCELLED"
)

type Order struct {
	ID            string         `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	UserID        string         `gorm:"index;not null" json:"user_id"`
	TotalAmount   float64        `gorm:"not null" json:"total_amount"`
	Status        OrderStatus    `gorm:"index;not null;default:CREATED" json:"status"`
	Items         []OrderItem    `gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE" json:"items"`
	TransactionID string         `json:"transaction_id"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

type OrderItem struct {
	ID        string  `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	OrderID   string  `gorm:"index;not null" json:"order_id"`
	ProductID string  `gorm:"not null" json:"product_id"`
	Quantity  int     `gorm:"not null" json:"quantity"`
	UnitPrice float64 `gorm:"not null" json:"unit_price"`
}
