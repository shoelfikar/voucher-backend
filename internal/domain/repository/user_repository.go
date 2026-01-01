package repository

import "github.com/shoelfikar/voucher-management-system/internal/domain/entity"

// UserRepository defines the interface for user data operations
type UserRepository interface {
	FindByEmail(email string) (*entity.User, error)
	Create(user *entity.User) error
}
