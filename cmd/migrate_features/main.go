package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// Simple migration script to add vouchers, shippings, and related tables

	// Read env or use default for migration
	dsn := "clothing_app:StrongPassword@123@tcp(127.0.0.1:3306)/clothing_store?parseTime=true"
	if os.Getenv("DB_PASSWORD") != "" {
		dsn = fmt.Sprintf(
			"%s:%s@tcp(%s:%s)/%s?parseTime=true",
			os.Getenv("DB_USER"),
			os.Getenv("DB_PASSWORD"),
			os.Getenv("DB_HOST"),
			os.Getenv("DB_PORT"),
			os.Getenv("DB_NAME"),
		)
	}

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create vouchers table
	vouchersQuery := `
	CREATE TABLE IF NOT EXISTS vouchers (
		id BIGINT AUTO_INCREMENT PRIMARY KEY,
		code VARCHAR(50) UNIQUE NOT NULL,
		description TEXT,
		type VARCHAR(20) NOT NULL,
		value DECIMAL(15, 2) NOT NULL,
		min_order DECIMAL(15, 2) DEFAULT 0,
		max_discount DECIMAL(15, 2) DEFAULT 0,
		usage_limit INT DEFAULT 1,
		used_count INT DEFAULT 0,
		start_date DATETIME NOT NULL,
		end_date DATETIME NOT NULL,
		is_active BOOLEAN DEFAULT TRUE,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		INDEX idx_code (code),
		INDEX idx_active (is_active)
	);`

	_, err = db.Exec(vouchersQuery)
	if err != nil {
		log.Fatal("Failed to create vouchers table:", err)
	}
	fmt.Println("✓ vouchers table created successfully")

	// Create voucher_usage table
	voucherUsageQuery := `
	CREATE TABLE IF NOT EXISTS voucher_usage (
		id BIGINT AUTO_INCREMENT PRIMARY KEY,
		voucher_id BIGINT NOT NULL,
		order_id BIGINT NOT NULL,
		user_id BIGINT,
		discount DECIMAL(15, 2) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (voucher_id) REFERENCES vouchers(id),
		FOREIGN KEY (order_id) REFERENCES orders(id),
		INDEX idx_voucher_id (voucher_id),
		INDEX idx_order_id (order_id)
	);`

	_, err = db.Exec(voucherUsageQuery)
	if err != nil {
		log.Fatal("Failed to create voucher_usage table:", err)
	}
	fmt.Println("✓ voucher_usage table created successfully")

	// Create shippings table
	shippingsQuery := `
	CREATE TABLE IF NOT EXISTS shippings (
		id BIGINT AUTO_INCREMENT PRIMARY KEY,
		order_id BIGINT NOT NULL UNIQUE,
		method VARCHAR(50) NOT NULL,
		provider VARCHAR(50) NOT NULL,
		tracking_code VARCHAR(100),
		fee DECIMAL(15, 2) NOT NULL,
		status VARCHAR(50) DEFAULT 'pending',
		address TEXT NOT NULL,
		city VARCHAR(100) NOT NULL,
		district VARCHAR(100) NOT NULL,
		ward VARCHAR(100),
		receiver_name VARCHAR(255) NOT NULL,
		receiver_phone VARCHAR(20) NOT NULL,
		note TEXT,
		estimated_date DATETIME,
		delivered_at DATETIME,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		FOREIGN KEY (order_id) REFERENCES orders(id),
		INDEX idx_tracking_code (tracking_code),
		INDEX idx_status (status)
	);`

	_, err = db.Exec(shippingsQuery)
	if err != nil {
		log.Fatal("Failed to create shippings table:", err)
	}
	fmt.Println("✓ shippings table created successfully")

	// Create shipping_history table
	shippingHistoryQuery := `
	CREATE TABLE IF NOT EXISTS shipping_history (
		id BIGINT AUTO_INCREMENT PRIMARY KEY,
		shipping_id BIGINT NOT NULL,
		status VARCHAR(50) NOT NULL,
		location VARCHAR(255),
		note TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (shipping_id) REFERENCES shippings(id),
		INDEX idx_shipping_id (shipping_id)
	);`

	_, err = db.Exec(shippingHistoryQuery)
	if err != nil {
		log.Fatal("Failed to create shipping_history table:", err)
	}
	fmt.Println("✓ shipping_history table created successfully")

	fmt.Println("\n🎉 All tables created successfully!")
	fmt.Println("Tables created:")
	fmt.Println("  - vouchers")
	fmt.Println("  - voucher_usage")
	fmt.Println("  - shippings")
	fmt.Println("  - shipping_history")
}
