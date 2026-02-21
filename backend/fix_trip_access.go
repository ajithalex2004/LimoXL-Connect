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

	fmt.Println("Fixing trip_access foreign key constraint...")

	// Drop the existing foreign key constraint
	_, err = db.Exec(`
		ALTER TABLE trip_access 
		DROP CONSTRAINT IF EXISTS trip_access_company_id_fkey;
	`)
	if err != nil {
		log.Fatal("Error dropping constraint:", err)
	}

	fmt.Println("✓ Dropped old foreign key constraint")

	// For now, we'll just remove the constraint entirely since we're using outsource_companies
	// In a production system, you might want to add a new constraint to outsource_companies
	// or use a polymorphic relationship

	fmt.Println("\n✓ trip_access table updated successfully!")
	fmt.Println("Note: Foreign key constraint removed. Direct assignments will now work with outsource_companies.")
}
