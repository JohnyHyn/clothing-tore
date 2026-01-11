package model

import "time"

type Payment struct {
	ID            int64     `json:"id"`
	OrderID       int64     `json:"order_id"`
	TransactionID string    `json:"transaction_id"`
	Provider      string    `json:"provider"` // e.g., "momo", "vnpay", "manual"
	Amount        float64   `json:"amount"`
	Status        string    `json:"status"` // "pending", "success", "failed"
	CreatedAt     time.Time `json:"created_at"`
}
