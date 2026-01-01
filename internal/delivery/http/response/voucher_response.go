package response

import (
	"time"

	"github.com/shoelfikar/voucher-management-system/internal/domain/entity"
)

// VoucherResponse represents a single voucher in response
type VoucherResponse struct {
	ID              uint    `json:"id"`
	VoucherCode     string  `json:"voucher_code"`
	DiscountPercent float64 `json:"discount_percent"`
	ExpiryDate      string  `json:"expiry_date"`
	CreatedAt       string  `json:"created_at"`
	UpdatedAt       string  `json:"updated_at"`
}

// VoucherListResponse represents a list of vouchers with pagination
type VoucherListResponse struct {
	Vouchers   []VoucherResponse `json:"vouchers"`
	Pagination PaginationMeta    `json:"pagination"`
}

// PaginationMeta represents pagination metadata
type PaginationMeta struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

// ToVoucherResponse converts entity.Voucher to VoucherResponse
func ToVoucherResponse(voucher *entity.Voucher) VoucherResponse {
	return VoucherResponse{
		ID:              voucher.ID,
		VoucherCode:     voucher.VoucherCode,
		DiscountPercent: voucher.DiscountPercent,
		ExpiryDate:      voucher.ExpiryDate.Format("2006-01-02"),
		CreatedAt:       voucher.CreatedAt.Format(time.RFC3339),
		UpdatedAt:       voucher.UpdatedAt.Format(time.RFC3339),
	}
}

// ToVoucherListResponse converts a list of vouchers to VoucherListResponse
func ToVoucherListResponse(vouchers []*entity.Voucher) []VoucherResponse {
	responses := make([]VoucherResponse, len(vouchers))
	for i, voucher := range vouchers {
		responses[i] = ToVoucherResponse(voucher)
	}
	return responses
}

// BuildVoucherListResponse builds a complete voucher list response with pagination
func BuildVoucherListResponse(vouchers []*entity.Voucher, page, limit int, total int64) VoucherListResponse {
	totalPages := int(total / int64(limit))
	if total%int64(limit) > 0 {
		totalPages++
	}

	return VoucherListResponse{
		Vouchers: ToVoucherListResponse(vouchers),
		Pagination: PaginationMeta{
			Page:       page,
			Limit:      limit,
			Total:      total,
			TotalPages: totalPages,
		},
	}
}
