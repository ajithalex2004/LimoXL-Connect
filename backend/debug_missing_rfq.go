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

	fmt.Println("--- Recent Trips (Last 10) ---")
	// Removed company_id, added fulfillment_company_id if needed, but for now just basic fields
	rows, err := db.Query(`
		SELECT id, reference_no, status, visibility, created_at, vehicle_type_requested
		FROM trips 
		ORDER BY created_at DESC 
		LIMIT 10
	`)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var id, ref, status, visibility, createdAt string
		var vehicleType sql.NullString
		rows.Scan(&id, &ref, &status, &visibility, &createdAt, &vehicleType)
		fmt.Printf("Ref: %s | Status: %s | Vis: %s | Created: %s | Veh: %s\n",
			ref, status, visibility, createdAt, vehicleType.String)
	}

	fmt.Println("\n--- Rejected Offers for Rider Arabia ---")
	// Rider Arabia ID: 24a8f9b1-5bdd-459e-826c-586946f176e6
	rows, err = db.Query(`
		SELECT trip_id, status 
		FROM trip_offers 
		WHERE supplier_company_id = '24a8f9b1-5bdd-459e-826c-586946f176e6'
	`)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		var tid, stat string
		rows.Scan(&tid, &stat)
		fmt.Printf("TripID: %s | Status: %s\n", tid, stat)
		count++
	}
	if count == 0 {
		fmt.Println("No offers found for Rider Arabia.")
	}
}
