package tests

import (
	"context"
	"errors"
	"testing"

	"vertice-backend/internal/domain"
	"vertice-backend/internal/service"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockOrderRepo struct {
	mock.Mock
}

func (m *MockOrderRepo) Create(ctx context.Context, order *domain.Order) error {
	args := m.Called(ctx, order)
	return args.Error(0)
}

func (m *MockOrderRepo) FindByIDAndUserID(ctx context.Context, id uint, userID uint) (*domain.Order, error) {
	args := m.Called(ctx, id, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Order), args.Error(1)
}

func (m *MockOrderRepo) FindByUserID(ctx context.Context, userID uint) ([]*domain.Order, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Order), args.Error(1)
}

func (m *MockOrderRepo) Update(ctx context.Context, order *domain.Order, userID uint) error {
	args := m.Called(ctx, order, userID)
	return args.Error(0)
}

func (m *MockOrderRepo) Delete(ctx context.Context, id uint, userID uint) error {
	args := m.Called(ctx, id, userID)
	return args.Error(0)
}

func TestCreateOrder_Success(t *testing.T) {
	mockOrderRepo := new(MockOrderRepo)
	mockProductRepo := new(MockProductRepo)
	orderService := service.NewOrderService(mockOrderRepo, mockProductRepo)

	product := &domain.Product{
		ID:     1,
		UserID: 1,
		Name:   "Test Product",
		Price:  10.0,
		Stock:  5,
	}
	mockProductRepo.On("FindByIDAndUserID", mock.Anything, uint(1), uint(1)).Return(product, nil)
	mockProductRepo.On("Update", mock.Anything, mock.AnythingOfType("*domain.Product"), uint(1)).Return(nil)
	mockOrderRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.Order")).Return(nil)

	expectedOrder := &domain.Order{
		ID:          1,
		UserID:      1,
		Status:      domain.OrderStatusPending,
		TotalAmount: 20.0,
		Items: []domain.OrderItem{
			{
				ID:        1,
				ProductID: 1,
				Quantity:  2,
				UnitPrice: 10.0,
				Subtotal:  20.0,
				Product:   *product,
			},
		},
	}
	mockOrderRepo.On("FindByIDAndUserID", mock.Anything, mock.AnythingOfType("uint"), mock.AnythingOfType("uint")).Return(expectedOrder, nil)

	req := service.CreateOrderRequest{
		Items: []service.OrderItemRequest{
			{ProductID: 1, Quantity: 2},
		},
	}

	order, err := orderService.CreateOrder(context.Background(), 1, req)

	assert.NoError(t, err)
	assert.Equal(t, uint(1), order.UserID)
	assert.Equal(t, domain.OrderStatusPending, order.Status)
	assert.Equal(t, 20.0, order.TotalAmount)
	assert.Len(t, order.Items, 1)
	assert.Equal(t, uint(1), order.Items[0].ProductID)
	assert.Equal(t, 2, order.Items[0].Quantity)
	assert.Equal(t, 10.0, order.Items[0].UnitPrice)
	assert.Equal(t, 20.0, order.Items[0].Subtotal)

	mockOrderRepo.AssertExpectations(t)
	mockProductRepo.AssertExpectations(t)
}

func TestCreateOrder_Error_EmptyItems(t *testing.T) {
	mockOrderRepo := new(MockOrderRepo)
	mockProductRepo := new(MockProductRepo)
	orderService := service.NewOrderService(mockOrderRepo, mockProductRepo)

	req := service.CreateOrderRequest{
		Items: []service.OrderItemRequest{},
	}

	_, err := orderService.CreateOrder(context.Background(), 1, req)

	assert.Error(t, err)
	assert.Equal(t, "order must have at least one item", err.Error())
}

func TestCreateOrder_Error_InvalidQuantity(t *testing.T) {
	mockOrderRepo := new(MockOrderRepo)
	mockProductRepo := new(MockProductRepo)
	orderService := service.NewOrderService(mockOrderRepo, mockProductRepo)

	req := service.CreateOrderRequest{
		Items: []service.OrderItemRequest{
			{ProductID: 1, Quantity: 0},
		},
	}

	_, err := orderService.CreateOrder(context.Background(), 1, req)

	assert.Error(t, err)
	assert.Equal(t, "quantity must be greater than 0", err.Error())
}

