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

	fmt.Println("=== Checking invoices table schema ===")

	query := `
		SELECT column_name, data_type 
		FROM information_schema.columns 
		WHERE table_name = 'invoices' 
		ORDER BY ordinal_position
	`

	rows, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	fmt.Println("\nColumns in 'invoices' table:")
	for rows.Next() {
		var colName, dataType string
		rows.Scan(&colName, &dataType)
		fmt.Printf("  - %s (%s)\n", colName, dataType)
	}
}
