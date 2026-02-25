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

	fmt.Println("Adding missing columns to vehicles table...")

	// Add vehicle_class column if it doesn't exist
	_, err = db.Exec(`
		ALTER TABLE vehicles 
		ADD COLUMN IF NOT EXISTS vehicle_class TEXT;
	`)
	if err != nil {
		log.Fatal("Error adding vehicle_class column:", err)
	}
	fmt.Println("✓ Added vehicle_class column")

	// Add vehicle_group column if it doesn't exist
	_, err = db.Exec(`
		ALTER TABLE vehicles 
		ADD COLUMN IF NOT EXISTS vehicle_group TEXT;
	`)
	if err != nil {
		log.Fatal("Error adding vehicle_group column:", err)
	}
	fmt.Println("✓ Added vehicle_group column")

	// Add model column if it doesn't exist
	_, err = db.Exec(`
		ALTER TABLE vehicles 
		ADD COLUMN IF NOT EXISTS model TEXT;
	`)
	if err != nil {
		log.Fatal("Error adding model column:", err)
	}
	fmt.Println("✓ Added model column")

	// Rename plate_number to license_plate if needed
	_, err = db.Exec(`
		ALTER TABLE vehicles 
		RENAME COLUMN plate_number TO license_plate;
	`)
	if err != nil {
		// Column might already be named license_plate or might not exist
		fmt.Println("Note: plate_number column already renamed or doesn't exist")
	} else {
		fmt.Println("✓ Renamed plate_number to license_plate")
	}

	fmt.Println("\n✓ Database schema updated successfully!")
}
