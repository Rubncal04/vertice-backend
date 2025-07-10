package repository

import (
	"context"
	"vertice-backend/config"
	"vertice-backend/internal/domain"

	"gorm.io/gorm"
)

type OrderGormRepository struct {
	db *gorm.DB
}

func NewOrderGormRepository() domain.OrderRepository {
	return &OrderGormRepository{db: config.DB}
}

func (r *OrderGormRepository) Create(ctx context.Context, order *domain.Order) error {
	return r.db.WithContext(ctx).Create(order).Error
}

func (r *OrderGormRepository) FindByIDAndUserID(ctx context.Context, id, userID uint) (*domain.Order, error) {
	var order domain.Order
	err := r.db.WithContext(ctx).
		Preload("Items.Product").
		Preload("User").
		Where("id = ? AND user_id = ?", id, userID).
		First(&order).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *OrderGormRepository) FindByUserID(ctx context.Context, userID uint) ([]*domain.Order, error) {
	var orders []*domain.Order
	err := r.db.WithContext(ctx).
		Preload("Items.Product").
		Preload("User").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&orders).Error
	if err != nil {
		return nil, err
	}
	return orders, nil
}

func (r *OrderGormRepository) Update(ctx context.Context, order *domain.Order, userID uint) error {
	return r.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", order.ID, userID).
		Save(order).Error
}

func (r *OrderGormRepository) Delete(ctx context.Context, id, userID uint) error {
	return r.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", id, userID).
		Delete(&domain.Order{}).Error
}
