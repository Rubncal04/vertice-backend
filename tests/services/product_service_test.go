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

// Mock of ProductRepository
type MockProductRepo struct {
	mock.Mock
}

func (m *MockProductRepo) Create(ctx context.Context, product *domain.Product) error {
	args := m.Called(ctx, product)
	return args.Error(0)
}

func (m *MockProductRepo) FindByIDAndUserID(ctx context.Context, id uint, userID uint) (*domain.Product, error) {
	args := m.Called(ctx, id, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Product), args.Error(1)
}

func (m *MockProductRepo) FindByUserID(ctx context.Context, userID uint) ([]*domain.Product, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Product), args.Error(1)
}

func (m *MockProductRepo) FindByCodeAndUserID(ctx context.Context, code string, userID uint) (*domain.Product, error) {
	args := m.Called(ctx, code, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Product), args.Error(1)
}

func (m *MockProductRepo) Update(ctx context.Context, product *domain.Product, userID uint) error {
	args := m.Called(ctx, product, userID)
	return args.Error(0)
}

func (m *MockProductRepo) Delete(ctx context.Context, id uint, userID uint) error {
	args := m.Called(ctx, id, userID)
	return args.Error(0)
}

func TestCreateProduct_Success(t *testing.T) {
	mockRepo := new(MockProductRepo)
	service := service.NewProductService(mockRepo)

	// Mock FindByCodeAndUserID to return nil (no existing product)
	mockRepo.On("FindByCodeAndUserID", mock.Anything, "PROD001", uint(1)).Return(nil, errors.New("not found"))
	// Mock Create to succeed
	mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.Product")).Return(nil)

	product, err := service.CreateProduct(context.Background(), 1, "PROD001", "Test Product", "Test Description", 99.99, 10)

	assert.NoError(t, err)
	assert.Equal(t, uint(1), product.UserID)
	assert.Equal(t, "PROD001", product.Code)
	assert.Equal(t, "Test Product", product.Name)
	assert.Equal(t, "Test Description", product.Description)
	assert.Equal(t, 99.99, product.Price)
	assert.Equal(t, 10, product.Stock)
	mockRepo.AssertExpectations(t)
}

func TestCreateProduct_ValidationError_EmptyCode(t *testing.T) {
	mockRepo := new(MockProductRepo)
	service := service.NewProductService(mockRepo)

	_, err := service.CreateProduct(context.Background(), 1, "", "Test Product", "Test Description", 99.99, 10)

	assert.Error(t, err)
	assert.Equal(t, "code and name are required", err.Error())
}

func TestCreateProduct_ValidationError_EmptyName(t *testing.T) {
	mockRepo := new(MockProductRepo)
	service := service.NewProductService(mockRepo)

	_, err := service.CreateProduct(context.Background(), 1, "PROD001", "", "Test Description", 99.99, 10)

	assert.Error(t, err)
	assert.Equal(t, "code and name are required", err.Error())
}

func TestCreateProduct_ValidationError_NegativePrice(t *testing.T) {
	mockRepo := new(MockProductRepo)
	service := service.NewProductService(mockRepo)

	_, err := service.CreateProduct(context.Background(), 1, "PROD001", "Test Product", "Test Description", -10.0, 10)

	assert.Error(t, err)
	assert.Equal(t, "price cannot be negative", err.Error())
}

func TestCreateProduct_ValidationError_NegativeStock(t *testing.T) {
	mockRepo := new(MockProductRepo)
	service := service.NewProductService(mockRepo)

	_, err := service.CreateProduct(context.Background(), 1, "PROD001", "Test Product", "Test Description", 99.99, -5)

	assert.Error(t, err)
	assert.Equal(t, "stock cannot be negative", err.Error())
}

