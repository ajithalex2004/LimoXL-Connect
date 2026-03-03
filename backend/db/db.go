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
		"db/migrations/04_rfq_support.sql",
		"db/migrations/05_detailed_fleet.sql",
		"db/migrations/06_nui_masters.sql",
	}
	for _, path := range migrations {
		data, err := os.ReadFile(path)
		if err != nil {
			log.Printf("CRITICAL: Could not read migration %s: %v", path, err)
			return fmt.Errorf("could not read migration %s: %w", path, err)
		}
		log.Printf("Applying migration: %s", path)
		if _, err := DB.Exec(string(data)); err != nil {
			log.Printf("CRITICAL: Migration %s failed: %v", path, err)
			return fmt.Errorf("migration %s failed: %w", path, err)
		}
		log.Printf("Migration applied successfully: %s", path)
	}

	return nil
}

// SeedAdmin ensures a default admin company and user exist in the database.
// It is safe to call on every startup — it uses INSERT ... ON CONFLICT DO NOTHING.
func SeedAdmin() {
	adminEmail := "admin@limoxlink.com"
	adminPassword := "admin123"

	// 1. Diagnostics: Ensure Schema is correct (Idempotent)
	_, _ = DB.Exec(`ALTER TABLE users ALTER COLUMN company_id DROP NOT NULL`)
	_, _ = DB.Exec(`ALTER TABLE users ADD COLUMN IF NOT EXISTS is_super_admin BOOLEAN DEFAULT FALSE`)

	// 2. Hash the password fresh every startup
	hashed, err := bcrypt.GenerateFromPassword([]byte(adminPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("SeedAdmin: failed to hash password: %v", err)
		return
	}

	// 3. Upsert the SuperAdmin user (NO DEPENDENCY ON COMPANY)
	superEmail := "superadmin@limoxlink.com"
	_, err = DB.Exec(`
		INSERT INTO users (company_id, role, email, password_hash, name, is_super_admin)
		VALUES (NULL, 'SUPER_ADMIN', $1, $2, 'Super Admin', true)
		ON CONFLICT (email) DO UPDATE
		SET password_hash = EXCLUDED.password_hash,
		    role = 'SUPER_ADMIN',
		    is_super_admin = true
	`, superEmail, string(hashed))
	if err != nil {
		log.Printf("SeedAdmin: failed to upsert superadmin user: %v", err)
	} else {
		log.Printf("SeedAdmin: superadmin user '%s' is ready", superEmail)
	}

	// 4. Upsert the operator company
	var companyID string
	// Check if exists first since we don't have unique constraint on name usually
	err = DB.QueryRow(`SELECT id FROM companies WHERE name = 'LimoXL Operator' LIMIT 1`).Scan(&companyID)
	if err != nil {
		// Does not exist, create it
		err = DB.QueryRow(`
			INSERT INTO companies (name, type, verified)
			VALUES ('LimoXL Operator', 'DEMAND', true)
			RETURNING id
		`).Scan(&companyID)
		if err != nil {
			log.Printf("SeedAdmin: failed to create operator company: %v", err)
			return
		}
	}

	// 5. Upsert the admin user
	_, err = DB.Exec(`
		INSERT INTO users (company_id, role, email, password_hash, name)
		VALUES ($1, 'ADMIN', $2, $3, 'Admin')
		ON CONFLICT (email) DO UPDATE
		SET password_hash = EXCLUDED.password_hash,
		    company_id    = EXCLUDED.company_id
	`, companyID, adminEmail, string(hashed))
	if err != nil {
		log.Printf("SeedAdmin: failed to upsert admin user: %v", err)
	} else {
		log.Printf("SeedAdmin: admin user '%s' is ready", adminEmail)
	}

	// 6. Upsert Tenant entry for the operator
	var tenantID string
	err = DB.QueryRow(`
		INSERT INTO tenants (company_id, name, slug, status, plan)
		VALUES ($1, 'LimoXL Operator', 'limoxl-operator', 'ACTIVE', 'PROFESSIONAL')
		ON CONFLICT (company_id) DO UPDATE SET name = EXCLUDED.name
		RETURNING id
	`, companyID).Scan(&tenantID)
	if err != nil {
		_ = DB.QueryRow(`SELECT id FROM tenants WHERE company_id = $1`, companyID).Scan(&tenantID)
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
