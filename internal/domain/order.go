package domain

import (
	"context"
	"time"
)

type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "pending"
	OrderStatusConfirmed OrderStatus = "confirmed"
	OrderStatusShipped   OrderStatus = "shipped"
	OrderStatusDelivered OrderStatus = "delivered"
	OrderStatusCancelled OrderStatus = "cancelled"
)

type Order struct {
	ID          uint        `json:"id" gorm:"primaryKey"`
	UserID      uint        `json:"user_id" gorm:"not null"`
	User        User        `json:"user" gorm:"foreignKey:UserID"`
	Status      OrderStatus `json:"status" gorm:"type:varchar(20);default:'pending'"`
	TotalAmount float64     `json:"total_amount" gorm:"not null"`
	Items       []OrderItem `json:"items" gorm:"foreignKey:OrderID"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}

type OrderItem struct {
	ID        uint    `json:"id" gorm:"primaryKey"`
	OrderID   uint    `json:"order_id" gorm:"not null"`
	Order     Order   `json:"order" gorm:"foreignKey:OrderID"`
	ProductID uint    `json:"product_id" gorm:"not null"`
	Product   Product `json:"product" gorm:"foreignKey:ProductID"`
	Quantity  int     `json:"quantity" gorm:"not null"`
	UnitPrice float64 `json:"unit_price" gorm:"not null"`
	Subtotal  float64 `json:"subtotal" gorm:"not null"`
}

type OrderRepository interface {
	Create(ctx context.Context, order *Order) error
	FindByIDAndUserID(ctx context.Context, id, userID uint) (*Order, error)
	FindByUserID(ctx context.Context, userID uint) ([]*Order, error)
	Update(ctx context.Context, order *Order, userID uint) error
	Delete(ctx context.Context, id, userID uint) error
}
