package service

import (
	"errors"
	"testing"
	"time"

	"github.com/shoelfikar/voucher-management-system/internal/delivery/http/request"
	"github.com/shoelfikar/voucher-management-system/internal/domain/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// MockVoucherRepository is a mock implementation of VoucherRepository
type MockVoucherRepository struct {
	mock.Mock
}

func (m *MockVoucherRepository) FindAll(page, limit int, search, sortBy, sortOrder string) ([]*entity.Voucher, int64, error) {
	args := m.Called(page, limit, search, sortBy, sortOrder)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*entity.Voucher), args.Get(1).(int64), args.Error(2)
}

func (m *MockVoucherRepository) FindByID(id uint) (*entity.Voucher, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Voucher), args.Error(1)
}

func (m *MockVoucherRepository) Create(voucher *entity.Voucher) error {
	args := m.Called(voucher)
	return args.Error(0)
}

func (m *MockVoucherRepository) Update(voucher *entity.Voucher) error {
	args := m.Called(voucher)
	return args.Error(0)
}

func (m *MockVoucherRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockVoucherRepository) FindByVoucherCode(code string) (*entity.Voucher, error) {
	args := m.Called(code)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Voucher), args.Error(1)
}

func (m *MockVoucherRepository) BulkCreate(vouchers []*entity.Voucher) error {
	args := m.Called(vouchers)
	return args.Error(0)
}

func (m *MockVoucherRepository) CheckDuplicateCodes(codes []string) ([]string, error) {
	args := m.Called(codes)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]string), args.Error(1)
}

// Test Create Voucher
func TestVoucherService_Create_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockVoucherRepository)
	voucherService := NewVoucherService(mockRepo)

	tomorrow := time.Now().Add(24 * time.Hour).Format("2006-01-02")
	req := &request.CreateVoucherRequest{
		VoucherCode:     "TEST123",
		DiscountPercent: 10.0,
		ExpiryDate:      tomorrow,
	}

	mockRepo.On("FindByVoucherCode", req.VoucherCode).Return((*entity.Voucher)(nil), nil)
	mockRepo.On("Create", mock.AnythingOfType("*entity.Voucher")).Return(nil)

	// Act
	voucher, err := voucherService.Create(req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, voucher)
	assert.Equal(t, req.VoucherCode, voucher.VoucherCode)
	assert.Equal(t, req.DiscountPercent, voucher.DiscountPercent)
	mockRepo.AssertExpectations(t)
}

func TestVoucherService_Create_DuplicateCode(t *testing.T) {
	// Arrange
	mockRepo := new(MockVoucherRepository)
	voucherService := NewVoucherService(mockRepo)

	tomorrow := time.Now().Add(24 * time.Hour).Format("2006-01-02")
	req := &request.CreateVoucherRequest{
		VoucherCode:     "TEST123",
		DiscountPercent: 10.0,
		ExpiryDate:      tomorrow,
	}

	existingVoucher := &entity.Voucher{
		ID:              1,
		VoucherCode:     "TEST123",
		DiscountPercent: 10.0,
	}

	mockRepo.On("FindByVoucherCode", req.VoucherCode).Return(existingVoucher, nil)

	// Act
	voucher, err := voucherService.Create(req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, voucher)
	assert.Contains(t, err.Error(), "already exists")
	mockRepo.AssertExpectations(t)
}

func TestVoucherService_Create_InvalidDateFormat(t *testing.T) {
	// Arrange
	mockRepo := new(MockVoucherRepository)
	voucherService := NewVoucherService(mockRepo)

	req := &request.CreateVoucherRequest{
		VoucherCode:     "TEST123",
		DiscountPercent: 10.0,
		ExpiryDate:      "invalid-date",
	}

	mockRepo.On("FindByVoucherCode", req.VoucherCode).Return((*entity.Voucher)(nil), nil)

	// Act
	voucher, err := voucherService.Create(req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, voucher)
	assert.Contains(t, err.Error(), "invalid date format")
	mockRepo.AssertExpectations(t)
}

func TestVoucherService_Create_PastExpiryDate(t *testing.T) {
	// Arrange
	mockRepo := new(MockVoucherRepository)
	voucherService := NewVoucherService(mockRepo)

	yesterday := time.Now().Add(-24 * time.Hour).Format("2006-01-02")
	req := &request.CreateVoucherRequest{
		VoucherCode:     "TEST123",
		DiscountPercent: 10.0,
		ExpiryDate:      yesterday,
	}

	mockRepo.On("FindByVoucherCode", req.VoucherCode).Return((*entity.Voucher)(nil), nil)

	// Act
	voucher, err := voucherService.Create(req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, voucher)
	assert.Contains(t, err.Error(), "must be today or in the future")
	mockRepo.AssertExpectations(t)
}

// Test Update Voucher
func TestVoucherService_Update_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockVoucherRepository)
	voucherService := NewVoucherService(mockRepo)

	tomorrow := time.Now().Add(24 * time.Hour).Format("2006-01-02")
	voucherID := uint(1)

	existingVoucher := &entity.Voucher{
		ID:              voucherID,
		VoucherCode:     "OLD123",
		DiscountPercent: 10.0,
	}

	req := &request.UpdateVoucherRequest{
		VoucherCode:     "NEW123",
		DiscountPercent: 15.0,
		ExpiryDate:      tomorrow,
	}

	mockRepo.On("FindByID", voucherID).Return(existingVoucher, nil)
	mockRepo.On("FindByVoucherCode", req.VoucherCode).Return(nil, gorm.ErrRecordNotFound)
	mockRepo.On("Update", mock.AnythingOfType("*entity.Voucher")).Return(nil)

	// Act
	voucher, err := voucherService.Update(voucherID, req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, voucher)
	assert.Equal(t, req.VoucherCode, voucher.VoucherCode)
	assert.Equal(t, req.DiscountPercent, voucher.DiscountPercent)
	mockRepo.AssertExpectations(t)
}

