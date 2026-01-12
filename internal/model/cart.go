package model

import "time"

type CartItem struct {
	ID          int64     `json:"id"`
	UserID      int64     `json:"user_id"`
	ProductID   int64     `json:"product_id"`
	ProductName string    `json:"product_name,omitempty"`
	Price       float64   `json:"price,omitempty"`
	ImageURL    string    `json:"image_url,omitempty"`
	Quantity    int       `json:"quantity"`
	CreatedAt   time.Time `json:"created_at"`
}
