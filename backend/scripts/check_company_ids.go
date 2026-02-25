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

	fmt.Println("Checking company IDs...")

	// List companies
	rows, err := db.Query(`SELECT id, name, type FROM companies WHERE type IN ('SUPPLY', 'BOTH') ORDER BY name`)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("\n=== COMPANIES TABLE ===")
	for rows.Next() {
		var id, name, ctype string
		rows.Scan(&id, &name, &ctype)
		fmt.Printf("ID: %s | Name: %s | Type: %s\n", id, name, ctype)
	}
	rows.Close()

	// List outsource companies
	rows, err = db.Query(`SELECT id, name FROM outsource_companies WHERE deleted_at IS NULL ORDER BY name`)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("\n=== OUTSOURCE_COMPANIES TABLE ===")
	for rows.Next() {
		var id, name string
		rows.Scan(&id, &name)
		fmt.Printf("ID: %s | Name: %s\n", id, name)
	}
	rows.Close()

	fmt.Println("\n=== SOLUTION ===")
	fmt.Println("The outsource_companies table should reference the same IDs as companies table.")
	fmt.Println("For now, use the companies table IDs when assigning trips.")
}
