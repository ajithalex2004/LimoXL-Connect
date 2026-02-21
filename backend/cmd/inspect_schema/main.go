package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func main() {
	// Hardcoded for debugging since user env might vary,
	// but better to use the specific URL they requested if known,
	// or try to load from same source.
	// User logs showed: postgresql://neondb_owner:npg_7ndWFKRYEOt6@ep-calm-heart-a15voo2a-pooler.ap-southeast-1.aws.neon.tech/neondb?sslmode=require
	connStr := "postgresql://neondb_owner:npg_7ndWFKRYEOt6@ep-calm-heart-a15voo2a-pooler.ap-southeast-1.aws.neon.tech/neondb?sslmode=require"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("Cannot connect:", err)
	}

	fmt.Println("Connected. Querying vehicles table info...")

	rows, err := db.Query(`
		SELECT column_name, data_type, is_nullable
		FROM information_schema.columns 
		WHERE table_name = 'drivers'
		ORDER BY ordinal_position;
	`)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	fmt.Println("--- Columns in 'drivers' ---")
	fmt.Printf("%-20s %-20s %s\n", "Column", "Type", "Nullable")
	found := false
	for rows.Next() {
		found = true
		var colName, colType, isNullable string
		if err := rows.Scan(&colName, &colType, &isNullable); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%-20s %-20s %s\n", colName, colType, isNullable)
	}
	if !found {
		fmt.Println("No columns found! Does the table 'vehicles' exist?")
	}
	fmt.Println("-----------------------------")
}
