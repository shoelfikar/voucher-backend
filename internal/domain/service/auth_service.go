package service

import "github.com/shoelfikar/voucher-management-system/internal/domain/entity"

// AuthService defines the interface for authentication operations
type AuthService interface {
	// Login authenticates a user and returns a token
	Login(email, password string) (string, *entity.User, error)

	// Register new user
	Register(email, password string) (string, error)
}
