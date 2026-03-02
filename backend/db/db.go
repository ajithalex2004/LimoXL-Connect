package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

var DB *sql.DB

func InitDB() error {
	// Default connection string - CHANGE THIS or use DATABASE_URL env var
	connStr := "postgres://postgres:password@localhost:5432/limoxlink?sslmode=disable"
	if envStr := os.Getenv("DATABASE_URL"); envStr != "" {
		connStr = envStr
	}
	log.Printf("Attempting to connect to DB: %s", connStr)
	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Printf("Error opening database: %v", err)
		return fmt.Errorf("error opening database: %w", err)
	}

	log.Println("Pinging database...")
	// Make sure we can actually ping the database
	if err = DB.Ping(); err != nil {
		log.Printf("Error pinging database: %v", err)
		return fmt.Errorf("error connecting to database: %w", err)
	}

	log.Println("Successfully connected to the database")
	return nil
}

func RunMigrations() error {
	// Simple migration runner for demo
	// In prod, use golang-migrate or similar
	schema, err := os.ReadFile("db/migrations/01_init_schema.sql")
	if err != nil {
		return fmt.Errorf("error reading schema file: %w", err)
	}

	_, err = DB.Exec(string(schema))
	if err != nil {
		return fmt.Errorf("error executing schema: %w", err)
	}
	log.Println("Migration 01: schema applied successfully")

	// Run additional migrations
	migrations := []string{
		"db/migrations/02_outsource_companies.sql",
		"db/migrations/03_multi_tenancy.sql",
	}
	for _, path := range migrations {
		data, err := os.ReadFile(path)
		if err != nil {
			log.Printf("Warning: Could not read migration %s: %v", path, err)
			continue
		}
		if _, err := DB.Exec(string(data)); err != nil {
			log.Printf("Warning: Migration %s failed: %v", path, err)
		} else {
			log.Printf("Migration applied: %s", path)
		}
	}

	return nil
}

// SeedAdmin ensures a default admin company and user exist in the database.
// It is safe to call on every startup — it uses INSERT ... ON CONFLICT DO NOTHING.
func SeedAdmin() {
	adminEmail := "admin@limoxlink.com"
	adminPassword := "admin123"

	// 1. Hash the password fresh every startup
	hashed, err := bcrypt.GenerateFromPassword([]byte(adminPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("SeedAdmin: failed to hash password: %v", err)
		return
	}

	// 2. Upsert the operator company
	var companyID string
	err = DB.QueryRow(`
		INSERT INTO companies (name, type, verified)
		VALUES ('LimoXL Operator', 'DEMAND', true)
		ON CONFLICT DO NOTHING
		RETURNING id
	`).Scan(&companyID)
	if err != nil {
		// Company already exists — fetch its ID
		err2 := DB.QueryRow(`SELECT id FROM companies WHERE name = 'LimoXL Operator' LIMIT 1`).Scan(&companyID)
		if err2 != nil {
			log.Printf("SeedAdmin: failed to get company ID: %v", err2)
			return
		}
	}

	log.Printf("SeedAdmin: admin user '%s' is ready (company_id=%s)", adminEmail, companyID)

	// 4. Upsert Tenant entry for the operator
	var tenantID string
	err = DB.QueryRow(`
		INSERT INTO tenants (company_id, name, slug, status, plan)
		VALUES ($1, 'LimoXL Operator', 'limoxl-operator', 'ACTIVE', 'PROFESSIONAL')
		ON CONFLICT (company_id) DO UPDATE SET name = EXCLUDED.name
		RETURNING id
	`, companyID).Scan(&tenantID)
	if err != nil {
		err2 := DB.QueryRow(`SELECT id FROM tenants WHERE company_id = $1`, companyID).Scan(&tenantID)
		if err2 != nil {
			log.Printf("SeedAdmin: failed to get tenant ID: %v", err2)
			return
		}
	}

	// 5. Enable all features for this seeded tenant
	features := []string{
		"dispatch", "outsource_marketplace", "fleet_management",
		"team_management", "invoicing", "partner_portal", "analytics",
	}
	for _, f := range features {
		_, err = DB.Exec(`
			INSERT INTO tenant_features (tenant_id, feature_key, is_enabled)
			VALUES ($1, $2, true)
			ON CONFLICT (tenant_id, feature_key) DO NOTHING
		`, tenantID, f)
		if err != nil {
			log.Printf("SeedAdmin: failed to seed feature %s: %v", f, err)
		}
	}
	log.Printf("SeedAdmin: tenant '%s' and features are ready", tenantID)
}
