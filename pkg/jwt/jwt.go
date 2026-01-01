package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWTService defines the interface for JWT operations
type JWTService interface {
	GenerateToken(email string) (string, error)
	ValidateToken(token string) (*Claims, error)
}

// Claims represents the JWT claims
type Claims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

// jwtService implements JWTService interface
type jwtService struct {
	secretKey  string
	expiration time.Duration
}

// NewJWTService creates a new JWT service instance
func NewJWTService(secretKey string, expiration time.Duration) JWTService {
	return &jwtService{
		secretKey:  secretKey,
		expiration: expiration,
	}
}

// GenerateToken generates a new JWT token for the given email
func (s *jwtService) GenerateToken(email string) (string, error) {
	claims := Claims{
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.expiration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateToken validates the JWT token and returns the claims
func (s *jwtService) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(s.secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
