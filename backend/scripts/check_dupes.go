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

	fmt.Println("=== Checking Companies and Users ===")

	rows, err := db.Query(`
		SELECT c.id, c.name, u.email, u.name as user_name
		FROM companies c
		LEFT JOIN users u ON c.id = u.company_id
		WHERE c.name LIKE 'AAT%' OR c.name LIKE 'Rider%'
		ORDER BY c.name
	`)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var id, name string
		var email, userName sql.NullString
		if err := rows.Scan(&id, &name, &email, &userName); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Company: %s \n  ID: %s\n  User Email: %v\n  User Name: %v\n",
			name, id, email.String, userName.String)
		fmt.Println("------------------------------------------------")
	}
}