func TestVoucherService_Update_NotFound(t *testing.T) {
	// Arrange
	mockRepo := new(MockVoucherRepository)
	voucherService := NewVoucherService(mockRepo)

	tomorrow := time.Now().Add(24 * time.Hour).Format("2006-01-02")
	voucherID := uint(999)

	req := &request.UpdateVoucherRequest{
		VoucherCode:     "NEW123",
		DiscountPercent: 15.0,
		ExpiryDate:      tomorrow,
	}

	mockRepo.On("FindByID", voucherID).Return(nil, gorm.ErrRecordNotFound)

	// Act
	voucher, err := voucherService.Update(voucherID, req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, voucher)
	assert.Contains(t, err.Error(), "not found")
	mockRepo.AssertExpectations(t)
}

// Test Delete Voucher
func TestVoucherService_Delete_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockVoucherRepository)
	voucherService := NewVoucherService(mockRepo)

	voucherID := uint(1)
	existingVoucher := &entity.Voucher{
		ID:              voucherID,
		VoucherCode:     "TEST123",
		DiscountPercent: 10.0,
	}

	mockRepo.On("FindByID", voucherID).Return(existingVoucher, nil)
	mockRepo.On("Delete", voucherID).Return(nil)

	// Act
	err := voucherService.Delete(voucherID)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestVoucherService_Delete_NotFound(t *testing.T) {
	// Arrange
	mockRepo := new(MockVoucherRepository)
	voucherService := NewVoucherService(mockRepo)

	voucherID := uint(999)

	mockRepo.On("FindByID", voucherID).Return(nil, gorm.ErrRecordNotFound)

	// Act
	err := voucherService.Delete(voucherID)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
	mockRepo.AssertExpectations(t)
}

// Test GetByID
func TestVoucherService_GetByID_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockVoucherRepository)
	voucherService := NewVoucherService(mockRepo)

	voucherID := uint(1)
	expectedVoucher := &entity.Voucher{
		ID:              voucherID,
		VoucherCode:     "TEST123",
		DiscountPercent: 10.0,
	}

	mockRepo.On("FindByID", voucherID).Return(expectedVoucher, nil)

	// Act
	voucher, err := voucherService.GetByID(voucherID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, voucher)
	assert.Equal(t, expectedVoucher.ID, voucher.ID)
	assert.Equal(t, expectedVoucher.VoucherCode, voucher.VoucherCode)
	mockRepo.AssertExpectations(t)
}

func TestVoucherService_GetByID_NotFound(t *testing.T) {
	// Arrange
	mockRepo := new(MockVoucherRepository)
	voucherService := NewVoucherService(mockRepo)

	voucherID := uint(999)

	mockRepo.On("FindByID", voucherID).Return(nil, gorm.ErrRecordNotFound)

	// Act
	voucher, err := voucherService.GetByID(voucherID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, voucher)
	assert.Contains(t, err.Error(), "not found")
	mockRepo.AssertExpectations(t)
}

// Test GetAll
func TestVoucherService_GetAll_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockVoucherRepository)
	voucherService := NewVoucherService(mockRepo)

	expectedVouchers := []*entity.Voucher{
		{ID: 1, VoucherCode: "TEST1", DiscountPercent: 10.0},
		{ID: 2, VoucherCode: "TEST2", DiscountPercent: 20.0},
	}
	expectedTotal := int64(2)

	mockRepo.On("FindAll", 1, 10, "", "created_at", "desc").Return(expectedVouchers, expectedTotal, nil)

	// Act
	vouchers, total, err := voucherService.GetAll(1, 10, "", "created_at", "desc")

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedVouchers, vouchers)
	assert.Equal(t, expectedTotal, total)
	mockRepo.AssertExpectations(t)
}

func TestVoucherService_GetAll_WithSearch(t *testing.T) {
	// Arrange
	mockRepo := new(MockVoucherRepository)
	voucherService := NewVoucherService(mockRepo)

	search := "TEST"
	expectedVouchers := []*entity.Voucher{
		{ID: 1, VoucherCode: "TEST1", DiscountPercent: 10.0},
	}
	expectedTotal := int64(1)

	mockRepo.On("FindAll", 1, 10, search, "created_at", "desc").Return(expectedVouchers, expectedTotal, nil)

	// Act
	vouchers, total, err := voucherService.GetAll(1, 10, search, "created_at", "desc")

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedVouchers, vouchers)
	assert.Equal(t, expectedTotal, total)
	mockRepo.AssertExpectations(t)
}

func TestVoucherService_GetAll_Error(t *testing.T) {
	// Arrange
	mockRepo := new(MockVoucherRepository)
	voucherService := NewVoucherService(mockRepo)

	expectedError := errors.New("database error")

	mockRepo.On("FindAll", 1, 10, "", "created_at", "desc").Return(nil, int64(0), expectedError)

	// Act
	vouchers, total, err := voucherService.GetAll(1, 10, "", "created_at", "desc")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, vouchers)
	assert.Equal(t, int64(0), total)
	assert.Equal(t, expectedError, err)
	mockRepo.AssertExpectations(t)
}
