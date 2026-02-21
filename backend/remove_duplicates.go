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

	// Find and remove duplicate trips (keeping only the most recent one)
	query := `
		DELETE FROM trips
		WHERE id IN (
			SELECT id
			FROM (
				SELECT id, reference_no,
					ROW_NUMBER() OVER (PARTITION BY reference_no ORDER BY created_at DESC) as rn
				FROM trips
			) t
			WHERE rn > 1
		)
	`

	result, err := db.Exec(query)
	if err != nil {
		log.Fatal("Failed to remove duplicates:", err)
	}

	rows, _ := result.RowsAffected()
	fmt.Printf("✓ Removed %d duplicate trips\n", rows)

	// Also clean up orphaned trip_offers
	_, err = db.Exec(`
		DELETE FROM trip_offers
		WHERE trip_id NOT IN (SELECT id FROM trips)
	`)
	if err != nil {
		log.Printf("Warning: Failed to clean orphaned offers: %v\n", err)
	} else {
		fmt.Println("✓ Cleaned up orphaned trip offers")
	}
}
