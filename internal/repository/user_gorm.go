package repository

import (
	"context"
	"vertice-backend/config"
	"vertice-backend/internal/domain"
)

type UserGormRepository struct{}

func NewUserGormRepository() *UserGormRepository {
	return &UserGormRepository{}
}

func (r *UserGormRepository) Create(ctx context.Context, user *domain.User) error {
	return config.DB.WithContext(ctx).Create(user).Error
}

func (r *UserGormRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User
	err := config.DB.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserGormRepository) FindByID(ctx context.Context, id uint) (*domain.User, error) {
	var user domain.User
	err := config.DB.WithContext(ctx).First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
