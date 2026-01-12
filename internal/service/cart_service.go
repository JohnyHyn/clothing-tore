package service

import (
	"clothing-store/internal/model"
	"database/sql"
)

type CartService struct {
	DB *sql.DB
}

func (s *CartService) AddToCart(userID, productID int64, quantity int) error {
	// Check if item already exists
	var existingQuantity int
	err := s.DB.QueryRow("SELECT quantity FROM cart_items WHERE user_id = ? AND product_id = ?", userID, productID).Scan(&existingQuantity)

	if err == sql.ErrNoRows {
		_, err = s.DB.Exec("INSERT INTO cart_items (user_id, product_id, quantity) VALUES (?, ?, ?)", userID, productID, quantity)
	} else if err == nil {
		_, err = s.DB.Exec("UPDATE cart_items SET quantity = quantity + ? WHERE user_id = ? AND product_id = ?", quantity, userID, productID)
	}

	return err
}

func (s *CartService) GetCart(userID int64) ([]model.CartItem, error) {
	query := `
		SELECT c.id, c.product_id, p.name, p.price, p.image_url, c.quantity, c.created_at
		FROM cart_items c
		JOIN products p ON c.product_id = p.id
		WHERE c.user_id = ?
	`
	rows, err := s.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []model.CartItem
	for rows.Next() {
		var item model.CartItem
		var imgURL sql.NullString
		err := rows.Scan(&item.ID, &item.ProductID, &item.ProductName, &item.Price, &imgURL, &item.Quantity, &item.CreatedAt)
		if err != nil {
			return nil, err
		}
		item.ImageURL = imgURL.String
		items = append(items, item)
	}
	return items, nil
}

func (s *CartService) UpdateQuantity(userID, cartItemID int64, quantity int) error {
	_, err := s.DB.Exec("UPDATE cart_items SET quantity = ? WHERE id = ? AND user_id = ?", quantity, cartItemID, userID)
	return err
}

func (s *CartService) RemoveItem(userID, cartItemID int64) error {
	_, err := s.DB.Exec("DELETE FROM cart_items WHERE id = ? AND user_id = ?", cartItemID, userID)
	return err
}

func (s *CartService) ClearCart(userID int64) error {
	_, err := s.DB.Exec("DELETE FROM cart_items WHERE user_id = ?", userID)
	return err
}