func TestCreateOrder_Error_ProductNotFound(t *testing.T) {
	mockOrderRepo := new(MockOrderRepo)
	mockProductRepo := new(MockProductRepo)
	orderService := service.NewOrderService(mockOrderRepo, mockProductRepo)

	mockProductRepo.On("FindByIDAndUserID", mock.Anything, uint(1), uint(1)).Return(nil, errors.New("not found"))

	req := service.CreateOrderRequest{
		Items: []service.OrderItemRequest{
			{ProductID: 1, Quantity: 2},
		},
	}

	_, err := orderService.CreateOrder(context.Background(), 1, req)

	assert.Error(t, err)
	assert.Equal(t, "product not found", err.Error())
	mockProductRepo.AssertExpectations(t)
}

func TestCreateOrder_Error_InsufficientStock(t *testing.T) {
	mockOrderRepo := new(MockOrderRepo)
	mockProductRepo := new(MockProductRepo)
	orderService := service.NewOrderService(mockOrderRepo, mockProductRepo)

	product := &domain.Product{
		ID:     1,
		UserID: 1,
		Name:   "Test Product",
		Price:  10.0,
		Stock:  1,
	}
	mockProductRepo.On("FindByIDAndUserID", mock.Anything, uint(1), uint(1)).Return(product, nil)

	req := service.CreateOrderRequest{
		Items: []service.OrderItemRequest{
			{ProductID: 1, Quantity: 5},
		},
	}

	_, err := orderService.CreateOrder(context.Background(), 1, req)

	assert.Error(t, err)
	assert.Equal(t, "insufficient stock for product: Test Product", err.Error())
	mockProductRepo.AssertExpectations(t)
}

func TestGetOrder_Success(t *testing.T) {
	mockOrderRepo := new(MockOrderRepo)
	mockProductRepo := new(MockProductRepo)
	orderService := service.NewOrderService(mockOrderRepo, mockProductRepo)

	expectedOrder := &domain.Order{
		ID:     1,
		UserID: 1,
		Status: domain.OrderStatusPending,
	}
	mockOrderRepo.On("FindByIDAndUserID", mock.Anything, uint(1), uint(1)).Return(expectedOrder, nil)

	order, err := orderService.GetOrder(context.Background(), 1, 1)

	assert.NoError(t, err)
	assert.Equal(t, expectedOrder, order)
	mockOrderRepo.AssertExpectations(t)
}

func TestGetOrder_Error_NotFound(t *testing.T) {
	mockOrderRepo := new(MockOrderRepo)
	mockProductRepo := new(MockProductRepo)
	orderService := service.NewOrderService(mockOrderRepo, mockProductRepo)

	mockOrderRepo.On("FindByIDAndUserID", mock.Anything, uint(1), uint(1)).Return(nil, errors.New("not found"))

	_, err := orderService.GetOrder(context.Background(), 1, 1)

	assert.Error(t, err)
	mockOrderRepo.AssertExpectations(t)
}

func TestGetOrdersByUser_Success(t *testing.T) {
	mockOrderRepo := new(MockOrderRepo)
	mockProductRepo := new(MockProductRepo)
	orderService := service.NewOrderService(mockOrderRepo, mockProductRepo)

	expectedOrders := []*domain.Order{
		{ID: 1, UserID: 1, Status: domain.OrderStatusPending},
		{ID: 2, UserID: 1, Status: domain.OrderStatusConfirmed},
	}
	mockOrderRepo.On("FindByUserID", mock.Anything, uint(1)).Return(expectedOrders, nil)

	orders, err := orderService.GetOrdersByUser(context.Background(), 1)

	assert.NoError(t, err)
	assert.Equal(t, expectedOrders, orders)
	mockOrderRepo.AssertExpectations(t)
}

func TestUpdateOrderStatus_Success(t *testing.T) {
	mockOrderRepo := new(MockOrderRepo)
	mockProductRepo := new(MockProductRepo)
	orderService := service.NewOrderService(mockOrderRepo, mockProductRepo)

	existingOrder := &domain.Order{
		ID:     1,
		UserID: 1,
		Status: domain.OrderStatusPending,
	}
	mockOrderRepo.On("FindByIDAndUserID", mock.Anything, uint(1), uint(1)).Return(existingOrder, nil)
	mockOrderRepo.On("Update", mock.Anything, mock.AnythingOfType("*domain.Order"), uint(1)).Return(nil)

	order, err := orderService.UpdateOrderStatus(context.Background(), 1, 1, domain.OrderStatusConfirmed)

	assert.NoError(t, err)
	assert.Equal(t, domain.OrderStatusConfirmed, order.Status)
	mockOrderRepo.AssertExpectations(t)
}

