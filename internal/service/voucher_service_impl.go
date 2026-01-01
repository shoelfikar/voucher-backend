package service

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"mime/multipart"
	"strconv"
	"strings"
	"time"

	"github.com/shoelfikar/voucher-management-system/internal/delivery/http/request"
	"github.com/shoelfikar/voucher-management-system/internal/domain/entity"
	"github.com/shoelfikar/voucher-management-system/internal/domain/repository"
	domainService "github.com/shoelfikar/voucher-management-system/internal/domain/service"
	"gorm.io/gorm"
)

// voucherServiceImpl implements domain service.VoucherService
type voucherServiceImpl struct {
	voucherRepo repository.VoucherRepository
}

// NewVoucherService creates a new voucher service instance
func NewVoucherService(voucherRepo repository.VoucherRepository) domainService.VoucherService {
	return &voucherServiceImpl{
		voucherRepo: voucherRepo,
	}
}

// GetAll retrieves all vouchers with pagination and filters
func (s *voucherServiceImpl) GetAll(page, limit int, search, sortBy, sortOrder string) ([]*entity.Voucher, int64, error) {
	return s.voucherRepo.FindAll(page, limit, search, sortBy, sortOrder)
}

// GetByID retrieves a voucher by ID
func (s *voucherServiceImpl) GetByID(id uint) (*entity.Voucher, error) {
	voucher, err := s.voucherRepo.FindByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("voucher not found")
		}
		return nil, err
	}
	return voucher, nil
}

// Create creates a new voucher with validation
func (s *voucherServiceImpl) Create(req *request.CreateVoucherRequest) (*entity.Voucher, error) {
	// Check if voucher code already exists
	existing, err := s.voucherRepo.FindByVoucherCode(req.VoucherCode)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("voucher code already exists")
	}

	// Parse expiry date
	expiryDate, err := time.Parse("2006-01-02", req.ExpiryDate)
	if err != nil {
		return nil, errors.New("invalid date format, expected YYYY-MM-DD")
	}

	// Validate expiry date is in the future or today
	// Get today's date at midnight in local timezone for proper comparison
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	// Convert expiry date to local timezone for comparison
	expiryDateLocal := time.Date(expiryDate.Year(), expiryDate.Month(), expiryDate.Day(), 0, 0, 0, 0, now.Location())
	if expiryDateLocal.Before(today) {
		return nil, errors.New("expiry date must be today or in the future")
	}

	// Create voucher entity
	voucher := &entity.Voucher{
		VoucherCode:     req.VoucherCode,
		DiscountPercent: req.DiscountPercent,
		ExpiryDate:      expiryDate,
	}

	// Save to database
	err = s.voucherRepo.Create(voucher)
	if err != nil {
		return nil, err
	}

	return voucher, nil
}

// Update updates an existing voucher with validation
func (s *voucherServiceImpl) Update(id uint, req *request.UpdateVoucherRequest) (*entity.Voucher, error) {
	// Check if voucher exists
	voucher, err := s.voucherRepo.FindByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("voucher not found")
		}
		return nil, err
	}

	// Check if voucher code is being changed and if new code already exists
	if req.VoucherCode != voucher.VoucherCode {
		existing, err := s.voucherRepo.FindByVoucherCode(req.VoucherCode)
		if err != nil && err != gorm.ErrRecordNotFound {
			return nil, err
		}
		if existing != nil {
			return nil, errors.New("voucher code already exists")
		}
	}

	// Parse expiry date
	expiryDate, err := time.Parse("2006-01-02", req.ExpiryDate)
	if err != nil {
		return nil, errors.New("invalid date format, expected YYYY-MM-DD")
	}

	// Validate expiry date is in the future or today
	// Get today's date at midnight in local timezone for proper comparison
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	// Convert expiry date to local timezone for comparison
	expiryDateLocal := time.Date(expiryDate.Year(), expiryDate.Month(), expiryDate.Day(), 0, 0, 0, 0, now.Location())
	if expiryDateLocal.Before(today) {
		return nil, errors.New("expiry date must be today or in the future")
	}

	// Update voucher fields
	voucher.VoucherCode = req.VoucherCode
	voucher.DiscountPercent = req.DiscountPercent
	voucher.ExpiryDate = expiryDate

	// Save to database
	err = s.voucherRepo.Update(voucher)
	if err != nil {
		return nil, err
	}

	return voucher, nil
}

