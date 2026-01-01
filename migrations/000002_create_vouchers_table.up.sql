CREATE TABLE vouchers (
    id BIGSERIAL PRIMARY KEY,
    voucher_code VARCHAR(50) UNIQUE NOT NULL,
    discount_percent DECIMAL(5,2) NOT NULL CHECK (discount_percent >= 1 AND discount_percent <= 100),
    expiry_date DATE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);

CREATE INDEX idx_voucher_code ON vouchers(voucher_code);
CREATE INDEX idx_expiry_date ON vouchers(expiry_date);
CREATE INDEX idx_deleted_at ON vouchers(deleted_at);
