package service

import (
	"context"
	"errors"
	"slices"
	"vertice-backend/internal/domain"
)

type OrderService struct {
	orderRepo   domain.OrderRepository
	productRepo domain.ProductRepository
}

func NewOrderService(orderRepo domain.OrderRepository, productRepo domain.ProductRepository) *OrderService {
	return &OrderService{
		orderRepo:   orderRepo,
		productRepo: productRepo,
	}
}

type CreateOrderRequest struct {
	Items []OrderItemRequest `json:"items"`
}

type OrderItemRequest struct {
	ProductID uint `json:"product_id"`
	Quantity  int  `json:"quantity"`
}

func (s *OrderService) CreateOrder(ctx context.Context, userID uint, req CreateOrderRequest) (*domain.Order, error) {
	if len(req.Items) == 0 {
		return nil, errors.New("order must have at least one item")
	}

	order := &domain.Order{
		UserID:      userID,
		Status:      domain.OrderStatusPending,
		TotalAmount: 0,
		Items:       []domain.OrderItem{},
	}

	for _, itemReq := range req.Items {
		if itemReq.Quantity <= 0 {
			return nil, errors.New("quantity must be greater than 0")
		}

		product, err := s.productRepo.FindByIDAndUserID(ctx, itemReq.ProductID, userID)
		if err != nil {
			return nil, errors.New("product not found")
		}

		if product.Stock < itemReq.Quantity {
			return nil, errors.New("insufficient stock for product: " + product.Name)
		}

		subtotal := float64(itemReq.Quantity) * product.Price

		orderItem := domain.OrderItem{
			ProductID: product.ID,
			Quantity:  itemReq.Quantity,
			UnitPrice: product.Price,
			Subtotal:  subtotal,
		}

		order.Items = append(order.Items, orderItem)
		order.TotalAmount += subtotal

		product.Stock -= itemReq.Quantity
		if err := s.productRepo.Update(ctx, product, userID); err != nil {
			return nil, err
		}
	}

	if err := s.orderRepo.Create(ctx, order); err != nil {
		return nil, err
	}

	return s.orderRepo.FindByIDAndUserID(ctx, order.ID, userID)
}

func (s *OrderService) GetOrder(ctx context.Context, id, userID uint) (*domain.Order, error) {
	return s.orderRepo.FindByIDAndUserID(ctx, id, userID)
}

func (s *OrderService) GetOrdersByUser(ctx context.Context, userID uint) ([]*domain.Order, error) {
	return s.orderRepo.FindByUserID(ctx, userID)
}

func (s *OrderService) UpdateOrderStatus(ctx context.Context, id, userID uint, status domain.OrderStatus) (*domain.Order, error) {
	order, err := s.orderRepo.FindByIDAndUserID(ctx, id, userID)
	if err != nil {
		return nil, errors.New("order not found")
	}

	if !isValidStatusTransition(order.Status, status) {
		return nil, errors.New("invalid status transition")
	}

	order.Status = status
	if err := s.orderRepo.Update(ctx, order, userID); err != nil {
		return nil, err
	}

	return order, nil
}

func (s *OrderService) CancelOrder(ctx context.Context, id, userID uint) (*domain.Order, error) {
	order, err := s.orderRepo.FindByIDAndUserID(ctx, id, userID)
	if err != nil {
		return nil, errors.New("order not found")
	}

	if order.Status == domain.OrderStatusCancelled {
		return nil, errors.New("order is already cancelled")
	}

	if order.Status == domain.OrderStatusDelivered {
		return nil, errors.New("cannot cancel delivered order")
	}

	if order.Status == domain.OrderStatusPending || order.Status == domain.OrderStatusConfirmed {
		for _, item := range order.Items {
			product, err := s.productRepo.FindByIDAndUserID(ctx, item.ProductID, userID)
			if err != nil {
				continue
			}
			product.Stock += item.Quantity
			s.productRepo.Update(ctx, product, userID)
		}
	}

	order.Status = domain.OrderStatusCancelled
	if err := s.orderRepo.Update(ctx, order, userID); err != nil {
		return nil, err
	}

	return order, nil
}

func (s *OrderService) DeleteOrder(ctx context.Context, id, userID uint) error {
	order, err := s.orderRepo.FindByIDAndUserID(ctx, id, userID)
	if err != nil {
		return errors.New("order not found")
	}

	if order.Status != domain.OrderStatusCancelled {
		return errors.New("can only delete cancelled orders")
	}

	return s.orderRepo.Delete(ctx, id, userID)
}

func isValidStatusTransition(current, new domain.OrderStatus) bool {
	validTransitions := map[domain.OrderStatus][]domain.OrderStatus{
		domain.OrderStatusPending: {
			domain.OrderStatusConfirmed,
			domain.OrderStatusCancelled,
		},
		domain.OrderStatusConfirmed: {
			domain.OrderStatusShipped,
			domain.OrderStatusCancelled,
		},
		domain.OrderStatusShipped: {
			domain.OrderStatusDelivered,
		},
		domain.OrderStatusDelivered: {},
		domain.OrderStatusCancelled: {},
	}

	allowedTransitions, exists := validTransitions[current]
	if !exists {
		return false
	}

	return slices.Contains(allowedTransitions, new)
}
