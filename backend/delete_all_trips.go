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

	fmt.Printf("Connecting to: %s\n\n", connStr)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Count trips before deletion
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM trips").Scan(&count)
	if err != nil {
		log.Fatal("Error counting trips:", err)
	}

	fmt.Printf("Found %d trips in the database\n", count)

	if count == 0 {
		fmt.Println("No trips to delete.")
		return
	}

	fmt.Println("\nDeleting all trips and related data...")

	// Delete in order to respect foreign key constraints
	// 1. Delete trip_access entries
	result, err := db.Exec("DELETE FROM trip_access")
	if err != nil {
		log.Fatal("Error deleting trip_access:", err)
	}
	rows, _ := result.RowsAffected()
	fmt.Printf("✓ Deleted %d trip_access entries\n", rows)

	// 2. Delete trip_offers
	result, err = db.Exec("DELETE FROM trip_offers")
	if err != nil {
		log.Fatal("Error deleting trip_offers:", err)
	}
	rows, _ = result.RowsAffected()
	fmt.Printf("✓ Deleted %d trip_offers\n", rows)

	// 3. Delete invoices
	result, err = db.Exec("DELETE FROM invoices")
	if err != nil {
		log.Fatal("Error deleting invoices:", err)
	}
	rows, _ = result.RowsAffected()
	fmt.Printf("✓ Deleted %d invoices\n", rows)

	// 4. Delete trips
	result, err = db.Exec("DELETE FROM trips")
	if err != nil {
		log.Fatal("Error deleting trips:", err)
	}
	rows, _ = result.RowsAffected()
	fmt.Printf("✓ Deleted %d trips\n", rows)

	// Reset sequences
	_, err = db.Exec("ALTER SEQUENCE trip_id_seq RESTART WITH 260001")
	if err != nil {
		fmt.Printf("Warning: Could not reset trip_id_seq: %v\n", err)
	} else {
		fmt.Println("✓ Reset trip_id_seq to 260001")
	}

	_, err = db.Exec("ALTER SEQUENCE rfq_id_seq RESTART WITH 1001")
	if err != nil {
		fmt.Printf("Warning: Could not reset rfq_id_seq: %v\n", err)
	} else {
		fmt.Println("✓ Reset rfq_id_seq to 1001")
	}

	fmt.Println("\n✓ All trips and related data have been deleted successfully!")
}