// Delete deletes a voucher by ID (soft delete)
func (s *voucherServiceImpl) Delete(id uint) error {
	// Check if voucher exists
	_, err := s.voucherRepo.FindByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New("voucher not found")
		}
		return err
	}

	// Soft delete
	return s.voucherRepo.Delete(id)
}

// ImportVouchers imports vouchers from CSV file
func (s *voucherServiceImpl) ImportVouchers(file multipart.File) (*domainService.ImportResult, error) {
	// Read CSV file
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV file: %w", err)
	}

	if len(records) < 2 {
		return nil, errors.New("CSV file is empty or has no data rows")
	}

	result := &domainService.ImportResult{
		TotalRows: len(records) - 1,
		Errors:    []domainService.ImportError{},
	}

	var vouchers []*entity.Voucher

	// Process each row (skip header)
	for i, record := range records[1:] {
		rowNum := i + 2

		voucher, err := s.parseCSVRow(record, rowNum)
		if err != nil {
			result.Errors = append(result.Errors, domainService.ImportError{
				Row:   rowNum,
				Error: err.Error(),
			})
			result.Failed++
			continue
		}

		vouchers = append(vouchers, voucher)
	}

	// Bulk insert valid vouchers
	if len(vouchers) > 0 {
		err = s.voucherRepo.BulkCreate(vouchers)
		if err != nil {
			return nil, fmt.Errorf("failed to insert vouchers: %w", err)
		}
		result.Success = len(vouchers)
	}

	return result, nil
}

// parseCSVRow parses a single CSV row and returns a Voucher entity
func (s *voucherServiceImpl) parseCSVRow(record []string, rowNum int) (*entity.Voucher, error) {
	// Validate column count
	if len(record) < 3 {
		return nil, fmt.Errorf("insufficient columns (expected 3: voucher_code, discount_percent, expiry_date)")
	}

	voucherCode := strings.TrimSpace(record[0])
	if voucherCode == "" {
		return nil, errors.New("voucher code is required")
	}
	if len(voucherCode) > 50 {
		return nil, errors.New("voucher code exceeds 50 characters")
	}

	// Check if voucher code already exists
	existing, err := s.voucherRepo.FindByVoucherCode(voucherCode)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("failed to check voucher code: %w", err)
	}
	if existing != nil {
		return nil, fmt.Errorf("voucher code '%s' already exists", voucherCode)
	}

	// Parse discount percent
	discountStr := strings.TrimSpace(record[1])
	discountPercent, err := strconv.ParseFloat(discountStr, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid discount percent '%s': must be a number", discountStr)
	}
	if discountPercent < 1 || discountPercent > 100 {
		return nil, fmt.Errorf("discount percent %.2f out of range (must be 1-100)", discountPercent)
	}

	// Parse expiry date
	expiryDateStr := strings.TrimSpace(record[2])
	expiryDate, err := time.Parse("2006-01-02", expiryDateStr)
	if err != nil {
		return nil, fmt.Errorf("invalid date format '%s': expected YYYY-MM-DD", expiryDateStr)
	}

	// Validate expiry date is in the future
	today := time.Now().Truncate(24 * time.Hour)
	if expiryDate.Before(today) {
		return nil, fmt.Errorf("expiry date %s must be today or in the future", expiryDateStr)
	}

	voucher := &entity.Voucher{
		VoucherCode:     voucherCode,
		DiscountPercent: discountPercent,
		ExpiryDate:      expiryDate,
	}

	return voucher, nil
}

