package service

import (
	"errors"
	"testing"

	"github.com/shoelfikar/voucher-management-system/internal/domain/entity"
	jwtPkg "github.com/shoelfikar/voucher-management-system/pkg/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserRepository is a mock implementation of UserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) FindByEmail(email string) (*entity.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserRepository) Create(user *entity.User) error {
	args := m.Called(user)
	return args.Error(0)
}

// MockJWTService is a mock implementation of JWTService
type MockJWTService struct {
	mock.Mock
}

func (m *MockJWTService) GenerateToken(email string) (string, error) {
	args := m.Called(email)
	return args.String(0), args.Error(1)
}

func (m *MockJWTService) ValidateToken(token string) (*jwtPkg.Claims, error) {
	args := m.Called(token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*jwtPkg.Claims), args.Error(1)
}

func TestAuthService_Login_Success(t *testing.T) {
	// Arrange
	mockUserRepo := new(MockUserRepository)
	mockJWTService := new(MockJWTService)

	authService := NewAuthService(mockUserRepo, mockJWTService)

	email := "test@example.com"
	password := "password123"
	expectedToken := "mock.jwt.token"

	mockJWTService.On("GenerateToken", email).Return(expectedToken, nil)

	// Act
	token, user, err := authService.Login(email, password)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, email, user.Email)
	assert.Equal(t, expectedToken, token)
	mockJWTService.AssertExpectations(t)
}

func TestAuthService_Login_JWT_GenerateToken_Error(t *testing.T) {
	// Arrange
	mockUserRepo := new(MockUserRepository)
	mockJWTService := new(MockJWTService)

	authService := NewAuthService(mockUserRepo, mockJWTService)

	email := "test@example.com"
	password := "password123"
	expectedError := errors.New("failed to generate token")

	mockJWTService.On("GenerateToken", email).Return("", expectedError)

	// Act
	token, user, err := authService.Login(email, password)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Empty(t, token)
	assert.Nil(t, user)
	mockJWTService.AssertExpectations(t)
}

func TestAuthService_Login_EmptyEmail(t *testing.T) {
	// Arrange
	mockUserRepo := new(MockUserRepository)
	mockJWTService := new(MockJWTService)

	authService := NewAuthService(mockUserRepo, mockJWTService)

	email := ""
	password := "password123"
	expectedToken := "mock.jwt.token"

	mockJWTService.On("GenerateToken", email).Return(expectedToken, nil)

	// Act
	token, user, err := authService.Login(email, password)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, email, user.Email)
	assert.Equal(t, expectedToken, token)
	mockJWTService.AssertExpectations(t)
}

func TestAuthService_Login_EmptyPassword(t *testing.T) {
	// Arrange
	mockUserRepo := new(MockUserRepository)
	mockJWTService := new(MockJWTService)

	authService := NewAuthService(mockUserRepo, mockJWTService)

	email := "test@example.com"
	password := ""
	expectedToken := "mock.jwt.token"

	mockJWTService.On("GenerateToken", email).Return(expectedToken, nil)

	// Act
	token, user, err := authService.Login(email, password)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, email, user.Email)
	assert.Equal(t, expectedToken, token)
	mockJWTService.AssertExpectations(t)
}
