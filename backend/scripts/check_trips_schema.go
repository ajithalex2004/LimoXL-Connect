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

	fmt.Println("=== Checking trips table schema ===")

	query := `
		SELECT column_name, data_type 
		FROM information_schema.columns 
		WHERE table_name = 'trips' 
		ORDER BY ordinal_position
	`

	rows, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	fmt.Println("\nColumns in 'trips' table:")
	hasPrice := false
	for rows.Next() {
		var colName, dataType string
		rows.Scan(&colName, &dataType)
		fmt.Printf("  - %s (%s)\n", colName, dataType)
		if colName == "price" {
			hasPrice = true
		}
	}

	if !hasPrice {
		fmt.Println("\n❌ 'price' column does NOT exist in trips table!")
		fmt.Println("This is why the query fails with 500 error")
	} else {
		fmt.Println("\n✓ 'price' column exists")
	}
}
