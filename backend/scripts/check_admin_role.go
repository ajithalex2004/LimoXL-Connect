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

	var id, email, role string
	err = db.QueryRow("SELECT id, email, role FROM users WHERE email = 'admin@limoxlink.com'").Scan(&id, &email, &role)
	if err != nil {
		log.Fatalf("User lookup failed: %v", err)
	}
	fmt.Printf("DEBUG_RESULT: User=%s, ID=%s, Role=%s\n", email, id, role)
}
