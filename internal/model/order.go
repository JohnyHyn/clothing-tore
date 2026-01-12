package model

type Order struct {
	ID             int64       `json:"id"`
	UserID         int64       `json:"user_id"`
	CustomerName   string      `json:"customer_name"`
	CustomerPhone  string      `json:"customer_phone"`
	Address        string      `json:"address"`
	TotalPrice     float64     `json:"total_price"`
	PaymentMethod  string      `json:"payment_method"`
	ShippingMethod string      `json:"shipping_method"`
	Status         string      `json:"status"`
	CreatedAt      string      `json:"created_at,omitempty"`
	Items          []OrderItem `json:"items"`
}
