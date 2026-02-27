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
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL not set")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Failed to open DB: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping DB: %v", err)
	}

	// Hash the password
	hash, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Failed to hash password: %v", err)
	}

	// Upsert company
	var companyID string
	err = db.QueryRow(`
		INSERT INTO companies (name, type, verified)
		VALUES ('LimoXL Operator', 'DEMAND', true)
		RETURNING id
	`).Scan(&companyID)
	if err != nil {
		// Likely already exists — fetch it
		err2 := db.QueryRow(`SELECT id FROM companies WHERE name = 'LimoXL Operator' ORDER BY created_at ASC LIMIT 1`).Scan(&companyID)
		if err2 != nil {
			log.Fatalf("Failed to get company: %v | %v", err, err2)
		}
	}

	// Upsert admin user
	_, err = db.Exec(`
		INSERT INTO users (company_id, role, email, password_hash, name)
		VALUES ($1, 'ADMIN', 'admin@limoxlink.com', $2, 'Admin')
		ON CONFLICT (email) DO UPDATE
		SET password_hash = EXCLUDED.password_hash,
		    company_id    = EXCLUDED.company_id
	`, companyID, string(hash))
	if err != nil {
		log.Fatalf("Failed to upsert admin: %v", err)
	}

	fmt.Printf("✓ Admin user seeded! company_id=%s  password_hash=%s\n", companyID, string(hash))
}
