package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	dsn := "clothing_app:StrongPassword@123@tcp(127.0.0.1:3306)/clothing_store"
	if os.Getenv("DB_NAME") != "" {
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
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

	password := "123456"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}

	users := []struct {
		Email string
		Role  string
	}{
		{"admin@shop.com", "admin"},
		{"staff@shop.com", "staff"},
	}

	for _, u := range users {
		_, err = db.Exec(`
			INSERT INTO users (email, password, role) 
			VALUES (?, ?, ?)
			ON DUPLICATE KEY UPDATE password = ?, role = ?
		`, u.Email, string(hashedPassword), u.Role, string(hashedPassword), u.Role)

		if err != nil {
			log.Printf("Failed to seed user %s: %v", u.Email, err)
		} else {
			fmt.Printf("User %s created/updated successfully!\n", u.Email)
		}
	}
}
