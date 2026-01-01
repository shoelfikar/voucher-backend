package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/shoelfikar/voucher-management-system/internal/delivery/http/request"
	"github.com/shoelfikar/voucher-management-system/internal/domain/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockAuthService is a mock implementation of AuthService
type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) Login(email, password string) (string, *entity.User, error) {
	args := m.Called(email, password)
	if args.Get(1) == nil {
		return args.String(0), nil, args.Error(2)
	}
	return args.String(0), args.Get(1).(*entity.User), args.Error(2)
}

func (m *MockAuthService) Register(email, password string) (string, error) {
	args := m.Called(email, password)
	return args.String(0), args.Error(1)
}

func setupAuthTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	return router
}

func TestAuthHandler_Login_Success(t *testing.T) {
	// Arrange
	mockAuthService := new(MockAuthService)
	authHandler := NewAuthHandler(mockAuthService)
	router := setupAuthTestRouter()
	router.POST("/login", authHandler.Login)

	loginReq := request.LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	user := &entity.User{
		ID:    1,
		Email: loginReq.Email,
	}

	mockAuthService.On("Login", loginReq.Email, loginReq.Password).Return("mock.jwt.token", user, nil)

	requestBody, _ := json.Marshal(loginReq)
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "success", response["status"])
	assert.NotNil(t, response["data"])

	data := response["data"].(map[string]interface{})
	assert.Equal(t, "mock.jwt.token", data["token"])
	assert.NotNil(t, data["user"])

	mockAuthService.AssertExpectations(t)
}

func TestAuthHandler_Login_InvalidJSON(t *testing.T) {
	// Arrange
	mockAuthService := new(MockAuthService)
	authHandler := NewAuthHandler(mockAuthService)
	router := setupAuthTestRouter()
	router.POST("/login", authHandler.Login)

	invalidJSON := []byte(`{"email": "test@example.com", "password":}`)
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(invalidJSON))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "error", response["status"])
}

func TestAuthHandler_Login_MissingEmail(t *testing.T) {
	// Arrange
	mockAuthService := new(MockAuthService)
	authHandler := NewAuthHandler(mockAuthService)
	router := setupAuthTestRouter()
	router.POST("/login", authHandler.Login)

	loginReq := request.LoginRequest{
		Email:    "",
		Password: "password123",
	}

	requestBody, _ := json.Marshal(loginReq)
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "error", response["status"])
}

func TestAuthHandler_Login_MissingPassword(t *testing.T) {
	// Arrange
	mockAuthService := new(MockAuthService)
	authHandler := NewAuthHandler(mockAuthService)
	router := setupAuthTestRouter()
	router.POST("/login", authHandler.Login)

	loginReq := request.LoginRequest{
		Email:    "test@example.com",
		Password: "",
	}

	requestBody, _ := json.Marshal(loginReq)
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "error", response["status"])
}

func TestAuthHandler_Login_ServiceError(t *testing.T) {
	// Arrange
	mockAuthService := new(MockAuthService)
	authHandler := NewAuthHandler(mockAuthService)
	router := setupAuthTestRouter()
	router.POST("/login", authHandler.Login)

	loginReq := request.LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	serviceError := errors.New("service error")
	mockAuthService.On("Login", loginReq.Email, loginReq.Password).Return("", nil, serviceError)

	requestBody, _ := json.Marshal(loginReq)
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "error", response["status"])

	mockAuthService.AssertExpectations(t)
}

func TestAuthHandler_Login_InvalidEmailFormat(t *testing.T) {
	// Arrange
	mockAuthService := new(MockAuthService)
	authHandler := NewAuthHandler(mockAuthService)
	router := setupAuthTestRouter()
	router.POST("/login", authHandler.Login)

	loginReq := request.LoginRequest{
		Email:    "invalid-email",
		Password: "password123",
	}

	requestBody, _ := json.Marshal(loginReq)
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "error", response["status"])
}
