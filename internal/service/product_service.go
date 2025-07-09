package service

import (
	"context"
	"errors"
	"vertice-backend/internal/domain"
)

type ProductService struct {
	repo domain.ProductRepository
}

func NewProductService(repo domain.ProductRepository) *ProductService {
	return &ProductService{repo: repo}
}

func (s *ProductService) CreateProduct(ctx context.Context, userID uint, code, name, description string, price float64, stock int) (*domain.Product, error) {
	if code == "" || name == "" {
		return nil, errors.New("code and name are required")
	}
	if price < 0 {
		return nil, errors.New("price cannot be negative")
	}
	if stock < 0 {
		return nil, errors.New("stock cannot be negative")
	}

	existingProduct, err := s.repo.FindByCodeAndUserID(ctx, code, userID)
	if err == nil && existingProduct != nil {
		return nil, errors.New("product code already exists for this user")
	}

	product := &domain.Product{
		UserID:      userID,
		Code:        code,
		Name:        name,
		Description: description,
		Price:       price,
		Stock:       stock,
	}

	if err := s.repo.Create(ctx, product); err != nil {
		return nil, err
	}

	return product, nil
}

func (s *ProductService) GetProduct(ctx context.Context, id, userID uint) (*domain.Product, error) {
	return s.repo.FindByIDAndUserID(ctx, id, userID)
}

func (s *ProductService) GetProductsByUser(ctx context.Context, userID uint) ([]*domain.Product, error) {
	return s.repo.FindByUserID(ctx, userID)
}

func (s *ProductService) GetProductByCode(ctx context.Context, code string, userID uint) (*domain.Product, error) {
	return s.repo.FindByCodeAndUserID(ctx, code, userID)
}

func (s *ProductService) UpdateProduct(ctx context.Context, id, userID uint, code, name, description *string, price *float64, stock *int) (*domain.Product, error) {
	existingProduct, err := s.repo.FindByIDAndUserID(ctx, id, userID)
	if err != nil {
		return nil, errors.New("product not found")
	}

	// Update only the fields that are provided (not nil)
	if code != nil {
		if *code == "" {
			return nil, errors.New("code cannot be empty")
		}
		if *code != existingProduct.Code {
			conflictingProduct, err := s.repo.FindByCodeAndUserID(ctx, *code, userID)
			if err == nil && conflictingProduct != nil {
				return nil, errors.New("product code already exists for this user")
			}
		}
		existingProduct.Code = *code
	}

	if name != nil {
		if *name == "" {
			return nil, errors.New("name cannot be empty")
		}
		existingProduct.Name = *name
	}

	if description != nil {
		existingProduct.Description = *description
	}

	if price != nil {
		if *price < 0 {
			return nil, errors.New("price cannot be negative")
		}
		existingProduct.Price = *price
	}

	if stock != nil {
		if *stock < 0 {
			return nil, errors.New("stock cannot be negative")
		}
		existingProduct.Stock = *stock
	}

	if err := s.repo.Update(ctx, existingProduct, userID); err != nil {
		return nil, err
	}

	return existingProduct, nil
}

func (s *ProductService) UpdateProductStock(ctx context.Context, id, userID uint, stockDelta int) (*domain.Product, error) {
	product, err := s.repo.FindByIDAndUserID(ctx, id, userID)
	if err != nil {
		return nil, errors.New("product not found")
	}
	if product.Stock+stockDelta < 0 {
		return nil, errors.New("stock cannot be negative")
	}
	product.Stock += stockDelta
	if err := s.repo.Update(ctx, product, userID); err != nil {
		return nil, err
	}
	return product, nil
}

func (s *ProductService) DeleteProduct(ctx context.Context, id, userID uint) error {
	_, err := s.repo.FindByIDAndUserID(ctx, id, userID)
	if err != nil {
		return errors.New("product not found")
	}

	return s.repo.Delete(ctx, id, userID)
}
