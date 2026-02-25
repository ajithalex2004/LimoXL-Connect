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
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	fmt.Println("=== Creating trip_invoices table ===")

	createTable := `
		CREATE TABLE IF NOT EXISTS trip_invoices (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			trip_id UUID NOT NULL REFERENCES trips(id),
			supplier_company_id UUID NOT NULL,
			invoice_number VARCHAR(100) NOT NULL,
			amount DECIMAL(10, 2) NOT NULL,
			status VARCHAR(50) NOT NULL DEFAULT 'PENDING',
			created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			UNIQUE(invoice_number)
		)
	`

	_, err = db.Exec(createTable)
	if err != nil {
		log.Fatal("Failed to create table:", err)
	}

	fmt.Println("✓ trip_invoices table created successfully")

	// Create index for faster lookups
	createIndex := `CREATE INDEX IF NOT EXISTS idx_trip_invoices_supplier ON trip_invoices(supplier_company_id)`
	_, err = db.Exec(createIndex)
	if err != nil {
		log.Printf("Warning: Failed to create index: %v\n", err)
	} else {
		fmt.Println("✓ Index created on supplier_company_id")
	}

	createIndex2 := `CREATE INDEX IF NOT EXISTS idx_trip_invoices_trip ON trip_invoices(trip_id)`
	_, err = db.Exec(createIndex2)
	if err != nil {
		log.Printf("Warning: Failed to create index: %v\n", err)
	} else {
		fmt.Println("✓ Index created on trip_id")
	}
}
