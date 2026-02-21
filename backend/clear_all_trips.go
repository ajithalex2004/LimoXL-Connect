There should be provision to resubmit Invoice and Close Invoice.package main

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

	fmt.Println("=== Clearing ALL trips from the system ===")

	// Delete all trip offers first (foreign key constraint)
	result, err := db.Exec("DELETE FROM trip_offers")
	if err != nil {
		log.Fatal("Failed to delete trip offers:", err)
	}
	offers, _ := result.RowsAffected()
	fmt.Printf("✓ Deleted %d trip offers\n", offers)

	// Delete all trip access records
	result, err = db.Exec("DELETE FROM trip_access")
	if err != nil {
		log.Printf("Warning: Failed to delete trip access: %v\n", err)
	} else {
		access, _ := result.RowsAffected()
		fmt.Printf("✓ Deleted %d trip access records\n", access)
	}

	// Delete all trips
	result, err = db.Exec("DELETE FROM trips")
	if err != nil {
		log.Fatal("Failed to delete trips:", err)
	}
	trips, _ := result.RowsAffected()
	fmt.Printf("✓ Deleted %d trips\n", trips)

	fmt.Println("\n✅ All trips have been permanently removed from the system")
}
