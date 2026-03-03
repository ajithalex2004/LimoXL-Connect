package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/lib/pq"
)

func loadEnv() string {
	f, err := os.Open(".env")
	if err != nil {
		fmt.Printf("Warning: .env not found: %v\n", err)
		return ""
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "DATABASE_URL=") {
			return strings.TrimPrefix(line, "DATABASE_URL=")
		}
	}
	return ""
}

func main() {
	connStr := loadEnv()
	if connStr == "" {
		connStr = os.Getenv("DATABASE_URL")
	}
	if connStr == "" {
		fmt.Println("CRITICAL: DATABASE_URL not found in .env or environment")
		return
	}

	fmt.Printf("Connecting to DB (length: %d)...\n", len(connStr))

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal("Ping failed:", err)
	}
	fmt.Println("✓ Connected successfully")

	fmt.Println("\n--- Database Diagnostics ---")

	// 1. Check trips count
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM trips").Scan(&count)
	if err != nil {
		log.Fatal("Error counting trips:", err)
	}
	fmt.Printf("Total Trips in DB: %d\n", count)

	// 2. Check trips sample
	rows, err := db.Query("SELECT id, reference_no, requesting_company_id, status, created_at FROM trips ORDER BY created_at DESC LIMIT 5")
	if err != nil {
		fmt.Printf("Error querying trips: %v\n", err)
	} else {
		fmt.Println("\nRecent Trips Sample:")
		for rows.Next() {
			var id, ref, coId, status, createdAt string
			rows.Scan(&id, &ref, &coId, &status, &createdAt)
			fmt.Printf("ID: %s | Ref: %s | CoID: %s | Status: %s | CreatedAt: %s\n", id, ref, coId, status, createdAt)
		}
		rows.Close()
	}

	// 3. Check schema columns for trips
	rows, err = db.Query(`
		SELECT column_name, data_type 
		FROM information_schema.columns 
		WHERE table_name = 'trips'
		ORDER BY ordinal_position
	`)
	if err == nil {
		fmt.Println("\nTrips Table columns:")
		for rows.Next() {
			var name, dtype string
			rows.Scan(&name, &dtype)
			fmt.Printf("- %s (%s)\n", name, dtype)
		}
		rows.Close()
	}

	// 4. Check companies
	rows, err = db.Query("SELECT id, name FROM companies LIMIT 10")
	if err == nil {
		fmt.Println("\nCompanies Sample:")
		for rows.Next() {
			var id, name string
			rows.Scan(&id, &name)
			fmt.Printf("- %s: %s\n", id, name)
		}
		rows.Close()
	}
}
