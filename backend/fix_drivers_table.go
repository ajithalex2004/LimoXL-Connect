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

	fmt.Println("Adding assigned_vehicle_id column to drivers table...")

	// Add assigned_vehicle_id column if it doesn't exist
	_, err = db.Exec(`
		ALTER TABLE drivers 
		ADD COLUMN IF NOT EXISTS assigned_vehicle_id UUID;
	`)
	if err != nil {
		log.Fatal("Error adding assigned_vehicle_id column:", err)
	}

	fmt.Println("✓ Added assigned_vehicle_id column")
	fmt.Println("\n✓ Database schema updated successfully!")
}
