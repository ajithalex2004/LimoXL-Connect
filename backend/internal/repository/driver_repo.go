package repository

import (
	"context"
	"database/sql"
	"limoxlink-backend/internal/models"
	"time"

	"github.com/google/uuid"
)

type DriverRepository interface {
	Create(ctx context.Context, driver *models.Driver) error
	ListByCompany(ctx context.Context, companyID uuid.UUID) ([]models.Driver, error)
}

type PostgresDriverRepo struct {
	DB *sql.DB
}

func NewPostgresDriverRepo(db *sql.DB) *PostgresDriverRepo {
	return &PostgresDriverRepo{DB: db}
}

func (r *PostgresDriverRepo) Create(ctx context.Context, driver *models.Driver) error {
	query := `
		INSERT INTO drivers (
			id, company_id, name, contact_number, license_number, 
			license_expiry, itc_permit_expiry, visa_expiry, created_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	driver.ID = uuid.New()
	driver.CreatedAt = time.Now()

	_, err := r.DB.ExecContext(ctx, query,
		driver.ID,
		driver.CompanyID,
		driver.Name,
		driver.Phone, // Maps to contact_number
		driver.LicenseNumber,
		driver.LicenseExpiry,
		driver.ITCPermitExpiry,
		driver.VisaExpiry,
		driver.CreatedAt,
	)
	return err
}

func (r *PostgresDriverRepo) ListByCompany(ctx context.Context, companyID uuid.UUID) ([]models.Driver, error) {
	query := `
		SELECT 
			id, company_id, name, contact_number, license_number, 
			license_expiry, itc_permit_expiry, visa_expiry,
			assigned_vehicle_id, created_at
		FROM drivers 
		WHERE company_id = $1
	`

	rows, err := r.DB.QueryContext(ctx, query, companyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var drivers []models.Driver
	for rows.Next() {
		var d models.Driver
		var phone, license, vehicleID sql.NullString
		if err := rows.Scan(
			&d.ID, &d.CompanyID, &d.Name, &phone, &license,
			&d.LicenseExpiry, &d.ITCPermitExpiry, &d.VisaExpiry,
			&vehicleID, &d.CreatedAt,
		); err != nil {
			return nil, err
		}
		d.Phone = phone.String
		d.LicenseNumber = license.String
		if vehicleID.Valid {
			if parsedID, err := uuid.Parse(vehicleID.String); err == nil {
				d.CurrentVehicleID = uuid.NullUUID{UUID: parsedID, Valid: true}
			}
		}
		drivers = append(drivers, d)
	}
	return drivers, nil
}
