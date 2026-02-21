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

	// Rider Arabia's ID
	partnerID := uuid.MustParse("24a8f9b1-5bdd-459e-826c-586946f176e6")

	fmt.Println("=== Testing ListUninvoicedTrips Query ===")

	query := `
		SELECT t.id, t.reference_no, t.pickup_zone, t.dropoff_zone, t.pickup_time, t.status, t.price
		FROM trips t
		LEFT JOIN invoices i ON t.id = i.trip_id
		WHERE t.fulfillment_company_id = $1 
		AND t.status = 'COMPLETED'
		AND i.id IS NULL
		ORDER BY t.pickup_time DESC
	`

	rows, err := db.Query(query, partnerID)
	if err != nil {
		fmt.Printf("❌ Query failed: %v\n", err)
		return
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		var id uuid.UUID
		var refNo string
		var pZone, dZone sql.NullString
		var pickupTime string
		var status string
		var price sql.NullFloat64

		err := rows.Scan(&id, &refNo, &pZone, &dZone, &pickupTime, &status, &price)
		if err != nil {
			fmt.Printf("❌ Scan failed: %v\n", err)
			return
		}
		count++
		fmt.Printf("%d. Ref: %s, Status: %s, Price: %.2f\n", count, refNo, status, price.Float64)
	}

	if count == 0 {
		fmt.Println("\n✓ Query executed successfully but returned 0 trips")
		fmt.Println("Reason: No trips with status='COMPLETED' and fulfillment_company_id set")
	} else {
		fmt.Printf("\n✓ Found %d uninvoiced completed trips\n", count)
	}
}
