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

	fmt.Println("=== Inspection for TR260001 ===")

	// 1. Check Trip Details
	var id, status, visibility, rfqNo string
	err = db.QueryRow(`
		SELECT id, status, visibility, rfq_number 
		FROM trips 
		WHERE reference_no = 'TR260001'
	`).Scan(&id, &status, &visibility, &rfqNo)

	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("Error: Trip TR260001 not found!")
			return
		}
		log.Fatal(err)
	}

	fmt.Printf("Trip ID: %s\n", id)
	fmt.Printf("Status: %s (Expected: MARKETPLACE_SEARCH or OFFERED)\n", status)
	fmt.Printf("Visibility: %s\n", visibility)
	fmt.Printf("RFQ Number: %s\n", rfqNo)

	// 2. Check Trip Access
	fmt.Println("\n=== Trip Access Records ===")
	rows, err := db.Query(`
		SELECT ta.company_id, c.name 
		FROM trip_access ta
		LEFT JOIN companies c ON ta.company_id = c.id
		WHERE ta.trip_id = $1
	`, id)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		var cID, cName string
		rows.Scan(&cID, &cName)
		fmt.Printf("- Granted to: %s (ID: %s)\n", cName, cID)
		count++
	}

	if count == 0 {
		fmt.Println("No specific access records found.")
		if visibility == "DIRECT" {
			fmt.Println("WARNING: Visibility is DIRECT but no access records exist!")
		}
	}

	// 3. User Check (Rider Arabia)
	fmt.Println("\n=== Rider Arabia Check ===")
	rows, err = db.Query(`
		SELECT c.id, c.name, u.email 
		FROM companies c
		LEFT JOIN users u ON c.id = u.company_id
		WHERE c.name LIKE 'Rider Arabia%'
	`)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var cID, cName string
		var uEmail sql.NullString
		rows.Scan(&cID, &cName, &uEmail)
		fmt.Printf("Company: %s (ID: %s) - User: %s\n", cName, cID, uEmail.String)
	}
}
