package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func main() {
	connStr := "postgresql://neondb_owner:npg_7ndWFKRYEOt6@ep-calm-heart-a15voo2a-pooler.ap-southeast-1.aws.neon.tech/neondb?sslmode=require"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to DB. Applying patch...")

	// 0. Enable Extensions
	extensions := []string{"uuid-ossp", "postgis"}
	for _, ext := range extensions {
		_, err := db.Exec(fmt.Sprintf("CREATE EXTENSION IF NOT EXISTS \"%s\";", ext))
		if err != nil {
			log.Printf("Error enabling extension %s: %v", ext, err)
		} else {
			fmt.Printf("Extension %s enabled.\n", ext)
		}
	}

	// 1. Check if vehicles table exists
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS vehicles (
		id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
		plate_number VARCHAR(50) UNIQUE NOT NULL,
		type VARCHAR(50) NOT NULL,
		capacity INT NOT NULL,
		status VARCHAR(50) NOT NULL DEFAULT 'OFFLINE',
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	);`)
	if err != nil {
		log.Printf("Error ensuring table exists: %v", err)
	}

	// 2. Add missing columns
	queries := []string{
		`ALTER TABLE vehicles ADD COLUMN IF NOT EXISTS company_id UUID;`, // The missing one
		`ALTER TABLE vehicles ADD COLUMN IF NOT EXISTS current_location GEOGRAPHY(POINT, 4326);`,
		`ALTER TABLE vehicles ADD COLUMN IF NOT EXISTS last_heartbeat TIMESTAMP WITH TIME ZONE;`,
	}

	for _, q := range queries {
		_, err := db.Exec(q)
		if err != nil {
			log.Printf("Error executing query %q: %v", q, err)
		} else {
			fmt.Printf("Executed: %s\n", q)
		}
	}

	// 3. Check for trip_offers table
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS trip_offers (
		id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
		trip_id UUID REFERENCES trips(id),
		supplier_company_id UUID REFERENCES companies(id),
		status VARCHAR(50) NOT NULL DEFAULT 'PENDING',
		score FLOAT DEFAULT 0.0,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
		expires_at TIMESTAMP WITH TIME ZONE
	);`)
	if err != nil {
		log.Printf("Error ensuring trip_offers table exists: %v", err)
	} else {
		fmt.Println("Checked/Created trip_offers table.")
	}

	fmt.Println("Patch complete.")
}
