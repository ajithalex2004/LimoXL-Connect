package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

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

	// Create 5 fresh RFQs
	for i := 1; i <= 5; i++ {
		tripID := uuid.New()
		ref := fmt.Sprintf("RFQ-NEW-%d", time.Now().UnixNano()+int64(i))
		pickupTime := time.Now().Add(time.Duration(i*24) * time.Hour) // 1, 2, 3... days from now

		_, err = db.Exec(`
			INSERT INTO trips (
				id, reference_no, status, visibility, created_at, updated_at, 
				pickup_zone, dropoff_zone, pickup_time, passenger_name, passenger_phone,
				requesting_company_id, vehicle_type_requested, price,
				pickup_location, dropoff_location
			) VALUES (
				$1, $2, 'MARKETPLACE_SEARCH', 'PUBLIC', NOW(), NOW(),
				$3, $4, $5, $6, '+971500000000',
				$7, 'Sedan', 0,
				ST_SetSRID(ST_MakePoint(55.27, 25.20), 4326),
				ST_SetSRID(ST_MakePoint(55.36, 25.25), 4326)
			)`,
			tripID, ref,
			fmt.Sprintf("Zone A-%d", i), fmt.Sprintf("Zone B-%d", i),
			pickupTime, fmt.Sprintf("Passenger %d", i),
			"a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a00") // Operator ID

		if err != nil {
			log.Printf("Failed to create trip %d: %v", i, err)
		} else {
			fmt.Printf("✓ Created Trip: %s (Ref: %s)\n", tripID, ref)
		}
	}
}
