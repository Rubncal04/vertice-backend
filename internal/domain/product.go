package domain

import (
	"context"
	"time"
)

type Product struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	UserID      uint      `gorm:"not null;uniqueIndex:idx_user_code" json:"user_id"`
	User        *User     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	Code        string    `gorm:"not null;uniqueIndex:idx_user_code" json:"code"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	Stock       int       `json:"stock"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ProductRepository interface {
	Create(ctx context.Context, product *Product) error
	FindByIDAndUserID(ctx context.Context, id uint, userID uint) (*Product, error)
	FindByUserID(ctx context.Context, userID uint) ([]*Product, error)
	FindByCodeAndUserID(ctx context.Context, code string, userID uint) (*Product, error)
	Update(ctx context.Context, product *Product, userID uint) error
	Delete(ctx context.Context, id uint, userID uint) error
}
