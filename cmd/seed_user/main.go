package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	// Connect to DB directly for seeding
	dsn := "clothing_app:StrongPassword@123@tcp(127.0.0.1:3306)/clothing_store"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	password := "123456"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}

	// Check columns
	rows, err := db.Query("SELECT * FROM users LIMIT 0")
	if err != nil {
		log.Fatal(err)
	}
	cols, _ := rows.Columns()
	fmt.Println("Columns:", cols)

	_, err = db.Exec(`
		INSERT INTO users (email, password, role) 
		VALUES (?, ?, ?)
		ON DUPLICATE KEY UPDATE password = ?
	`, "admin@shop.com", string(hashedPassword), "admin", string(hashedPassword))

	if err != nil {
		log.Fatal("Failed to seed user:", err)
	}

	fmt.Println("User admin@shop.com / 123456 created/updated successfully!")
}
