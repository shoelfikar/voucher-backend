package repository

import (
	"testing"
	"time"

	"github.com/shoelfikar/voucher-management-system/internal/domain/entity"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupVoucherTestDB(t *testing.T) *gorm.DB {
	// Use in-memory SQLite database for testing
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Auto migrate the schema
	err = db.AutoMigrate(&entity.Voucher{})
	if err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}

	return db
}

func createTestVoucher(code string, discount float64) *entity.Voucher {
	return &entity.Voucher{
		VoucherCode:     code,
		DiscountPercent: discount,
		ExpiryDate:      time.Now().Add(24 * time.Hour),
	}
}

// Test Create
func TestVoucherRepository_Create_Success(t *testing.T) {
	// Arrange
	db := setupVoucherTestDB(t)
	repo := NewVoucherRepository(db)

	voucher := createTestVoucher("TEST123", 10.0)

	// Act
	err := repo.Create(voucher)

	// Assert
	assert.NoError(t, err)
	assert.NotZero(t, voucher.ID)
	assert.NotZero(t, voucher.CreatedAt)
	assert.NotZero(t, voucher.UpdatedAt)
}

func TestVoucherRepository_Create_DuplicateCode(t *testing.T) {
	// Arrange
	db := setupVoucherTestDB(t)
	repo := NewVoucherRepository(db)

	voucher1 := createTestVoucher("TEST123", 10.0)
	voucher2 := createTestVoucher("TEST123", 20.0)

	// Act
	err1 := repo.Create(voucher1)
	err2 := repo.Create(voucher2)

	// Assert
	assert.NoError(t, err1)
	assert.Error(t, err2)
}

// Test FindByID
func TestVoucherRepository_FindByID_Success(t *testing.T) {
	// Arrange
	db := setupVoucherTestDB(t)
	repo := NewVoucherRepository(db)

	voucher := createTestVoucher("TEST123", 10.0)
	err := repo.Create(voucher)
	assert.NoError(t, err)

	// Act
	foundVoucher, err := repo.FindByID(voucher.ID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, foundVoucher)
	assert.Equal(t, voucher.ID, foundVoucher.ID)
	assert.Equal(t, voucher.VoucherCode, foundVoucher.VoucherCode)
	assert.Equal(t, voucher.DiscountPercent, foundVoucher.DiscountPercent)
}

func TestVoucherRepository_FindByID_NotFound(t *testing.T) {
	// Arrange
	db := setupVoucherTestDB(t)
	repo := NewVoucherRepository(db)

	// Act
	foundVoucher, err := repo.FindByID(999)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, foundVoucher)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
}

// Test FindByVoucherCode
func TestVoucherRepository_FindByVoucherCode_Success(t *testing.T) {
	// Arrange
	db := setupVoucherTestDB(t)
	repo := NewVoucherRepository(db)

	voucher := createTestVoucher("TEST123", 10.0)
	err := repo.Create(voucher)
	assert.NoError(t, err)

	// Act
	foundVoucher, err := repo.FindByVoucherCode("TEST123")

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, foundVoucher)
	assert.Equal(t, voucher.VoucherCode, foundVoucher.VoucherCode)
}

func TestVoucherRepository_FindByVoucherCode_NotFound(t *testing.T) {
	// Arrange
	db := setupVoucherTestDB(t)
	repo := NewVoucherRepository(db)

	// Act
	foundVoucher, err := repo.FindByVoucherCode("NONEXISTENT")

	// Assert
	assert.NoError(t, err)
	assert.Nil(t, foundVoucher)
}

// Test Update
func TestVoucherRepository_Update_Success(t *testing.T) {
	// Arrange
	db := setupVoucherTestDB(t)
	repo := NewVoucherRepository(db)

	voucher := createTestVoucher("TEST123", 10.0)
	err := repo.Create(voucher)
	assert.NoError(t, err)

	// Act
	voucher.VoucherCode = "UPDATED123"
	voucher.DiscountPercent = 20.0
	err = repo.Update(voucher)

	// Assert
	assert.NoError(t, err)

	// Verify update
	foundVoucher, err := repo.FindByID(voucher.ID)
	assert.NoError(t, err)
	assert.Equal(t, "UPDATED123", foundVoucher.VoucherCode)
	assert.Equal(t, 20.0, foundVoucher.DiscountPercent)
}

