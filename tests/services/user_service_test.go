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

// Mock of UserRepository
type MockUserRepo struct {
	mock.Mock
}

func (m *MockUserRepo) Create(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}
func (m *MockUserRepo) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}
func (m *MockUserRepo) FindByID(ctx context.Context, id uint) (*domain.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func TestRegister_Success(t *testing.T) {
	mockRepo := new(MockUserRepo)
	service := service.NewUserService(mockRepo)

	mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.User")).Return(nil)

	user, err := service.Register(context.Background(), "Test", "test@mail.com", "password123")
	assert.NoError(t, err)
	assert.Equal(t, "Test", user.Name)
	assert.Equal(t, "test@mail.com", user.Email)
	mockRepo.AssertExpectations(t)
}

func TestAuthenticate_Fail_NotFound(t *testing.T) {
	mockRepo := new(MockUserRepo)
	service := service.NewUserService(mockRepo)

	mockRepo.On("FindByEmail", mock.Anything, "notfound@mail.com").Return(nil, errors.New("not found"))

	_, err := service.Authenticate(context.Background(), "notfound@mail.com", "password123")
	assert.Error(t, err)
}
