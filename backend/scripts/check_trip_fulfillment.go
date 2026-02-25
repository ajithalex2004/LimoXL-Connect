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

	fmt.Println("=== Checking Trip Details ===")

	query := `
		SELECT id, reference_no, status, fulfillment_company_id, driver_link_token
		FROM trips
		WHERE reference_no = 'TR260005'
	`

	var id, refNo, status, token string
	var fulfillmentID sql.NullString

	err = db.QueryRow(query).Scan(&id, &refNo, &status, &fulfillmentID, &token)
	if err != nil {
		log.Fatal("Query failed:", err)
	}

	fmt.Printf("Trip ID: %s\n", id)
	fmt.Printf("Reference: %s\n", refNo)
	fmt.Printf("Status: %s\n", status)
	fmt.Printf("Fulfillment Company ID: %v (Valid: %v)\n", fulfillmentID.String, fulfillmentID.Valid)
	fmt.Printf("Driver Link Token: %s\n", token)

	if !fulfillmentID.Valid {
		fmt.Println("\n❌ ISSUE: fulfillment_company_id is NULL!")
		fmt.Println("This is why the trip doesn't appear in Partner's Assigned Trips list")
	} else {
		fmt.Println("\n✓ fulfillment_company_id is set correctly")
	}
}
