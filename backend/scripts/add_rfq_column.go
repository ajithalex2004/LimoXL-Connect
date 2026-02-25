package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func main() {
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		connStr = "postgres://postgres:password@localhost:5432/limoxlink?sslmode=disable"
	}

	fmt.Printf("Connecting to: %s\n\n", connStr)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	fmt.Println("Adding rfq_number column to trips table...")

	// Add rfq_number column if it doesn't exist
	_, err = db.Exec(`
		ALTER TABLE trips 
		ADD COLUMN IF NOT EXISTS rfq_number TEXT;
	`)
	if err != nil {
		log.Fatal("Error adding rfq_number column:", err)
	}

	fmt.Println("✓ Added rfq_number column")

	// Create sequence if it doesn't exist
	_, err = db.Exec(`
		CREATE SEQUENCE IF NOT EXISTS rfq_id_seq START 1000;
	`)
	if err != nil {
		log.Fatal("Error creating sequence:", err)
	}

	fmt.Println("✓ Created rfq_id_seq sequence")
	fmt.Println("\n✓ Database schema updated successfully!")
}
