package service

import (
	"clothing-store/internal/model"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type PaymentService struct {
	DB *sql.DB
}

// PaymentProvider is an interface that real providers (MoMo, ZaloPay) would implement
type PaymentProvider interface {
	GenerateURL(orderID int64, amount float64) (string, error)
	VerifyCallback(data map[string]interface{}) (bool, error)
}

// MockProvider for demonstration
type MockProvider struct{}

func (m *MockProvider) GenerateURL(orderID int64, amount float64) (string, error) {
	// In reality, this calls the provider's API
	return fmt.Sprintf("https://mock-payment-gateway.com/pay?orderId=%d&amount=%.2f", orderID, amount), nil
}

func (s *PaymentService) CreatePayment(orderID int64, providerName string) (string, error) {
	// 1. Get Order to check amount and status
	var amount float64
	var status string
	err := s.DB.QueryRow("SELECT total_price, status FROM orders WHERE id = ?", orderID).Scan(&amount, &status)
	if err != nil {
		return "", err
	}

	if status == "paid" {
		return "", errors.New("order is already paid")
	}

	// 2. Generate Payment Record
	tx, err := s.DB.Begin()
	if err != nil {
		return "", err
	}
	defer tx.Rollback()

	// Create a dummy transaction ID for now
	timestamp := time.Now().UnixNano()
	transID := fmt.Sprintf("TRANS_%d_%d", orderID, timestamp)

	_, err = tx.Exec(`
		INSERT INTO payments (order_id, transaction_id, provider, amount, status)
		VALUES (?, ?, ?, ?, 'pending')
	`, orderID, transID, providerName, amount)
	if err != nil {
		return "", err
	}

	if err := tx.Commit(); err != nil {
		return "", err
	}

	// 3. Call Provider to get URL
	// Simple switch for now, ideally use a factory map
	if providerName == "mock" {
		provider := &MockProvider{}
		return provider.GenerateURL(orderID, amount)
	}

	return "", errors.New("unsupported provider")
}

// ConfirmPayment is called by Webhook/IPN
func (s *PaymentService) ConfirmPayment(orderID int64, transactionID string, status string) error {
	tx, err := s.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Update Payment
	res, err := tx.Exec("UPDATE payments SET status = ? WHERE transaction_id = ?", status, transactionID)
	if err != nil {
		return err
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		return errors.New("transaction not found")
	}

	// Update Order if success
	if status == "success" {
		_, err = tx.Exec("UPDATE orders SET status = 'paid' WHERE id = ?", orderID)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (s *PaymentService) GetHistory(orderID int64) ([]model.Payment, error) {
	rows, err := s.DB.Query("SELECT id, order_id, transaction_id, provider, amount, status, created_at FROM payments WHERE order_id = ?", orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var payments []model.Payment
	for rows.Next() {
		var p model.Payment
		if err := rows.Scan(&p.ID, &p.OrderID, &p.TransactionID, &p.Provider, &p.Amount, &p.Status, &p.CreatedAt); err != nil {
			return nil, err
		}
		payments = append(payments, p)
	}
	return payments, nil
}
