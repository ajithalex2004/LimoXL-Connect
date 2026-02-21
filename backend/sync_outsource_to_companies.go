package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
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

	fmt.Println("Syncing outsource_companies to companies table...")

	// Get all outsource companies
	rows, err := db.Query(`
		SELECT id, name, email, contact_person
		FROM outsource_companies
		WHERE deleted_at IS NULL
	`)
	if err != nil {
		log.Fatal("Error querying outsource companies:", err)
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		var id, name string
		var email, contactPerson sql.NullString

		err := rows.Scan(&id, &name, &email, &contactPerson)
		if err != nil {
			log.Fatal("Error scanning row:", err)
		}

		// Check if company already exists in companies table
		var existingID string
		err = db.QueryRow(`SELECT id FROM companies WHERE name = $1`, name).Scan(&existingID)

		if err == sql.ErrNoRows {
			// Create company
			companyID := uuid.MustParse(id)
			_, err = db.Exec(`
				INSERT INTO companies (id, name, type, verified, created_at, updated_at)
				VALUES ($1, $2, 'SUPPLY', true, NOW(), NOW())
			`, companyID, name)
			if err != nil {
				fmt.Printf("Error creating company %s: %v\n", name, err)
				continue
			}
			fmt.Printf("✓ Created company: %s\n", name)

			// Create user if email is provided
			if email.Valid && email.String != "" {
				// Generate default password
				password := "Password123!"
				hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
				if err != nil {
					fmt.Printf("Error hashing password for %s: %v\n", name, err)
					continue
				}

				userName := contactPerson.String
				if userName == "" {
					userName = name + " Admin"
				}

				_, err = db.Exec(`
					INSERT INTO users (company_id, email, password_hash, name, role, password_change_required, created_at)
					VALUES ($1, $2, $3, $4, 'SUPPLIER_ADMIN', true, NOW())
				`, companyID, email.String, string(hashedPassword), userName)
				if err != nil {
					fmt.Printf("Error creating user for %s: %v\n", name, err)
				} else {
					fmt.Printf("  ✓ Created user: %s (email: %s, password: %s)\n", userName, email.String, password)
				}
			}
			count++
		} else if err != nil {
			fmt.Printf("Error checking company %s: %v\n", name, err)
		} else {
			fmt.Printf("- Company already exists: %s\n", name)
		}
	}

	fmt.Printf("\n✓ Synced %d outsource companies to companies table!\n", count)
	fmt.Println("\nDefault password for all new users: Password123!")
	fmt.Println("Users will be required to change password on first login.")
}
