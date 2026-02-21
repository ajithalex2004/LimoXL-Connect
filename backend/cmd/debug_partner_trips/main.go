package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func main() {
	// Database connection string
	connStr := "postgresql://neondb_owner:npg_7ndWFKRYEOt6@ep-calm-heart-a15voo2a-pooler.ap-southeast-1.aws.neon.tech/neondb?sslmode=require"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rfqs := []string{"RFQ-1005", "RFQ-1006"}
	partnerID := "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11"

	fmt.Printf("Checking visibility for Partner: %s\n", partnerID)

	for _, ref := range rfqs {
		fmt.Printf("\n--- Inspecting %s ---\n", ref)
		var id, status, visibility, reqCompID string
		// Check main trip details
		err := db.QueryRow(`
            SELECT id, status, visibility, requesting_company_id 
            FROM trips WHERE reference_no = $1`, ref).Scan(&id, &status, &visibility, &reqCompID)
		if err != nil {
			fmt.Printf("Trip %s not found: %v\n", ref, err)
			continue
		}
		fmt.Printf("ID: %s\n", id)
		fmt.Printf("Status: %s\n", status)
		fmt.Printf("Visibility: %s\n", visibility)
		fmt.Printf("Requesting Company: %s\n", reqCompID)

		// Check Trip Access if visibility is DIRECT
		if visibility == "DIRECT" || visibility == "PUBLIC" {
			rows, err := db.Query("SELECT company_id FROM trip_access WHERE trip_id = $1", id)
			if err != nil {
				fmt.Printf("Error checking trip_access: %v\n", err)
			} else {
				fmt.Println("Trip Access List:")
				count := 0
				for rows.Next() {
					var accessCompID string
					rows.Scan(&accessCompID)
					match := ""
					if accessCompID == partnerID {
						match = " (MATCHES TARGET PARTNER)"
					}
					fmt.Printf("  - %s%s\n", accessCompID, match)
					count++
				}
				if count == 0 {
					fmt.Println("  (No specific access entries found)")
				}
			}
		}
	}
}
