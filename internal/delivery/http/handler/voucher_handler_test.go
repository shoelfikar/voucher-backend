package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shoelfikar/voucher-management-system/internal/delivery/http/request"
	"github.com/shoelfikar/voucher-management-system/internal/domain/entity"
	"github.com/shoelfikar/voucher-management-system/internal/domain/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockVoucherService is a mock implementation of VoucherService
type MockVoucherService struct {
	mock.Mock
}

func (m *MockVoucherService) GetAll(page, limit int, search, sortBy, sortOrder string) ([]*entity.Voucher, int64, error) {
	args := m.Called(page, limit, search, sortBy, sortOrder)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*entity.Voucher), args.Get(1).(int64), args.Error(2)
}

func (m *MockVoucherService) GetByID(id uint) (*entity.Voucher, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Voucher), args.Error(1)
}

func (m *MockVoucherService) Create(req *request.CreateVoucherRequest) (*entity.Voucher, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Voucher), args.Error(1)
}

func (m *MockVoucherService) Update(id uint, req *request.UpdateVoucherRequest) (*entity.Voucher, error) {
	args := m.Called(id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Voucher), args.Error(1)
}

func (m *MockVoucherService) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockVoucherService) ImportVouchers(file multipart.File) (*service.ImportResult, error) {
	args := m.Called(file)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*service.ImportResult), args.Error(1)
}

func (m *MockVoucherService) ImportBatch(vouchers []request.CreateVoucherRequest) (*service.BatchImportResult, error) {
	args := m.Called(vouchers)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*service.BatchImportResult), args.Error(1)
}

func (m *MockVoucherService) ExportVouchers() ([]byte, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]byte), args.Error(1)
}

func setupVoucherTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	return router
}

// Test GetAll Vouchers
func TestVoucherHandler_GetAll_Success(t *testing.T) {
	// Arrange
	mockService := new(MockVoucherService)
	voucherHandler := NewVoucherHandler(mockService)
	router := setupVoucherTestRouter()
	router.GET("/vouchers", voucherHandler.GetAll)

	vouchers := []*entity.Voucher{
		{ID: 1, VoucherCode: "TEST1", DiscountPercent: 10.0},
		{ID: 2, VoucherCode: "TEST2", DiscountPercent: 20.0},
	}
	total := int64(2)

	mockService.On("GetAll", 1, 10, "", "created_at", "desc").Return(vouchers, total, nil)

	req, _ := http.NewRequest("GET", "/vouchers?page=1&limit=10&sort_by=created_at&sort_order=desc", nil)
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

	mockService.AssertExpectations(t)
}

func TestVoucherHandler_GetAll_WithSearch(t *testing.T) {
	// Arrange
	mockService := new(MockVoucherService)
	voucherHandler := NewVoucherHandler(mockService)
	router := setupVoucherTestRouter()
	router.GET("/vouchers", voucherHandler.GetAll)

	vouchers := []*entity.Voucher{
		{ID: 1, VoucherCode: "TEST1", DiscountPercent: 10.0},
	}
	total := int64(1)

	mockService.On("GetAll", 1, 10, "TEST", "created_at", "desc").Return(vouchers, total, nil)

	req, _ := http.NewRequest("GET", "/vouchers?page=1&limit=10&search=TEST&sort_by=created_at&sort_order=desc", nil)
	w := httptest.NewRecorder()

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestVoucherHandler_GetAll_ServiceError(t *testing.T) {
	// Arrange
	mockService := new(MockVoucherService)
	voucherHandler := NewVoucherHandler(mockService)
	router := setupVoucherTestRouter()
	router.GET("/vouchers", voucherHandler.GetAll)

	serviceError := errors.New("database error")
	mockService.On("GetAll", 1, 10, "", "created_at", "desc").Return(nil, int64(0), serviceError)

	req, _ := http.NewRequest("GET", "/vouchers", nil)
	w := httptest.NewRecorder()

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "error", response["status"])

	mockService.AssertExpectations(t)
}

