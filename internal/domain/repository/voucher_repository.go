package repository

import "github.com/shoelfikar/voucher-management-system/internal/domain/entity"

// VoucherRepository defines the interface for voucher data operations
type VoucherRepository interface {
	// FindAll retrieves all vouchers with pagination, search, and sorting
	FindAll(page, limit int, search, sortBy, sortOrder string) ([]*entity.Voucher, int64, error)

	// FindByID retrieves a voucher by ID
	FindByID(id uint) (*entity.Voucher, error)

	// Create creates a new voucher
	Create(voucher *entity.Voucher) error

	// Update updates an existing voucher
	Update(voucher *entity.Voucher) error

	// Delete soft deletes a voucher by ID
	Delete(id uint) error

	// FindByVoucherCode retrieves a voucher by voucher code
	FindByVoucherCode(code string) (*entity.Voucher, error)

	// BulkCreate creates multiple vouchers at once
	BulkCreate(vouchers []*entity.Voucher) error

	// CheckDuplicateCodes checks which voucher codes already exist
	CheckDuplicateCodes(codes []string) ([]string, error)
}
