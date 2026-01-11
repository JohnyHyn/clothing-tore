package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-sql-driver/mysql"
)

func mustEnv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatalf("Missing environment variable: %s", key)
	}
	return val
}

func Connect() *sql.DB {
	cfg := mysql.Config{
		User:                 mustEnv("DB_USER"),
		Passwd:               mustEnv("DB_PASSWORD"),
		Net:                  "tcp",
		Addr:                 fmt.Sprintf("%s:%s", mustEnv("DB_HOST"), mustEnv("DB_PORT")),
		DBName:               mustEnv("DB_NAME"),
		ParseTime:            true,
		AllowNativePasswords: true,
	}

	dsn := cfg.FormatDSN()

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}

	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err := db.Ping(); err != nil {
		log.Fatal("Cannot connect DB:", err)
	}

	log.Println("Connected to MySQL")
	return db
}
