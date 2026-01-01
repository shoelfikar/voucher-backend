package repository

import (
	"testing"

	"github.com/shoelfikar/voucher-management-system/internal/domain/entity"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	// Use in-memory SQLite database for testing
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Auto migrate the schema
	err = db.AutoMigrate(&entity.User{})
	if err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}

	return db
}

func TestUserRepository_Create_Success(t *testing.T) {
	// Arrange
	db := setupTestDB(t)
	repo := NewUserRepository(db)

	user := &entity.User{
		Email:    "test@example.com",
		Password: "hashed_password",
	}

	// Act
	err := repo.Create(user)

	// Assert
	assert.NoError(t, err)
	assert.NotZero(t, user.ID)
}

func TestUserRepository_Create_DuplicateEmail(t *testing.T) {
	// Arrange
	db := setupTestDB(t)
	repo := NewUserRepository(db)

	user1 := &entity.User{
		Email:    "test@example.com",
		Password: "hashed_password",
	}

	user2 := &entity.User{
		Email:    "test@example.com",
		Password: "another_password",
	}

	// Act
	err1 := repo.Create(user1)
	err2 := repo.Create(user2)

	// Assert
	assert.NoError(t, err1)
	assert.Error(t, err2)
}

func TestUserRepository_FindByEmail_Success(t *testing.T) {
	// Arrange
	db := setupTestDB(t)
	repo := NewUserRepository(db)

	user := &entity.User{
		Email:    "test@example.com",
		Password: "hashed_password",
	}

	err := repo.Create(user)
	assert.NoError(t, err)

	// Act
	foundUser, err := repo.FindByEmail("test@example.com")

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, foundUser)
	assert.Equal(t, user.Email, foundUser.Email)
	assert.Equal(t, user.Password, foundUser.Password)
}

func TestUserRepository_FindByEmail_NotFound(t *testing.T) {
	// Arrange
	db := setupTestDB(t)
	repo := NewUserRepository(db)

	// Act
	foundUser, err := repo.FindByEmail("nonexistent@example.com")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, foundUser)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
}

func TestUserRepository_FindByEmail_EmptyEmail(t *testing.T) {
	// Arrange
	db := setupTestDB(t)
	repo := NewUserRepository(db)

	// Act
	foundUser, err := repo.FindByEmail("")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, foundUser)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
}

func TestUserRepository_Multiple_Users(t *testing.T) {
	// Arrange
	db := setupTestDB(t)
	repo := NewUserRepository(db)

	users := []*entity.User{
		{Email: "user1@example.com", Password: "password1"},
		{Email: "user2@example.com", Password: "password2"},
		{Email: "user3@example.com", Password: "password3"},
	}

	// Act
	for _, user := range users {
		err := repo.Create(user)
		assert.NoError(t, err)
	}

	// Assert - Find each user
	for _, user := range users {
		foundUser, err := repo.FindByEmail(user.Email)
		assert.NoError(t, err)
		assert.NotNil(t, foundUser)
		assert.Equal(t, user.Email, foundUser.Email)
	}
}
