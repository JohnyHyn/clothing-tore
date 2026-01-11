package model

type Order struct {
	ID            int64       `json:"id"`
	CustomerName  string      `json:"customer_name"`
	CustomerPhone string      `json:"customer_phone"`
	TotalPrice    float64     `json:"total_price"`
	Status        string      `json:"status"`
	CreatedAt     string      `json:"created_at,omitempty"`
	Items         []OrderItem `json:"items"`
}
