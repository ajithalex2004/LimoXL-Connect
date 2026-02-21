package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() error {
	// Default connection string - CHANGE THIS or use DATABASE_URL env var
	connStr := "postgres://postgres:password@localhost:5432/limoxlink?sslmode=disable"
	if envStr := os.Getenv("DATABASE_URL"); envStr != "" {
		connStr = envStr
	}
	log.Printf("Attempting to connect to DB: %s", connStr)
	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Printf("Error opening database: %v", err)
		return fmt.Errorf("error opening database: %w", err)
	}

	log.Println("Pinging database...")
	// Make sure we can actually ping the database
	if err = DB.Ping(); err != nil {
		log.Printf("Error pinging database: %v", err)
		return fmt.Errorf("error connecting to database: %w", err)
	}

	log.Println("Successfully connected to the database")
	return nil
}

func RunMigrations() error {
	// Simple migration runner for demo
	// In prod, use golang-migrate or similar
	schema, err := os.ReadFile("db/migrations/01_init_schema.sql")
	if err != nil {
		return fmt.Errorf("error reading schema file: %w", err)
	}

	_, err = DB.Exec(string(schema))
	if err != nil {
		return fmt.Errorf("error executing schema: %w", err)
	}
	log.Println("Database schema applied successfully")
	return nil
}
