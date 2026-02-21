package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func main() {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL not set")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	refs := []string{"RFQ-1005", "RFQ-1006"}

	for _, ref := range refs {
		fmt.Printf("\n--- Debugging %s ---\n", ref)
		var tripID, reqCompID string
		err := db.QueryRow("SELECT id, requesting_company_id FROM trips WHERE reference_no = $1", ref).Scan(&tripID, &reqCompID)
		if err != nil {
			fmt.Printf("Trip not found: %v\n", err)
			continue
		}
		fmt.Printf("Trip ID: %s\n", tripID)
		fmt.Printf("Requesting Company ID: %s\n", reqCompID)

		// Check Offers
		rows, err := db.Query("SELECT id, supplier_company_id, status FROM trip_offers WHERE trip_id = $1", tripID)
		if err != nil {
			fmt.Printf("Error querying offers: %v\n", err)
			continue
		}
		defer rows.Close()

		offerCount := 0
		for rows.Next() {
			var offID, supID, status string
			rows.Scan(&offID, &supID, &status)
			fmt.Printf("  Offer Found: ID=%s, Supplier=%s, Status=%s\n", offID, supID, status)
			offerCount++
		}
		if offerCount == 0 {
			fmt.Println("  No offers found for this trip.")
		}
	}

	// Check Supplier Company
	supID := "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11"
	var compName string
	err = db.QueryRow("SELECT name FROM companies WHERE id = $1", supID).Scan(&compName)
	if err != nil {
		fmt.Printf("\nCRITICAL: Supplier Company %s NOT FOUND in companies table! %v\n", supID, err)
	} else {
		fmt.Printf("\nSupplier Company Found: %s\n", compName)
	}

	// Also check generic offers to see if any exist
	fmt.Println("\n--- Checking All Offers ---")
	rows, err := db.Query("SELECT id, trip_id, status FROM trip_offers LIMIT 5")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		defer rows.Close()
		for rows.Next() {
			var id, tid, stat string
			rows.Scan(&id, &tid, &stat)
			fmt.Printf("Offer: %s, Trip: %s, Status: %s\n", id, tid, stat)
		}
	}
}