func TestUpdateOrderStatus_Error_InvalidTransition(t *testing.T) {
	mockOrderRepo := new(MockOrderRepo)
	mockProductRepo := new(MockProductRepo)
	orderService := service.NewOrderService(mockOrderRepo, mockProductRepo)

	existingOrder := &domain.Order{
		ID:     1,
		UserID: 1,
		Status: domain.OrderStatusDelivered,
	}
	mockOrderRepo.On("FindByIDAndUserID", mock.Anything, uint(1), uint(1)).Return(existingOrder, nil)

	_, err := orderService.UpdateOrderStatus(context.Background(), 1, 1, domain.OrderStatusPending)

	assert.Error(t, err)
	assert.Equal(t, "invalid status transition", err.Error())
	mockOrderRepo.AssertExpectations(t)
}

func TestCancelOrder_Success(t *testing.T) {
	mockOrderRepo := new(MockOrderRepo)
	mockProductRepo := new(MockProductRepo)
	orderService := service.NewOrderService(mockOrderRepo, mockProductRepo)

	existingOrder := &domain.Order{
		ID:     1,
		UserID: 1,
		Status: domain.OrderStatusPending,
		Items: []domain.OrderItem{
			{ProductID: 1, Quantity: 2},
		},
	}
	mockOrderRepo.On("FindByIDAndUserID", mock.Anything, uint(1), uint(1)).Return(existingOrder, nil)

	product := &domain.Product{ID: 1, UserID: 1, Stock: 3}
	mockProductRepo.On("FindByIDAndUserID", mock.Anything, uint(1), uint(1)).Return(product, nil)
	mockProductRepo.On("Update", mock.Anything, mock.AnythingOfType("*domain.Product"), uint(1)).Return(nil)

	mockOrderRepo.On("Update", mock.Anything, mock.AnythingOfType("*domain.Order"), uint(1)).Return(nil)

	order, err := orderService.CancelOrder(context.Background(), 1, 1)

	assert.NoError(t, err)
	assert.Equal(t, domain.OrderStatusCancelled, order.Status)
	mockOrderRepo.AssertExpectations(t)
	mockProductRepo.AssertExpectations(t)
}

func TestCancelOrder_Error_AlreadyCancelled(t *testing.T) {
	mockOrderRepo := new(MockOrderRepo)
	mockProductRepo := new(MockProductRepo)
	orderService := service.NewOrderService(mockOrderRepo, mockProductRepo)

	existingOrder := &domain.Order{
		ID:     1,
		UserID: 1,
		Status: domain.OrderStatusCancelled,
	}
	mockOrderRepo.On("FindByIDAndUserID", mock.Anything, uint(1), uint(1)).Return(existingOrder, nil)

	_, err := orderService.CancelOrder(context.Background(), 1, 1)

	assert.Error(t, err)
	assert.Equal(t, "order is already cancelled", err.Error())
	mockOrderRepo.AssertExpectations(t)
}

func TestCancelOrder_Error_DeliveredOrder(t *testing.T) {
	mockOrderRepo := new(MockOrderRepo)
	mockProductRepo := new(MockProductRepo)
	orderService := service.NewOrderService(mockOrderRepo, mockProductRepo)

	existingOrder := &domain.Order{
		ID:     1,
		UserID: 1,
		Status: domain.OrderStatusDelivered,
	}
	mockOrderRepo.On("FindByIDAndUserID", mock.Anything, uint(1), uint(1)).Return(existingOrder, nil)

	_, err := orderService.CancelOrder(context.Background(), 1, 1)

	assert.Error(t, err)
	assert.Equal(t, "cannot cancel delivered order", err.Error())
	mockOrderRepo.AssertExpectations(t)
}

