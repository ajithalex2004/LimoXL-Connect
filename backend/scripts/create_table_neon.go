package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func main() {
	// Use the same connection string as the server
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

	// Read the SQL file
	sqlBytes, err := os.ReadFile("create_outsource_companies_table.sql")
	if err != nil {
		log.Fatal("Error reading SQL file:", err)
	}

	fmt.Println("Creating outsource_companies table...")

	// Execute the SQL
	_, err = db.Exec(string(sqlBytes))
	if err != nil {
		log.Fatal("Error creating table:", err)
	}

	fmt.Println("✓ Successfully created outsource_companies table")

	// Verify
	var exists bool
	err = db.QueryRow(`
		SELECT EXISTS (
			SELECT FROM information_schema.tables 
			WHERE table_schema = 'public' 
			AND table_name = 'outsource_companies'
		)
	`).Scan(&exists)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\nTable exists: %v\n", exists)
}