func TestVoucherRepository_Update_NotFound(t *testing.T) {
	// Arrange
	db := setupVoucherTestDB(t)
	repo := NewVoucherRepository(db)

	voucher := &entity.Voucher{
		ID:              999,
		VoucherCode:     "TEST123",
		DiscountPercent: 10.0,
		ExpiryDate:      time.Now().Add(24 * time.Hour),
	}

	// Act
	err := repo.Update(voucher)

	// Assert
	// GORM's Save (used in Update) performs an upsert operation:
	// it will create the record if it doesn't exist
	// This is the expected behavior, not an error
	assert.NoError(t, err)

	// Verify the record was created by Save
	foundVoucher, findErr := repo.FindByID(999)
	assert.NoError(t, findErr)
	assert.NotNil(t, foundVoucher)
	assert.Equal(t, uint(999), foundVoucher.ID)
	assert.Equal(t, "TEST123", foundVoucher.VoucherCode)
}

// Test Delete (Soft Delete)
func TestVoucherRepository_Delete_Success(t *testing.T) {
	// Arrange
	db := setupVoucherTestDB(t)
	repo := NewVoucherRepository(db)

	voucher := createTestVoucher("TEST123", 10.0)
	err := repo.Create(voucher)
	assert.NoError(t, err)

	// Act
	err = repo.Delete(voucher.ID)

	// Assert
	assert.NoError(t, err)

	// Verify soft delete - should not find with normal query
	foundVoucher, err := repo.FindByID(voucher.ID)
	assert.Error(t, err)
	assert.Nil(t, foundVoucher)

	// Verify record still exists with Unscoped
	var deletedVoucher entity.Voucher
	err = db.Unscoped().First(&deletedVoucher, voucher.ID).Error
	assert.NoError(t, err)
	assert.NotZero(t, deletedVoucher.DeletedAt)
}

// Test FindAll
func TestVoucherRepository_FindAll_Success(t *testing.T) {
	// Arrange
	db := setupVoucherTestDB(t)
	repo := NewVoucherRepository(db)

	vouchers := []*entity.Voucher{
		createTestVoucher("TEST1", 10.0),
		createTestVoucher("TEST2", 20.0),
		createTestVoucher("TEST3", 30.0),
	}

	for _, v := range vouchers {
		err := repo.Create(v)
		assert.NoError(t, err)
	}

	// Act
	foundVouchers, total, err := repo.FindAll(1, 10, "", "created_at", "asc")

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 3, len(foundVouchers))
	assert.Equal(t, int64(3), total)
}

func TestVoucherRepository_FindAll_WithPagination(t *testing.T) {
	// Arrange
	db := setupVoucherTestDB(t)
	repo := NewVoucherRepository(db)

	// Create 5 vouchers
	for i := 1; i <= 5; i++ {
		voucher := createTestVoucher(string(rune(i))+"TEST", float64(i*10))
		err := repo.Create(voucher)
		assert.NoError(t, err)
	}

	// Act - Get page 1 with limit 2
	page1Vouchers, total, err := repo.FindAll(1, 2, "", "created_at", "asc")

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 2, len(page1Vouchers))
	assert.Equal(t, int64(5), total)

	// Act - Get page 2 with limit 2
	page2Vouchers, total, err := repo.FindAll(2, 2, "", "created_at", "asc")

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 2, len(page2Vouchers))
	assert.Equal(t, int64(5), total)

	// Act - Get page 3 with limit 2
	page3Vouchers, total, err := repo.FindAll(3, 2, "", "created_at", "asc")

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 1, len(page3Vouchers))
	assert.Equal(t, int64(5), total)
}

func TestVoucherRepository_FindAll_WithSearch(t *testing.T) {
	// Arrange
	db := setupVoucherTestDB(t)
	repo := NewVoucherRepository(db)

	vouchers := []*entity.Voucher{
		createTestVoucher("SUMMER2024", 10.0),
		createTestVoucher("WINTER2024", 20.0),
		createTestVoucher("SUMMER2025", 30.0),
	}

	for _, v := range vouchers {
		err := repo.Create(v)
		assert.NoError(t, err)
	}

	// Act
	foundVouchers, total, err := repo.FindAll(1, 10, "SUMMER", "created_at", "asc")

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 2, len(foundVouchers))
	assert.Equal(t, int64(2), total)
}

