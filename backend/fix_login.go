package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
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

	// 1. Check if user exists
	var id, email, role, companyID string
	err = db.QueryRow("SELECT id, email, role, company_id FROM users WHERE email = 'ahmed@rider.ae'").Scan(&id, &email, &role, &companyID)
	if err != nil {
		log.Fatalf("User lookup failed: %v", err)
	}
	fmt.Printf("User Found: %s | ID: %s | Role: %s | Company: %s\n", email, id, role, companyID)

	// 2. Reset password to ensure access
	newPass := "Password123!"
	bytes, _ := bcrypt.GenerateFromPassword([]byte(newPass), 14) // using cost 14
	_, err = db.Exec("UPDATE users SET password_hash = $1 WHERE id = $2", string(bytes), id)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("✓ Password forcefully reset to: %s\n", newPass)
}
