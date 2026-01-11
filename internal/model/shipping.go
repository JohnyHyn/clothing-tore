package model

import "time"

type Shipping struct {
	ID            int64     `json:"id"`
	OrderID       int64     `json:"order_id"`
	Method        string    `json:"method"`                   // "standard", "express", "same_day"
	Provider      string    `json:"provider"`                 // "ghn", "ghtk", "viettel_post"
	TrackingCode  string    `json:"tracking_code"`            // Mã vận đơn
	Fee           float64   `json:"fee"`                      // Phí ship
	Status        string    `json:"status"`                   // "pending", "picked_up", "in_transit", "delivered", "failed"
	Address       string    `json:"address"`                  // Địa chỉ giao hàng
	City          string    `json:"city"`                     // Thành phố
	District      string    `json:"district"`                 // Quận/Huyện
	Ward          string    `json:"ward,omitempty"`           // Phường/Xã
	ReceiverName  string    `json:"receiver_name"`            // Tên người nhận
	ReceiverPhone string    `json:"receiver_phone"`           // SĐT người nhận
	Note          string    `json:"note,omitempty"`           // Ghi chú
	EstimatedDate time.Time `json:"estimated_date,omitempty"` // Ngày giao dự kiến
	DeliveredAt   time.Time `json:"delivered_at,omitempty"`   // Ngày giao thực tế
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type ShippingHistory struct {
	ID         int64     `json:"id"`
	ShippingID int64     `json:"shipping_id"`
	Status     string    `json:"status"`
	Location   string    `json:"location,omitempty"` // Vị trí hiện tại
	Note       string    `json:"note,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
}
