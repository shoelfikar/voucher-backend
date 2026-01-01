package service

import (
	"github.com/shoelfikar/voucher-management-system/internal/domain/entity"
	"github.com/shoelfikar/voucher-management-system/internal/domain/repository"
	domainService "github.com/shoelfikar/voucher-management-system/internal/domain/service"
	"github.com/shoelfikar/voucher-management-system/pkg/jwt"
)

// authServiceImpl implements domain service.AuthService
type authServiceImpl struct {
	userRepo   repository.UserRepository
	jwtService jwt.JWTService
}

// NewAuthService creates a new auth service instance
func NewAuthService(userRepo repository.UserRepository, jwtService jwt.JWTService) domainService.AuthService {
	return &authServiceImpl{
		userRepo:   userRepo,
		jwtService: jwtService,
	}
}

// Login authenticates a user with dummy validation and returns a JWT token
func (s *authServiceImpl) Login(email, password string) (string, *entity.User, error) {
	// Dummy validation - accept any email/password combination
	// In production, you should:
	// 1. Find user by email from database
	// _, err := s.userRepo.FindByEmail(email)
	// if err != nil {
	// 	return "", nil, err
	// }
	// 2. Compare hashed password with bcrypt
	// 3. Return error if credentials are invalid

	user := &entity.User{
		Email: email,
	}

	token, err := s.jwtService.GenerateToken(email)
	if err != nil {
		return "", nil, err
	}

	return token, user, nil
}

func (s *authServiceImpl) Register(email, password string) (string, error) {
	return "", nil
}
