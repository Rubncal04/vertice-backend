package repository

import (
	"context"
	"vertice-backend/config"
	"vertice-backend/internal/domain"
)

type ProductGormRepository struct{}

func NewProductGormRepository() *ProductGormRepository {
	return &ProductGormRepository{}
}

func (r *ProductGormRepository) Create(ctx context.Context, product *domain.Product) error {
	return config.DB.WithContext(ctx).Create(product).Error
}

func (r *ProductGormRepository) FindByIDAndUserID(ctx context.Context, id uint, userID uint) (*domain.Product, error) {
	var product domain.Product
	err := config.DB.WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).First(&product).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *ProductGormRepository) FindByUserID(ctx context.Context, userID uint) ([]*domain.Product, error) {
	var products []*domain.Product
	err := config.DB.WithContext(ctx).Where("user_id = ?", userID).Find(&products).Error
	if err != nil {
		return nil, err
	}
	return products, nil
}

func (r *ProductGormRepository) FindByCodeAndUserID(ctx context.Context, code string, userID uint) (*domain.Product, error) {
	var product domain.Product
	err := config.DB.WithContext(ctx).Where("code = ? AND user_id = ?", code, userID).First(&product).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *ProductGormRepository) Update(ctx context.Context, product *domain.Product, userID uint) error {
	return config.DB.WithContext(ctx).Model(&domain.Product{}).
		Where("id = ? AND user_id = ?", product.ID, userID).
		Updates(product).Error
}

func (r *ProductGormRepository) Delete(ctx context.Context, id uint, userID uint) error {
	return config.DB.WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).Delete(&domain.Product{}).Error
}
