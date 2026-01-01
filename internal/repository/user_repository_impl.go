package repository

import (
	"github.com/shoelfikar/voucher-management-system/internal/domain/entity"
	"github.com/shoelfikar/voucher-management-system/internal/domain/repository"
	"gorm.io/gorm"
)

// userRepositoryImpl implements repository.UserRepository
type userRepositoryImpl struct {
	db *gorm.DB
}

// NewUserRepository creates a new user repository instance
func NewUserRepository(db *gorm.DB) repository.UserRepository {
	return &userRepositoryImpl{db: db}
}

// FindByEmail finds a user by email
func (r *userRepositoryImpl) FindByEmail(email string) (*entity.User, error) {
	var user entity.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Create creates a new user
func (r *userRepositoryImpl) Create(user *entity.User) error {
	return r.db.Create(user).Error
}
