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

	fmt.Println("Cleaning up duplicate companies...")

	// Valid Company IDs (ones with users)
	validAAT := "b9bf3f64-a9d9-43b9-b051-40d004f3e455"
	validRider := "24a8f9b1-5bdd-459e-826c-586946f176e6"

	// Delete others
	res, err := db.Exec(`
		DELETE FROM companies 
		WHERE (name LIKE 'AAT%' OR name LIKE 'Rider%')
		AND id NOT IN ($1, $2)
	`, validAAT, validRider)

	if err != nil {
		// If delete fails due to constraints (e.g. they have assets), we might need to handle that,
		// but given they are likely ghosts, it should be fine.
		log.Fatal("Error deleting duplicates:", err)
	}

	rows, _ := res.RowsAffected()
	fmt.Printf("✓ Deleted %d duplicate company records\n", rows)
}
