package service_test

import (
	"clothing-store/internal/model"
	"clothing-store/internal/service"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestCreateOrder_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	s := &service.OrderService{DB: db}

	order := &model.Order{
		CustomerName:  "Test Customer",
		CustomerPhone: "1234567890",
		Items: []model.OrderItem{
			{ProductID: 1, Quantity: 2},
		},
	}

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT stock, price FROM products WHERE id = \\? AND status='active'").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"stock", "price"}).AddRow(10, 100.0))

	mock.ExpectExec("UPDATE products SET stock = stock - \\? WHERE id = \\?").
		WithArgs(2, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectExec("INSERT INTO orders").
		WithArgs("Test Customer", "1234567890", 200.0).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectExec("INSERT INTO order_items").
		WithArgs(1, 1, 2, 100.0).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	err = s.CreateOrder(order)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), order.ID)
	assert.Equal(t, 200.0, order.TotalPrice)
}