// Test GetByID
func TestVoucherHandler_GetByID_Success(t *testing.T) {
	// Arrange
	mockService := new(MockVoucherService)
	voucherHandler := NewVoucherHandler(mockService)
	router := setupVoucherTestRouter()
	router.GET("/vouchers/:id", voucherHandler.GetByID)

	voucher := &entity.Voucher{
		ID:              1,
		VoucherCode:     "TEST123",
		DiscountPercent: 10.0,
	}

	mockService.On("GetByID", uint(1)).Return(voucher, nil)

	req, _ := http.NewRequest("GET", "/vouchers/1", nil)
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

	mockService.AssertExpectations(t)
}

func TestVoucherHandler_GetByID_InvalidID(t *testing.T) {
	// Arrange
	mockService := new(MockVoucherService)
	voucherHandler := NewVoucherHandler(mockService)
	router := setupVoucherTestRouter()
	router.GET("/vouchers/:id", voucherHandler.GetByID)

	req, _ := http.NewRequest("GET", "/vouchers/invalid", nil)
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

func TestVoucherHandler_GetByID_NotFound(t *testing.T) {
	// Arrange
	mockService := new(MockVoucherService)
	voucherHandler := NewVoucherHandler(mockService)
	router := setupVoucherTestRouter()
	router.GET("/vouchers/:id", voucherHandler.GetByID)

	notFoundError := errors.New("voucher not found")
	mockService.On("GetByID", uint(999)).Return(nil, notFoundError)

	req, _ := http.NewRequest("GET", "/vouchers/999", nil)
	w := httptest.NewRecorder()

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusNotFound, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "error", response["status"])

	mockService.AssertExpectations(t)
}

// Test Create Voucher
func TestVoucherHandler_Create_Success(t *testing.T) {
	// Arrange
	mockService := new(MockVoucherService)
	voucherHandler := NewVoucherHandler(mockService)
	router := setupVoucherTestRouter()
	router.POST("/vouchers", voucherHandler.Create)

	tomorrow := time.Now().Add(24 * time.Hour).Format("2006-01-02")
	createReq := request.CreateVoucherRequest{
		VoucherCode:     "TEST123",
		DiscountPercent: 10.0,
		ExpiryDate:      tomorrow,
	}

	createdVoucher := &entity.Voucher{
		ID:              1,
		VoucherCode:     createReq.VoucherCode,
		DiscountPercent: createReq.DiscountPercent,
	}

	mockService.On("Create", mock.AnythingOfType("*request.CreateVoucherRequest")).Return(createdVoucher, nil)

	requestBody, _ := json.Marshal(createReq)
	req, _ := http.NewRequest("POST", "/vouchers", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusCreated, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "success", response["status"])
	assert.NotNil(t, response["data"])

	mockService.AssertExpectations(t)
}

