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

	// Find duplicates by ID (not reference_no)
	fmt.Println("=== Checking for duplicate IDs ===")
	rows, err := db.Query(`
		SELECT id, COUNT(*) as count
		FROM trips
		GROUP BY id
		HAVING COUNT(*) > 1
	`)
	if err != nil {
		log.Fatal(err)
	}

	dupCount := 0
	for rows.Next() {
		var id string
		var count int
		rows.Scan(&id, &count)
		fmt.Printf("Duplicate ID: %s (appears %d times)\n", id, count)
		dupCount++
	}
	rows.Close()

	if dupCount == 0 {
		fmt.Println("✓ No duplicate IDs found")
	}

	// Check for trips with same reference_no but different IDs
	fmt.Println("\n=== Checking for duplicate reference numbers ===")
	rows, err = db.Query(`
		SELECT reference_no, COUNT(DISTINCT id) as id_count, array_agg(id::text) as ids
		FROM trips
		GROUP BY reference_no
		HAVING COUNT(DISTINCT id) > 1
	`)
	if err != nil {
		log.Fatal(err)
	}

	dupRefCount := 0
	for rows.Next() {
		var refNo string
		var idCount int
		var ids string
		rows.Scan(&refNo, &idCount, &ids)
		fmt.Printf("Duplicate Ref: %s has %d different IDs\n", refNo, idCount)
		dupRefCount++
	}
	rows.Close()

	if dupRefCount > 0 {
		fmt.Printf("\n⚠️  Found %d reference numbers with multiple IDs\n", dupRefCount)
		fmt.Println("Cleaning up duplicates (keeping most recent)...")

		// Delete older duplicates
		result, err := db.Exec(`
			DELETE FROM trips
			WHERE id IN (
				SELECT id
				FROM (
					SELECT id, reference_no,
						ROW_NUMBER() OVER (PARTITION BY reference_no ORDER BY created_at DESC) as rn
					FROM trips
				) t
				WHERE rn > 1
			)
		`)
		if err != nil {
			log.Fatal("Failed to clean duplicates:", err)
		}

		deleted, _ := result.RowsAffected()
		fmt.Printf("✓ Deleted %d duplicate trips\n", deleted)
	} else {
		fmt.Println("✓ No duplicate reference numbers found")
	}
}
