package model

import "time"

type Voucher struct {
	ID          int64     `json:"id"`
	Code        string    `json:"code"`         // Mã voucher, e.g., "SUMMER2026"
	Description string    `json:"description"`  // Mô tả voucher
	Type        string    `json:"type"`         // "percentage" hoặc "fixed"
	Value       float64   `json:"value"`        // Giá trị giảm (% hoặc số tiền)
	MinOrder    float64   `json:"min_order"`    // Đơn hàng tối thiểu
	MaxDiscount float64   `json:"max_discount"` // Giảm tối đa (cho percentage)
	UsageLimit  int       `json:"usage_limit"`  // Số lần sử dụng tối đa
	UsedCount   int       `json:"used_count"`   // Đã dùng bao nhiêu lần
	StartDate   time.Time `json:"start_date"`   // Ngày bắt đầu
	EndDate     time.Time `json:"end_date"`     // Ngày kết thúc
	IsActive    bool      `json:"is_active"`    // Còn hoạt động không
	CreatedAt   time.Time `json:"created_at"`
}

type VoucherUsage struct {
	ID        int64     `json:"id"`
	VoucherID int64     `json:"voucher_id"`
	OrderID   int64     `json:"order_id"`
	UserID    int64     `json:"user_id,omitempty"`
	Discount  float64   `json:"discount"` // Số tiền đã giảm
	CreatedAt time.Time `json:"created_at"`
}