func TestVoucherHandler_Create_InvalidJSON(t *testing.T) {
	// Arrange
	mockService := new(MockVoucherService)
	voucherHandler := NewVoucherHandler(mockService)
	router := setupVoucherTestRouter()
	router.POST("/vouchers", voucherHandler.Create)

	invalidJSON := []byte(`{"voucher_code": "TEST123", "discount_percent":}`)
	req, _ := http.NewRequest("POST", "/vouchers", bytes.NewBuffer(invalidJSON))
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

func TestVoucherHandler_Create_ValidationError(t *testing.T) {
	// Arrange
	mockService := new(MockVoucherService)
	voucherHandler := NewVoucherHandler(mockService)
	router := setupVoucherTestRouter()
	router.POST("/vouchers", voucherHandler.Create)

	createReq := request.CreateVoucherRequest{
		VoucherCode:     "",
		DiscountPercent: 10.0,
		ExpiryDate:      "2025-12-31",
	}

	requestBody, _ := json.Marshal(createReq)
	req, _ := http.NewRequest("POST", "/vouchers", bytes.NewBuffer(requestBody))
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

func TestVoucherHandler_Create_ServiceError(t *testing.T) {
	// Arrange
	mockService := new(MockVoucherService)
	voucherHandler := NewVoucherHandler(mockService)
	router := setupVoucherTestRouter()
	router.POST("/vouchers", voucherHandler.Create)

	tomorrow := time.Now().Add(24 * time.Hour).Format("2006-01-02")
	createReq := request.CreateVoucherRequest{
		VoucherCode:     "TEST123",
		DiscountPercent: 10.0,
		ExpiryDate:      tomorrow,
	}

	serviceError := errors.New("voucher code already exists")
	mockService.On("Create", mock.AnythingOfType("*request.CreateVoucherRequest")).Return(nil, serviceError)

	requestBody, _ := json.Marshal(createReq)
	req, _ := http.NewRequest("POST", "/vouchers", bytes.NewBuffer(requestBody))
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

	mockService.AssertExpectations(t)
}

// Test Update Voucher
func TestVoucherHandler_Update_Success(t *testing.T) {
	// Arrange
	mockService := new(MockVoucherService)
	voucherHandler := NewVoucherHandler(mockService)
	router := setupVoucherTestRouter()
	router.PUT("/vouchers/:id", voucherHandler.Update)

	tomorrow := time.Now().Add(24 * time.Hour).Format("2006-01-02")
	updateReq := request.UpdateVoucherRequest{
		VoucherCode:     "UPDATED123",
		DiscountPercent: 15.0,
		ExpiryDate:      tomorrow,
	}

	updatedVoucher := &entity.Voucher{
		ID:              1,
		VoucherCode:     updateReq.VoucherCode,
		DiscountPercent: updateReq.DiscountPercent,
	}

	mockService.On("Update", uint(1), mock.AnythingOfType("*request.UpdateVoucherRequest")).Return(updatedVoucher, nil)

	requestBody, _ := json.Marshal(updateReq)
	req, _ := http.NewRequest("PUT", "/vouchers/1", bytes.NewBuffer(requestBody))
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

	mockService.AssertExpectations(t)
}

func TestVoucherHandler_Update_InvalidID(t *testing.T) {
	// Arrange
	mockService := new(MockVoucherService)
	voucherHandler := NewVoucherHandler(mockService)
	router := setupVoucherTestRouter()
	router.PUT("/vouchers/:id", voucherHandler.Update)

	updateReq := request.UpdateVoucherRequest{
		VoucherCode:     "UPDATED123",
		DiscountPercent: 15.0,
		ExpiryDate:      "2025-12-31",
	}

	requestBody, _ := json.Marshal(updateReq)
	req, _ := http.NewRequest("PUT", "/vouchers/invalid", bytes.NewBuffer(requestBody))
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

// Test Delete Voucher
func TestVoucherHandler_Delete_Success(t *testing.T) {
	// Arrange
	mockService := new(MockVoucherService)
	voucherHandler := NewVoucherHandler(mockService)
	router := setupVoucherTestRouter()
	router.DELETE("/vouchers/:id", voucherHandler.Delete)

	mockService.On("Delete", uint(1)).Return(nil)

	req, _ := http.NewRequest("DELETE", "/vouchers/1", nil)
	w := httptest.NewRecorder()

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "success", response["status"])

	mockService.AssertExpectations(t)
}

func TestVoucherHandler_Delete_InvalidID(t *testing.T) {
	// Arrange
	mockService := new(MockVoucherService)
	voucherHandler := NewVoucherHandler(mockService)
	router := setupVoucherTestRouter()
	router.DELETE("/vouchers/:id", voucherHandler.Delete)

	req, _ := http.NewRequest("DELETE", "/vouchers/invalid", nil)
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

func TestVoucherHandler_Delete_NotFound(t *testing.T) {
	// Arrange
	mockService := new(MockVoucherService)
	voucherHandler := NewVoucherHandler(mockService)
	router := setupVoucherTestRouter()
	router.DELETE("/vouchers/:id", voucherHandler.Delete)

	notFoundError := errors.New("voucher not found")
	mockService.On("Delete", uint(999)).Return(notFoundError)

	req, _ := http.NewRequest("DELETE", "/vouchers/999", nil)
	w := httptest.NewRecorder()

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusNotFound, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "error", response["status"])

	mockService.AssertExpectations(t)
}
