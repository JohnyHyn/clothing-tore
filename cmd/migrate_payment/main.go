package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// Simple migration script to add payments table

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

	query := `
	CREATE TABLE IF NOT EXISTS payments (
		id BIGINT AUTO_INCREMENT PRIMARY KEY,
		order_id BIGINT NOT NULL,
		transaction_id VARCHAR(255),
		provider VARCHAR(50) NOT NULL,
		amount DECIMAL(15, 2) NOT NULL,
		status VARCHAR(50) DEFAULT 'pending',
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (order_id) REFERENCES orders(id)
	);`

	_, err = db.Exec(query)
	if err != nil {
		log.Fatal("Migration failed:", err)
	}

	fmt.Println("Migration successful: payments table created.")
}
