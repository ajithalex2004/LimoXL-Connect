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

	fmt.Println("=== Checking for duplicate invoice numbers ===\n")

	query := `
		SELECT id, invoice_number, amount, status, created_at
		FROM trip_invoices
		ORDER BY created_at DESC
		LIMIT 10
	`

	rows, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	fmt.Println("Recent invoices:")
	count := 0
	for rows.Next() {
		var id, invoiceNum, status string
		var amount float64
		var createdAt string
		rows.Scan(&id, &invoiceNum, &amount, &status, &createdAt)
		count++
		fmt.Printf("%d. %s - Amount: %.2f - Status: %s - Created: %s\n",
			count, invoiceNum, amount, status, createdAt)
	}

	if count == 0 {
		fmt.Println("No invoices found in database")
	}

	// Check for the specific invoice number
	var exists bool
	checkQuery := `SELECT EXISTS(SELECT 1 FROM trip_invoices WHERE invoice_number = 'INV-2026-RIDER-001')`
	db.QueryRow(checkQuery).Scan(&exists)

	if exists {
		fmt.Println("\n❌ Invoice 'INV-2026-RIDER-001' already exists!")
		fmt.Println("This is why the browser submission fails.")
		fmt.Println("\nSolution: Use a different invoice number or delete the existing one.")
	} else {
		fmt.Println("\n✓ Invoice 'INV-2026-RIDER-001' does not exist")
	}
}