func TestVoucherRepository_FindAll_WithSorting(t *testing.T) {
	// Arrange
	db := setupVoucherTestDB(t)
	repo := NewVoucherRepository(db)

	vouchers := []*entity.Voucher{
		createTestVoucher("C_VOUCHER", 10.0),
		createTestVoucher("A_VOUCHER", 20.0),
		createTestVoucher("B_VOUCHER", 30.0),
	}

	for _, v := range vouchers {
		time.Sleep(1 * time.Millisecond) // Ensure different timestamps
		err := repo.Create(v)
		assert.NoError(t, err)
	}

	// Act - Sort by voucher_code ascending
	foundVouchers, _, err := repo.FindAll(1, 10, "", "voucher_code", "asc")

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 3, len(foundVouchers))
	assert.Equal(t, "A_VOUCHER", foundVouchers[0].VoucherCode)
	assert.Equal(t, "B_VOUCHER", foundVouchers[1].VoucherCode)
	assert.Equal(t, "C_VOUCHER", foundVouchers[2].VoucherCode)
}

func TestVoucherRepository_FindAll_ExcludesDeleted(t *testing.T) {
	// Arrange
	db := setupVoucherTestDB(t)
	repo := NewVoucherRepository(db)

	vouchers := []*entity.Voucher{
		createTestVoucher("TEST1", 10.0),
		createTestVoucher("TEST2", 20.0),
		createTestVoucher("TEST3", 30.0),
	}

	for _, v := range vouchers {
		err := repo.Create(v)
		assert.NoError(t, err)
	}

	// Delete one voucher
	err := repo.Delete(vouchers[1].ID)
	assert.NoError(t, err)

	// Act
	foundVouchers, total, err := repo.FindAll(1, 10, "", "created_at", "asc")

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 2, len(foundVouchers))
	assert.Equal(t, int64(2), total)
}

// Test BulkCreate
func TestVoucherRepository_BulkCreate_Success(t *testing.T) {
	// Arrange
	db := setupVoucherTestDB(t)
	repo := NewVoucherRepository(db)

	vouchers := []*entity.Voucher{
		createTestVoucher("BULK1", 10.0),
		createTestVoucher("BULK2", 20.0),
		createTestVoucher("BULK3", 30.0),
	}

	// Act
	err := repo.BulkCreate(vouchers)

	// Assert
	assert.NoError(t, err)

	// Verify all were created
	foundVouchers, total, err := repo.FindAll(1, 10, "", "created_at", "asc")
	assert.NoError(t, err)
	assert.Equal(t, 3, len(foundVouchers))
	assert.Equal(t, int64(3), total)
}

// Test CheckDuplicateCodes
func TestVoucherRepository_CheckDuplicateCodes_Success(t *testing.T) {
	// Arrange
	db := setupVoucherTestDB(t)
	repo := NewVoucherRepository(db)

	// Create existing vouchers
	existingVouchers := []*entity.Voucher{
		createTestVoucher("EXISTING1", 10.0),
		createTestVoucher("EXISTING2", 20.0),
	}

	for _, v := range existingVouchers {
		err := repo.Create(v)
		assert.NoError(t, err)
	}

	// Act - Check for duplicates
	codes := []string{"EXISTING1", "NEW1", "EXISTING2", "NEW2"}
	duplicates, err := repo.CheckDuplicateCodes(codes)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 2, len(duplicates))
	assert.Contains(t, duplicates, "EXISTING1")
	assert.Contains(t, duplicates, "EXISTING2")
}

func TestVoucherRepository_CheckDuplicateCodes_NoDuplicates(t *testing.T) {
	// Arrange
	db := setupVoucherTestDB(t)
	repo := NewVoucherRepository(db)

	// Act - Check for duplicates with no existing vouchers
	codes := []string{"NEW1", "NEW2", "NEW3"}
	duplicates, err := repo.CheckDuplicateCodes(codes)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 0, len(duplicates))
}
