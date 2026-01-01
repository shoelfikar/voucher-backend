package request

// CreateVoucherRequest represents the request to create a new voucher
type CreateVoucherRequest struct {
	VoucherCode     string  `json:"voucher_code" binding:"required,max=50"`
	DiscountPercent float64 `json:"discount_percent" binding:"required,min=1,max=100"`
	ExpiryDate      string  `json:"expiry_date" binding:"required"`
}

// UpdateVoucherRequest represents the request to update an existing voucher
type UpdateVoucherRequest struct {
	VoucherCode     string  `json:"voucher_code" binding:"required,max=50"`
	DiscountPercent float64 `json:"discount_percent" binding:"required,min=1,max=100"`
	ExpiryDate      string  `json:"expiry_date" binding:"required"`
}

// BatchUploadRequest represents the request to upload a batch of vouchers
type BatchUploadRequest struct {
	Vouchers []CreateVoucherRequest `json:"vouchers" binding:"required"`
}