func TestCreateProduct_Error_CodeAlreadyExists(t *testing.T) {
	mockRepo := new(MockProductRepo)
	service := service.NewProductService(mockRepo)

	existingProduct := &domain.Product{ID: 1, UserID: 1, Code: "PROD001", Name: "Existing Product"}
	mockRepo.On("FindByCodeAndUserID", mock.Anything, "PROD001", uint(1)).Return(existingProduct, nil).Maybe()

	_, err := service.CreateProduct(context.Background(), 1, "PROD001", "Test Product", "Test Description", 99.99, 10)

	assert.Error(t, err)
	assert.Equal(t, "product code already exists for this user", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestGetProduct_Success(t *testing.T) {
	mockRepo := new(MockProductRepo)
	service := service.NewProductService(mockRepo)

	expectedProduct := &domain.Product{ID: 1, UserID: 1, Code: "PROD001", Name: "Test Product"}
	mockRepo.On("FindByIDAndUserID", mock.Anything, uint(1), uint(1)).Return(expectedProduct, nil)

	product, err := service.GetProduct(context.Background(), 1, 1)

	assert.NoError(t, err)
	assert.Equal(t, expectedProduct, product)
	mockRepo.AssertExpectations(t)
}

func TestGetProduct_Error_NotFound(t *testing.T) {
	mockRepo := new(MockProductRepo)
	service := service.NewProductService(mockRepo)

	mockRepo.On("FindByIDAndUserID", mock.Anything, uint(1), uint(1)).Return(nil, errors.New("not found"))

	_, err := service.GetProduct(context.Background(), 1, 1)

	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}

func TestGetProductsByUser_Success(t *testing.T) {
	mockRepo := new(MockProductRepo)
	service := service.NewProductService(mockRepo)

	expectedProducts := []*domain.Product{
		{ID: 1, UserID: 1, Code: "PROD001", Name: "Product 1"},
		{ID: 2, UserID: 1, Code: "PROD002", Name: "Product 2"},
	}
	mockRepo.On("FindByUserID", mock.Anything, uint(1)).Return(expectedProducts, nil)

	products, err := service.GetProductsByUser(context.Background(), 1)

	assert.NoError(t, err)
	assert.Equal(t, expectedProducts, products)
	mockRepo.AssertExpectations(t)
}

func TestGetProductByCode_Success(t *testing.T) {
	mockRepo := new(MockProductRepo)
	service := service.NewProductService(mockRepo)

	expectedProduct := &domain.Product{ID: 1, UserID: 1, Code: "PROD001", Name: "Test Product"}
	mockRepo.On("FindByCodeAndUserID", mock.Anything, "PROD001", uint(1)).Return(expectedProduct, nil)

	product, err := service.GetProductByCode(context.Background(), "PROD001", 1)

	assert.NoError(t, err)
	assert.Equal(t, expectedProduct, product)
	mockRepo.AssertExpectations(t)
}

func TestUpdateProduct_Success(t *testing.T) {
	mockRepo := new(MockProductRepo)
	service := service.NewProductService(mockRepo)

	existingProduct := &domain.Product{ID: 1, UserID: 1, Code: "PROD001", Name: "Old Name"}
	mockRepo.On("FindByIDAndUserID", mock.Anything, uint(1), uint(1)).Return(existingProduct, nil)
	mockRepo.On("FindByCodeAndUserID", mock.Anything, "PROD001", uint(1)).Return(existingProduct, nil).Maybe()
	mockRepo.On("Update", mock.Anything, mock.AnythingOfType("*domain.Product"), uint(1)).Return(nil)

	code := "PROD001"
	name := "New Name"
	description := "New Description"
	price := 150.0
	stock := 20
	product, err := service.UpdateProduct(context.Background(), 1, 1, &code, &name, &description, &price, &stock)

	assert.NoError(t, err)
	assert.Equal(t, "New Name", product.Name)
	assert.Equal(t, "New Description", product.Description)
	assert.Equal(t, 150.0, product.Price)
	assert.Equal(t, 20, product.Stock)
	mockRepo.AssertExpectations(t)
}

func TestUpdateProduct_Error_ProductNotFound(t *testing.T) {
	mockRepo := new(MockProductRepo)
	service := service.NewProductService(mockRepo)

	mockRepo.On("FindByIDAndUserID", mock.Anything, uint(1), uint(1)).Return(nil, errors.New("not found"))

	code := "PROD001"
	name := "New Name"
	description := "New Description"
	price := 150.0
	stock := 20
	_, err := service.UpdateProduct(context.Background(), 1, 1, &code, &name, &description, &price, &stock)

	assert.Error(t, err)
	assert.Equal(t, "product not found", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestUpdateProduct_Error_CodeConflict(t *testing.T) {
	mockRepo := new(MockProductRepo)
	service := service.NewProductService(mockRepo)

	existingProduct := &domain.Product{ID: 1, UserID: 1, Code: "PROD001", Name: "Old Name"}
	conflictingProduct := &domain.Product{ID: 2, UserID: 1, Code: "PROD002", Name: "Other Product"}

	mockRepo.On("FindByIDAndUserID", mock.Anything, uint(1), uint(1)).Return(existingProduct, nil)
	mockRepo.On("FindByCodeAndUserID", mock.Anything, "PROD002", uint(1)).Return(conflictingProduct, nil)

	code := "PROD002"
	name := "New Name"
	description := "New Description"
	price := 150.0
	stock := 20
	_, err := service.UpdateProduct(context.Background(), 1, 1, &code, &name, &description, &price, &stock)

	assert.Error(t, err)
	assert.Equal(t, "product code already exists for this user", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestDeleteProduct_Success(t *testing.T) {
	mockRepo := new(MockProductRepo)
	service := service.NewProductService(mockRepo)

	existingProduct := &domain.Product{ID: 1, UserID: 1, Code: "PROD001", Name: "Test Product"}
	mockRepo.On("FindByIDAndUserID", mock.Anything, uint(1), uint(1)).Return(existingProduct, nil)
	mockRepo.On("Delete", mock.Anything, uint(1), uint(1)).Return(nil)

	err := service.DeleteProduct(context.Background(), 1, 1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestDeleteProduct_Error_ProductNotFound(t *testing.T) {
	mockRepo := new(MockProductRepo)
	service := service.NewProductService(mockRepo)

	mockRepo.On("FindByIDAndUserID", mock.Anything, uint(1), uint(1)).Return(nil, errors.New("not found"))

	err := service.DeleteProduct(context.Background(), 1, 1)

	assert.Error(t, err)
	assert.Equal(t, "product not found", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestUpdateProductStock_Success_Add(t *testing.T) {
	mockRepo := new(MockProductRepo)
	service := service.NewProductService(mockRepo)

	mockRepo.On("FindByIDAndUserID", mock.Anything, uint(1), uint(1)).Return(&domain.Product{ID: 1, UserID: 1, Stock: 10}, nil)
	mockRepo.On("Update", mock.Anything, mock.AnythingOfType("*domain.Product"), uint(1)).Return(nil)

	updated, err := service.UpdateProductStock(context.Background(), 1, 1, 5)
	assert.NoError(t, err)
	assert.Equal(t, uint(1), updated.ID)
	assert.Equal(t, 15, updated.Stock)
	mockRepo.AssertExpectations(t)
}

func TestUpdateProductStock_Success_Subtract(t *testing.T) {
	mockRepo := new(MockProductRepo)
	service := service.NewProductService(mockRepo)

	mockRepo.On("FindByIDAndUserID", mock.Anything, uint(1), uint(1)).Return(&domain.Product{ID: 1, UserID: 1, Stock: 10}, nil)
	mockRepo.On("Update", mock.Anything, mock.AnythingOfType("*domain.Product"), uint(1)).Return(nil)

	updated, err := service.UpdateProductStock(context.Background(), 1, 1, -3)
	assert.NoError(t, err)
	assert.Equal(t, uint(1), updated.ID)
	assert.Equal(t, 7, updated.Stock)
	mockRepo.AssertExpectations(t)
}

func TestUpdateProductStock_Error_NegativeStock(t *testing.T) {
	mockRepo := new(MockProductRepo)
	service := service.NewProductService(mockRepo)

	mockRepo.On("FindByIDAndUserID", mock.Anything, uint(1), uint(1)).Return(&domain.Product{ID: 1, UserID: 1, Stock: 2}, nil)

	updated, err := service.UpdateProductStock(context.Background(), 1, 1, -5)
	assert.Error(t, err)
	assert.Nil(t, updated)
	assert.Equal(t, "stock cannot be negative", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestUpdateProductStock_Error_ProductNotFound(t *testing.T) {
	mockRepo := new(MockProductRepo)
	service := service.NewProductService(mockRepo)

	mockRepo.On("FindByIDAndUserID", mock.Anything, uint(1), uint(1)).Return(nil, errors.New("not found"))

	updated, err := service.UpdateProductStock(context.Background(), 1, 1, 5)
	assert.Error(t, err)
	assert.Nil(t, updated)
	assert.Equal(t, "product not found", err.Error())
	mockRepo.AssertExpectations(t)
}

// Tests for UpdateProduct
func TestUpdateProduct_Success_UpdateNameOnly(t *testing.T) {
	mockRepo := new(MockProductRepo)
	service := service.NewProductService(mockRepo)

	existingProduct := &domain.Product{ID: 1, UserID: 1, Code: "PROD001", Name: "Old Name", Description: "Old Description", Price: 100.0, Stock: 10}
	mockRepo.On("FindByIDAndUserID", mock.Anything, uint(1), uint(1)).Return(existingProduct, nil)
	mockRepo.On("Update", mock.Anything, mock.AnythingOfType("*domain.Product"), uint(1)).Return(nil)

	newName := "New Name"
	product, err := service.UpdateProduct(context.Background(), 1, 1, nil, &newName, nil, nil, nil)

	assert.NoError(t, err)
	assert.Equal(t, "New Name", product.Name)
	assert.Equal(t, "Old Description", product.Description) // Should remain unchanged
	assert.Equal(t, 100.0, product.Price)                   // Should remain unchanged
	assert.Equal(t, 10, product.Stock)                      // Should remain unchanged
	mockRepo.AssertExpectations(t)
}

func TestUpdateProduct_Success_UpdatePriceAndStockOnly(t *testing.T) {
	mockRepo := new(MockProductRepo)
	service := service.NewProductService(mockRepo)

	existingProduct := &domain.Product{ID: 1, UserID: 1, Code: "PROD001", Name: "Test Product", Description: "Test Description", Price: 100.0, Stock: 10}
	mockRepo.On("FindByIDAndUserID", mock.Anything, uint(1), uint(1)).Return(existingProduct, nil)
	mockRepo.On("Update", mock.Anything, mock.AnythingOfType("*domain.Product"), uint(1)).Return(nil)

	newPrice := 150.0
	newStock := 25
	product, err := service.UpdateProduct(context.Background(), 1, 1, nil, nil, nil, &newPrice, &newStock)

	assert.NoError(t, err)
	assert.Equal(t, "Test Product", product.Name)            // Should remain unchanged
	assert.Equal(t, "Test Description", product.Description) // Should remain unchanged
	assert.Equal(t, 150.0, product.Price)                    // Should be updated
	assert.Equal(t, 25, product.Stock)                       // Should be updated
	mockRepo.AssertExpectations(t)
}

func TestUpdateProduct_Success_UpdateCodeOnly(t *testing.T) {
	mockRepo := new(MockProductRepo)
	service := service.NewProductService(mockRepo)

	existingProduct := &domain.Product{ID: 1, UserID: 1, Code: "PROD001", Name: "Test Product", Description: "Test Description", Price: 100.0, Stock: 10}
	mockRepo.On("FindByIDAndUserID", mock.Anything, uint(1), uint(1)).Return(existingProduct, nil)
	mockRepo.On("FindByCodeAndUserID", mock.Anything, "PROD002", uint(1)).Return(nil, errors.New("not found"))
	mockRepo.On("Update", mock.Anything, mock.AnythingOfType("*domain.Product"), uint(1)).Return(nil)

	newCode := "PROD002"
	product, err := service.UpdateProduct(context.Background(), 1, 1, &newCode, nil, nil, nil, nil)

	assert.NoError(t, err)
	assert.Equal(t, "PROD002", product.Code)                 // Should be updated
	assert.Equal(t, "Test Product", product.Name)            // Should remain unchanged
	assert.Equal(t, "Test Description", product.Description) // Should remain unchanged
	assert.Equal(t, 100.0, product.Price)                    // Should remain unchanged
	assert.Equal(t, 10, product.Stock)                       // Should remain unchanged
	mockRepo.AssertExpectations(t)
}

func TestUpdateProductPartial_Error_ProductNotFound(t *testing.T) {
	mockRepo := new(MockProductRepo)
	service := service.NewProductService(mockRepo)

	mockRepo.On("FindByIDAndUserID", mock.Anything, uint(1), uint(1)).Return(nil, errors.New("not found"))

	newName := "New Name"
	_, err := service.UpdateProduct(context.Background(), 1, 1, nil, &newName, nil, nil, nil)

	assert.Error(t, err)
	assert.Equal(t, "product not found", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestUpdateProduct_Error_EmptyCode(t *testing.T) {
	mockRepo := new(MockProductRepo)
	service := service.NewProductService(mockRepo)

	existingProduct := &domain.Product{ID: 1, UserID: 1, Code: "PROD001", Name: "Test Product"}
	mockRepo.On("FindByIDAndUserID", mock.Anything, uint(1), uint(1)).Return(existingProduct, nil)

	emptyCode := ""
	_, err := service.UpdateProduct(context.Background(), 1, 1, &emptyCode, nil, nil, nil, nil)

	assert.Error(t, err)
	assert.Equal(t, "code cannot be empty", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestUpdateProduct_Error_EmptyName(t *testing.T) {
	mockRepo := new(MockProductRepo)
	service := service.NewProductService(mockRepo)

	existingProduct := &domain.Product{ID: 1, UserID: 1, Code: "PROD001", Name: "Test Product"}
	mockRepo.On("FindByIDAndUserID", mock.Anything, uint(1), uint(1)).Return(existingProduct, nil)

	emptyName := ""
	_, err := service.UpdateProduct(context.Background(), 1, 1, nil, &emptyName, nil, nil, nil)

	assert.Error(t, err)
	assert.Equal(t, "name cannot be empty", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestUpdateProduct_Error_NegativePrice(t *testing.T) {
	mockRepo := new(MockProductRepo)
	service := service.NewProductService(mockRepo)

	existingProduct := &domain.Product{ID: 1, UserID: 1, Code: "PROD001", Name: "Test Product"}
	mockRepo.On("FindByIDAndUserID", mock.Anything, uint(1), uint(1)).Return(existingProduct, nil)

	negativePrice := -10.0
	_, err := service.UpdateProduct(context.Background(), 1, 1, nil, nil, nil, &negativePrice, nil)

	assert.Error(t, err)
	assert.Equal(t, "price cannot be negative", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestUpdateProduct_Error_NegativeStock(t *testing.T) {
	mockRepo := new(MockProductRepo)
	service := service.NewProductService(mockRepo)

	existingProduct := &domain.Product{ID: 1, UserID: 1, Code: "PROD001", Name: "Test Product"}
	mockRepo.On("FindByIDAndUserID", mock.Anything, uint(1), uint(1)).Return(existingProduct, nil)

	negativeStock := -5
	_, err := service.UpdateProduct(context.Background(), 1, 1, nil, nil, nil, nil, &negativeStock)

	assert.Error(t, err)
	assert.Equal(t, "stock cannot be negative", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestUpdateProductPartial_Error_CodeConflict(t *testing.T) {
	mockRepo := new(MockProductRepo)
	service := service.NewProductService(mockRepo)

	existingProduct := &domain.Product{ID: 1, UserID: 1, Code: "PROD001", Name: "Test Product"}
	conflictingProduct := &domain.Product{ID: 2, UserID: 1, Code: "PROD002", Name: "Other Product"}

	mockRepo.On("FindByIDAndUserID", mock.Anything, uint(1), uint(1)).Return(existingProduct, nil)
	mockRepo.On("FindByCodeAndUserID", mock.Anything, "PROD002", uint(1)).Return(conflictingProduct, nil)

	newCode := "PROD002"
	_, err := service.UpdateProduct(context.Background(), 1, 1, &newCode, nil, nil, nil, nil)

	assert.Error(t, err)
	assert.Equal(t, "product code already exists for this user", err.Error())
	mockRepo.AssertExpectations(t)
}
