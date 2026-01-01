package service

import (
	"mime/multipart"

	"github.com/shoelfikar/voucher-management-system/internal/delivery/http/request"
	"github.com/shoelfikar/voucher-management-system/internal/domain/entity"
)

// ImportResult represents the result of CSV import
type ImportResult struct {
	TotalRows int           `json:"total_rows"`
	Success   int           `json:"success"`
	Failed    int           `json:"failed"`
	Errors    []ImportError `json:"errors,omitempty"`
}

// ImportError represents an error during CSV import
type ImportError struct {
	Row   int    `json:"row"`
	Error string `json:"error"`
}

// BatchImportResult represents the result of batch import
type BatchImportResult struct {
	TotalReceived  int      `json:"total_received"`
	Inserted       int      `json:"inserted"`
	Duplicates     int      `json:"duplicates"`
	DuplicateCodes []string `json:"duplicate_codes"`
	Errors         []string `json:"errors"`
}

// VoucherService defines the interface for voucher business logic
type VoucherService interface {
	// GetAll retrieves all vouchers with pagination and filters
	GetAll(page, limit int, search, sortBy, sortOrder string) ([]*entity.Voucher, int64, error)

	// GetByID retrieves a voucher by ID
	GetByID(id uint) (*entity.Voucher, error)

	// Create creates a new voucher with validation
	Create(req *request.CreateVoucherRequest) (*entity.Voucher, error)

	// Update updates an existing voucher with validation
	Update(id uint, req *request.UpdateVoucherRequest) (*entity.Voucher, error)

	// Delete deletes a voucher by ID
	Delete(id uint) error

	// ImportVouchers imports vouchers from CSV file
	ImportVouchers(file multipart.File) (*ImportResult, error)

	// ImportBatch imports a batch of vouchers with duplicate checking
	ImportBatch(vouchers []request.CreateVoucherRequest) (*BatchImportResult, error)

	// ExportVouchers exports all vouchers to CSV format
	ExportVouchers() ([]byte, error)
}