func TestDeleteOrder_Success(t *testing.T) {
	mockOrderRepo := new(MockOrderRepo)
	mockProductRepo := new(MockProductRepo)
	orderService := service.NewOrderService(mockOrderRepo, mockProductRepo)

	existingOrder := &domain.Order{
		ID:     1,
		UserID: 1,
		Status: domain.OrderStatusCancelled,
	}
	mockOrderRepo.On("FindByIDAndUserID", mock.Anything, uint(1), uint(1)).Return(existingOrder, nil)
	mockOrderRepo.On("Delete", mock.Anything, uint(1), uint(1)).Return(nil)

	err := orderService.DeleteOrder(context.Background(), 1, 1)

	assert.NoError(t, err)
	mockOrderRepo.AssertExpectations(t)
}

func TestDeleteOrder_Error_NotCancelled(t *testing.T) {
	mockOrderRepo := new(MockOrderRepo)
	mockProductRepo := new(MockProductRepo)
	orderService := service.NewOrderService(mockOrderRepo, mockProductRepo)

	existingOrder := &domain.Order{
		ID:     1,
		UserID: 1,
		Status: domain.OrderStatusPending,
	}
	mockOrderRepo.On("FindByIDAndUserID", mock.Anything, uint(1), uint(1)).Return(existingOrder, nil)

	err := orderService.DeleteOrder(context.Background(), 1, 1)

	assert.Error(t, err)
	assert.Equal(t, "can only delete cancelled orders", err.Error())
	mockOrderRepo.AssertExpectations(t)
}

func TestCreateOrder_MultipleItems_Success(t *testing.T) {
	mockOrderRepo := new(MockOrderRepo)
	mockProductRepo := new(MockProductRepo)
	orderService := service.NewOrderService(mockOrderRepo, mockProductRepo)

	product1 := &domain.Product{ID: 1, UserID: 1, Name: "Prod1", Price: 10.0, Stock: 10}
	product2 := &domain.Product{ID: 2, UserID: 1, Name: "Prod2", Price: 20.0, Stock: 10}
	mockProductRepo.On("FindByIDAndUserID", mock.Anything, uint(1), uint(1)).Return(product1, nil)
	mockProductRepo.On("FindByIDAndUserID", mock.Anything, uint(2), uint(1)).Return(product2, nil)
	mockProductRepo.On("Update", mock.Anything, mock.AnythingOfType("*domain.Product"), uint(1)).Return(nil)
	mockOrderRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.Order")).Return(nil)

	expectedOrder := &domain.Order{
		ID:          1,
		UserID:      1,
		Status:      domain.OrderStatusPending,
		TotalAmount: 50.0,
		Items: []domain.OrderItem{
			{
				ID:        1,
				ProductID: 1,
				Quantity:  2,
				UnitPrice: 10.0,
				Subtotal:  20.0,
				Product:   *product1,
			},
			{
				ID:        2,
				ProductID: 2,
				Quantity:  1,
				UnitPrice: 20.0,
				Subtotal:  20.0,
				Product:   *product2,
			},
		},
	}
	mockOrderRepo.On("FindByIDAndUserID", mock.Anything, mock.AnythingOfType("uint"), mock.AnythingOfType("uint")).Return(expectedOrder, nil)

	req := service.CreateOrderRequest{
		Items: []service.OrderItemRequest{
			{ProductID: 1, Quantity: 2},
			{ProductID: 2, Quantity: 1},
		},
	}

	order, err := orderService.CreateOrder(context.Background(), 1, req)

	assert.NoError(t, err)
	assert.Equal(t, uint(1), order.UserID)
	assert.Equal(t, domain.OrderStatusPending, order.Status)
	assert.Equal(t, 50.0, order.TotalAmount)
	assert.Len(t, order.Items, 2)
	assert.Equal(t, uint(1), order.Items[0].ProductID)
	assert.Equal(t, 2, order.Items[0].Quantity)
	assert.Equal(t, 10.0, order.Items[0].UnitPrice)
	assert.Equal(t, 20.0, order.Items[0].Subtotal)
	assert.Equal(t, uint(2), order.Items[1].ProductID)
	assert.Equal(t, 1, order.Items[1].Quantity)
	assert.Equal(t, 20.0, order.Items[1].UnitPrice)
	assert.Equal(t, 20.0, order.Items[1].Subtotal)

	mockOrderRepo.AssertExpectations(t)
	mockProductRepo.AssertExpectations(t)
}
