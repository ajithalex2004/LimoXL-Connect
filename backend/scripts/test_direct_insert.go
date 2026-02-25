package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/google/uuid"
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

	fmt.Println("=== Testing Direct INSERT into trip_invoices ===")

	tripID := uuid.MustParse("f857f1ae-bb78-4ef9-8909-6290d1c96ac3")
	supplierID := uuid.MustParse("24a8f9b1-5bdd-459e-826c-586946f176e6")

	query := `
		INSERT INTO trip_invoices (trip_id, supplier_company_id, invoice_number, amount, status, created_at)
		VALUES ($1, $2, $3, $4, $5, NOW())
	`

	_, err = db.Exec(query, tripID, supplierID, "INV-TEST-001", 300.00, "PENDING")
	if err != nil {
		fmt.Printf("❌ INSERT failed: %v\n", err)
		return
	}

	fmt.Println("✓ Invoice inserted successfully!")

	// Verify
	var count int
	db.QueryRow("SELECT COUNT(*) FROM trip_invoices").Scan(&count)
	fmt.Printf("Total invoices in table: %d\n", count)
}
