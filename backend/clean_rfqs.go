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

	// Delete trip offers related to these trips first to avoid foreign key violations (if any, though usually cascade or independent)
	// But `trip_offers` references `trips` via `trip_id`.
	// Let's delete offers for trips that are MARKETPLACE_SEARCH first.
	_, err = db.Exec(`
		DELETE FROM trip_offers 
		WHERE trip_id IN (SELECT id FROM trips WHERE status = 'MARKETPLACE_SEARCH')
	`)
	if err != nil {
		log.Printf("Error clearing related offers: %v", err)
	}

	// Delete the trips
	res, err := db.Exec("DELETE FROM trips WHERE status = 'MARKETPLACE_SEARCH'")
	if err != nil {
		log.Fatal(err)
	}

	count, _ := res.RowsAffected()
	fmt.Printf("✓ Deleted %d trips with status MARKETPLACE_SEARCH\n", count)
}
