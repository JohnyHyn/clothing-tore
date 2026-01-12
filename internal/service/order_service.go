package service

import (
	"clothing-store/internal/model"
	"database/sql"
	"errors"
	"strings"
)

type OrderService struct {
	DB *sql.DB
}

func (s *OrderService) CreateOrder(order *model.Order) error {
	tx, err := s.DB.Begin()
	if err != nil {
		return err
	}

	var total float64 = 0

	for i := range order.Items {
		item := &order.Items[i]

		var stock int
		var price float64

		err := tx.QueryRow(
			"SELECT stock, price FROM products WHERE id = ? AND status='active'",
			item.ProductID,
		).Scan(&stock, &price)

		if err != nil {
			tx.Rollback()
			return errors.New("product not found")
		}

		if stock < item.Quantity {
			tx.Rollback()
			return errors.New("not enough stock")
		}

		item.Price = price
		total += price * float64(item.Quantity)

		_, err = tx.Exec(
			"UPDATE products SET stock = stock - ? WHERE id = ?",
			item.Quantity,
			item.ProductID,
		)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	result, err := tx.Exec(
		"INSERT INTO orders (user_id, customer_name, customer_phone, address, total_price, payment_method, shipping_method) VALUES (?, ?, ?, ?, ?, ?, ?)",
		order.UserID,
		order.CustomerName,
		order.CustomerPhone,
		order.Address,
		total,
		order.PaymentMethod,
		order.ShippingMethod,
	)
	if err != nil {
		tx.Rollback()
		return err
	}

	orderID, _ := result.LastInsertId()

	for _, item := range order.Items {
		_, err := tx.Exec(
			"INSERT INTO order_items (order_id, product_id, quantity, price) VALUES (?, ?, ?, ?)",
			orderID,
			item.ProductID,
			item.Quantity,
			item.Price,
		)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()
	order.ID = orderID
	order.TotalPrice = total
	return nil
}
func (s *OrderService) CancelOrder(orderID int64) error {
	tx, err := s.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var status string
	err = tx.QueryRow("SELECT status FROM orders WHERE id = ?", orderID).Scan(&status)
	if err != nil {
		return errors.New("order not found")
	}

	if status != "pending" {
		return errors.New("only pending orders can be cancelled")
	}

	rows, err := tx.Query("SELECT product_id, quantity FROM order_items WHERE order_id = ?", orderID)
	if err != nil {
		return err
	}
	defer rows.Close()

	type item struct {
		pID int64
		qty int
	}
	var items []item
	for rows.Next() {
		var i item
		rows.Scan(&i.pID, &i.qty)
		items = append(items, i)
	}
	rows.Close()

	for _, i := range items {
		_, err = tx.Exec("UPDATE products SET stock = stock + ? WHERE id = ?", i.qty, i.pID)
		if err != nil {
			return err
		}
	}

	_, err = tx.Exec("UPDATE orders SET status = 'cancelled' WHERE id = ?", orderID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (s *OrderService) DeleteOrder(orderID int64) error {
	tx, err := s.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var status string
	err = tx.QueryRow("SELECT status FROM orders WHERE id = ? FOR UPDATE", orderID).Scan(&status)
	if err != nil {
		return err
	}

	if status != "pending" && status != "cancelled" {
		return errors.New("cannot delete order that is already paid or processed")
	}

	// Restore stock if it was pending (already subtracted)
	if status == "pending" {
		rows, err := tx.Query("SELECT product_id, quantity FROM order_items WHERE order_id = ?", orderID)
		if err != nil {
			return err
		}
		defer rows.Close()

		type item struct {
			pID int64
			qty int
		}
		var items []item
		for rows.Next() {
			var i item
			rows.Scan(&i.pID, &i.qty)
			items = append(items, i)
		}
		rows.Close()

		for _, i := range items {
			_, err = tx.Exec("UPDATE products SET stock = stock + ? WHERE id = ?", i.qty, i.pID)
			if err != nil {
				return err
			}
		}
	}

	// Delete items first
	_, err = tx.Exec("DELETE FROM order_items WHERE order_id = ?", orderID)
	if err != nil {
		return err
	}

	// Delete order
	_, err = tx.Exec("DELETE FROM orders WHERE id = ?", orderID)
	if err != nil {
		return err
	}

	return tx.Commit()
}
func (s *OrderService) PayOrder(orderID int64) error {
	tx, err := s.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var status string
	err = tx.QueryRow(
		"SELECT status FROM orders WHERE id = ? FOR UPDATE",
		orderID,
	).Scan(&status)
	if err != nil {
		return err
	}

	if status != "pending" {
		return errors.New("only pending orders can be paid")
	}

	_, err = tx.Exec(
		"UPDATE orders SET status = 'paid' WHERE id = ?",
		orderID,
	)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (s *OrderService) ApproveOrder(orderID int64) error {
	tx, err := s.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var status string
	err = tx.QueryRow("SELECT status FROM orders WHERE id = ? FOR UPDATE", orderID).Scan(&status)
	if err != nil {
		return err
	}

	if status != "chờ xử lý" {
		return errors.New("chỉ có thể duyệt đơn hàng đang chờ xử lý")
	}

	_, err = tx.Exec("UPDATE orders SET status = 'đã xử lý' WHERE id = ?", orderID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (s *OrderService) GetByID(orderID int64) (*model.Order, error) {
	order := &model.Order{}
	query := "SELECT id, customer_name, customer_phone, total_price, status FROM orders WHERE id = ?"
	err := s.DB.QueryRow(query, orderID).Scan(
		&order.ID, &order.CustomerName, &order.CustomerPhone, &order.TotalPrice, &order.Status,
	)
	if err != nil {
		return nil, err
	}

	rows, err := s.DB.Query("SELECT product_id, quantity, price FROM order_items WHERE order_id = ?", orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var item model.OrderItem
		if err := rows.Scan(&item.ProductID, &item.Quantity, &item.Price); err != nil {
			return nil, err
		}
		order.Items = append(order.Items, item)
	}

	return order, nil
}
func (s *OrderService) GetOrderByID(orderID int64) (*model.Order, error) {
	order := &model.Order{}

	err := s.DB.QueryRow(`
		SELECT id, user_id, customer_name, customer_phone, address, total_price, payment_method, shipping_method, status
		FROM orders
		WHERE id = ?
	`, orderID).Scan(
		&order.ID,
		&order.UserID,
		&order.CustomerName,
		&order.CustomerPhone,
		&order.Address,
		&order.TotalPrice,
		&order.PaymentMethod,
		&order.ShippingMethod,
		&order.Status,
	)

	if err != nil {
		return nil, err
	}

	rows, err := s.DB.Query(`
		SELECT oi.product_id, p.name, oi.price, oi.quantity
		FROM order_items oi
		JOIN products p ON oi.product_id = p.id
		WHERE oi.order_id = ?
	`, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var item model.OrderItem
		if err := rows.Scan(
			&item.ProductID,
			&item.ProductName,
			&item.Price,
			&item.Quantity,
		); err != nil {
			return nil, err
		}
		order.Items = append(order.Items, item)
	}

	return order, nil
}
func (s *OrderService) RefundOrder(orderID int64, reason string) error {
	tx, err := s.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var status string
	var total float64

	err = tx.QueryRow(`
		SELECT status, total_price
		FROM orders
		WHERE id = ?
		FOR UPDATE
	`, orderID).Scan(&status, &total)
	if err != nil {
		return err
	}

	if status != "paid" {
		return errors.New("only paid orders can be refunded")
	}

	// restore stock
	rows, err := tx.Query(`
		SELECT product_id, quantity
		FROM order_items
		WHERE order_id = ?
	`, orderID)
	if err != nil {
		return err
	}
	defer rows.Close()

	type orderItem struct {
		ProductID int64
		Quantity  int
	}
	var items []orderItem

	for rows.Next() {
		var item orderItem
		if err := rows.Scan(&item.ProductID, &item.Quantity); err != nil {
			return err
		}
		items = append(items, item)
	}
	rows.Close()

	for _, item := range items {
		_, err := tx.Exec(`
			UPDATE products
			SET stock = stock + ?
			WHERE id = ?
		`, item.Quantity, item.ProductID)
		if err != nil {
			return err
		}
	}

	// insert refund record
	_, err = tx.Exec(`
		INSERT INTO refunds (order_id, amount, reason)
		VALUES (?, ?, ?)
	`, orderID, total, reason)
	if err != nil {
		return err
	}

	// update order status
	_, err = tx.Exec(`
		UPDATE orders
		SET status = 'cancelled'
		WHERE id = ?
	`, orderID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// ListOrders retrieves a paginated list of orders with optional filters
func (s *OrderService) ListOrders(page, limit int, status, dateFrom, dateTo string) ([]model.Order, int, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit

	// Build WHERE clause dynamically
	var whereClauses []string
	var args []interface{}

	if status != "" {
		whereClauses = append(whereClauses, "status = ?")
		args = append(args, status)
	}

	if dateFrom != "" {
		whereClauses = append(whereClauses, "created_at >= ?")
		args = append(args, dateFrom)
	}

	if dateTo != "" {
		whereClauses = append(whereClauses, "created_at <= ?")
		args = append(args, dateTo)
	}

	whereClause := ""
	if len(whereClauses) > 0 {
		whereClause = "WHERE " + strings.Join(whereClauses, " AND ")
	}

	// Get total count
	var total int
	countQuery := "SELECT COUNT(*) FROM orders " + whereClause
	err := s.DB.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Get paginated orders
	query := `
		SELECT id, user_id, customer_name, customer_phone, address, total_price, payment_method, shipping_method, status, created_at
		FROM orders
		` + whereClause + `
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`

	// Add limit and offset to args
	queryArgs := append(args, limit, offset)

	rows, err := s.DB.Query(query, queryArgs...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var orders []model.Order
	for rows.Next() {
		var order model.Order
		err := rows.Scan(
			&order.ID,
			&order.UserID,
			&order.CustomerName,
			&order.CustomerPhone,
			&order.Address,
			&order.TotalPrice,
			&order.PaymentMethod,
			&order.ShippingMethod,
			&order.Status,
			&order.CreatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		orders = append(orders, order)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}

func (s *OrderService) ListUserOrders(userID int64, page, limit int) ([]model.Order, int, error) {
	offset := (page - 1) * limit

	var total int
	err := s.DB.QueryRow("SELECT COUNT(*) FROM orders WHERE user_id = ?", userID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query := `
		SELECT id, user_id, customer_name, customer_phone, address, total_price, payment_method, shipping_method, status, created_at
		FROM orders
		WHERE user_id = ?
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`
	rows, err := s.DB.Query(query, userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var orders []model.Order
	for rows.Next() {
		var order model.Order
		err := rows.Scan(
			&order.ID,
			&order.UserID,
			&order.CustomerName,
			&order.CustomerPhone,
			&order.Address,
			&order.TotalPrice,
			&order.PaymentMethod,
			&order.ShippingMethod,
			&order.Status,
			&order.CreatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		orders = append(orders, order)
	}

	return orders, total, nil
}
