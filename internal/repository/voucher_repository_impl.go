package repository

import (
	"github.com/shoelfikar/voucher-management-system/internal/domain/entity"
	"github.com/shoelfikar/voucher-management-system/internal/domain/repository"
	"gorm.io/gorm"
)

// voucherRepositoryImpl implements repository.VoucherRepository
type voucherRepositoryImpl struct {
	db *gorm.DB
}

// NewVoucherRepository creates a new voucher repository instance
func NewVoucherRepository(db *gorm.DB) repository.VoucherRepository {
	return &voucherRepositoryImpl{db: db}
}

// FindAll retrieves all vouchers with pagination, search, and sorting
func (r *voucherRepositoryImpl) FindAll(page, limit int, search, sortBy, sortOrder string) ([]*entity.Voucher, int64, error) {
	var vouchers []*entity.Voucher
	var total int64

	offset := (page - 1) * limit

	query := r.db.Model(&entity.Voucher{})

	if search != "" {
		query = query.Where("LOWER(voucher_code) LIKE LOWER(?)", "%"+search+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if sortBy != "" {
		order := sortBy + " " + sortOrder
		query = query.Order(order)
	} else {
		query = query.Order("created_at desc")
	}

	// Pagination
	err := query.Offset(offset).Limit(limit).Find(&vouchers).Error
	if err != nil {
		return nil, 0, err
	}

	return vouchers, total, nil
}

// FindByID retrieves a voucher by ID
func (r *voucherRepositoryImpl) FindByID(id uint) (*entity.Voucher, error) {
	var voucher entity.Voucher
	err := r.db.First(&voucher, id).Error
	if err != nil {
		return nil, err
	}
	return &voucher, nil
}

// Create creates a new voucher
func (r *voucherRepositoryImpl) Create(voucher *entity.Voucher) error {
	return r.db.Create(voucher).Error
}

// Update updates an existing voucher
func (r *voucherRepositoryImpl) Update(voucher *entity.Voucher) error {
	return r.db.Save(voucher).Error
}

// Delete soft deletes a voucher by ID
func (r *voucherRepositoryImpl) Delete(id uint) error {
	return r.db.Delete(&entity.Voucher{}, id).Error
}

// FindByVoucherCode retrieves a voucher by voucher code
func (r *voucherRepositoryImpl) FindByVoucherCode(code string) (*entity.Voucher, error) {
	var voucher entity.Voucher
	err := r.db.Where("voucher_code = ?", code).First(&voucher).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &voucher, nil
}

// BulkCreate creates multiple vouchers at once
func (r *voucherRepositoryImpl) BulkCreate(vouchers []*entity.Voucher) error {
	return r.db.Create(&vouchers).Error
}

// CheckDuplicateCodes checks which voucher codes already exist
func (r *voucherRepositoryImpl) CheckDuplicateCodes(codes []string) ([]string, error) {
	var existingCodes []string

	err := r.db.Model(&entity.Voucher{}).
		Where("voucher_code IN ?", codes).
		Pluck("voucher_code", &existingCodes).
		Error

	if err != nil {
		return nil, err
	}

	return existingCodes, nil
}