// ExportVouchers exports all vouchers to CSV format
func (s *voucherServiceImpl) ExportVouchers() ([]byte, error) {
	vouchers, _, err := s.voucherRepo.FindAll(1, 100000, "", "created_at", "asc")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch vouchers: %w", err)
	}

	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	header := []string{"voucher_code", "discount_percent", "expiry_date"}
	if err := writer.Write(header); err != nil {
		return nil, fmt.Errorf("failed to write CSV header: %w", err)
	}

	for _, voucher := range vouchers {
		record := []string{
			voucher.VoucherCode,
			fmt.Sprintf("%.2f", voucher.DiscountPercent),
			voucher.ExpiryDate.Format("2006-01-02"),
		}
		if err := writer.Write(record); err != nil {
			return nil, fmt.Errorf("failed to write CSV row: %w", err)
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return nil, fmt.Errorf("failed to flush CSV writer: %w", err)
	}

	return buf.Bytes(), nil
}

// ImportBatch imports a batch of vouchers with duplicate checking
func (s *voucherServiceImpl) ImportBatch(vouchers []request.CreateVoucherRequest) (*domainService.BatchImportResult, error) {
	result := &domainService.BatchImportResult{
		TotalReceived:  len(vouchers),
		DuplicateCodes: []string{},
		Errors:         []string{},
	}

	// Step 1: Extract all voucher codes
	voucherCodes := make([]string, len(vouchers))
	for i, v := range vouchers {
		voucherCodes[i] = v.VoucherCode
	}

	// Step 2: Check duplicates with IN query
	existingCodes, err := s.voucherRepo.CheckDuplicateCodes(voucherCodes)
	if err != nil {
		return nil, err
	}

	// Step 3: Create map for quick lookup
	duplicateMap := make(map[string]bool)
	for _, code := range existingCodes {
		duplicateMap[code] = true
	}

	// Step 4: Filter valid vouchers
	validVouchers := []*entity.Voucher{}
	for _, voucherReq := range vouchers {
		// Check if duplicate
		if duplicateMap[voucherReq.VoucherCode] {
			result.Duplicates++
			result.DuplicateCodes = append(result.DuplicateCodes, voucherReq.VoucherCode)
			continue
		}

		// Validate and convert
		voucher, err := s.validateAndConvert(&voucherReq)
		if err != nil {
			result.Errors = append(result.Errors,
				fmt.Sprintf("Code %s: %s", voucherReq.VoucherCode, err.Error()))
			continue
		}

		validVouchers = append(validVouchers, voucher)
	}

	// Step 5: Bulk insert valid vouchers
	if len(validVouchers) > 0 {
		err = s.voucherRepo.BulkCreate(validVouchers)
		if err != nil {
			return nil, err
		}
		result.Inserted = len(validVouchers)
	}

	return result, nil
}

// validateAndConvert validates a voucher request and converts it to entity
func (s *voucherServiceImpl) validateAndConvert(req *request.CreateVoucherRequest) (*entity.Voucher, error) {
	// Validate voucher code
	if req.VoucherCode == "" {
		return nil, errors.New("voucher code is required")
	}
	if len(req.VoucherCode) > 50 {
		return nil, errors.New("voucher code exceeds 50 characters")
	}

	// Validate discount percent
	if req.DiscountPercent < 1 || req.DiscountPercent > 100 {
		return nil, fmt.Errorf("discount percent %.2f out of range (must be 1-100)", req.DiscountPercent)
	}

	// Parse expiry date
	expiryDate, err := time.Parse("2006-01-02", req.ExpiryDate)
	if err != nil {
		return nil, fmt.Errorf("invalid date format '%s': expected YYYY-MM-DD", req.ExpiryDate)
	}

	// Validate expiry date is in the future
	today := time.Now().Truncate(24 * time.Hour)
	if expiryDate.Before(today) {
		return nil, fmt.Errorf("expiry date %s must be today or in the future", req.ExpiryDate)
	}

	voucher := &entity.Voucher{
		VoucherCode:     req.VoucherCode,
		DiscountPercent: req.DiscountPercent,
		ExpiryDate:      expiryDate,
	}

	return voucher, nil
}
