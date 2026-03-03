package repository

import (
	"context"
	"database/sql"
	"limoxlink-backend/internal/models"
	"time"

	"github.com/google/uuid"
)

type PostgresDriverRepo struct {
	DB *sql.DB
}

func NewPostgresDriverRepo(db *sql.DB) *PostgresDriverRepo {
	return &PostgresDriverRepo{DB: db}
}

func (r *PostgresDriverRepo) Create(ctx context.Context, d *models.Driver) error {
	query := `
		INSERT INTO drivers (
			id, company_id, name, phone, license_number, 
			license_expiry, itc_permit_expiry, visa_expiry,
			first_name, last_name, date_of_birth, nationality, emirates_id,
			date_of_join, dallas_id, communication_language, hierarchy, driver_type,
			created_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19)
	`
	d.ID = uuid.New()
	d.CreatedAt = time.Now()

	_, err := r.DB.ExecContext(ctx, query,
		d.ID, d.CompanyID, d.Name, d.Phone, d.LicenseNumber,
		d.LicenseExpiry, d.ITCPermitExpiry, d.VisaExpiry,
		d.FirstName, d.LastName, d.DateOfBirth, d.Nationality, d.EmiratesID,
		d.DateOfJoin, d.DallasID, d.CommunicationLanguage, d.Hierarchy, d.DriverType,
		d.CreatedAt,
	)
	return err
}

func (r *PostgresDriverRepo) GetByID(ctx context.Context, id uuid.UUID) (*models.Driver, error) {
	query := `
		SELECT 
			id, company_id, name, phone, license_number, 
			license_expiry, itc_permit_expiry, visa_expiry,
			first_name, last_name, date_of_birth, nationality, emirates_id,
			date_of_join, dallas_id, communication_language, hierarchy, driver_type,
			current_vehicle_id, created_at
		FROM drivers WHERE id = $1
	`
	var d models.Driver
	var phone, license, vID, fName, lName, nation, eID, dID, lang, hierarchy, dType sql.NullString
	err := r.DB.QueryRowContext(ctx, query, id).Scan(
		&d.ID, &d.CompanyID, &d.Name, &phone, &license,
		&d.LicenseExpiry, &d.ITCPermitExpiry, &d.VisaExpiry,
		&fName, &lName, &d.DateOfBirth, &nation, &eID,
		&d.DateOfJoin, &dID, &lang, &hierarchy, &dType,
		&vID, &d.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	d.Phone = phone.String
	d.LicenseNumber = license.String
	d.FirstName = fName.String
	d.LastName = lName.String
	d.Nationality = nation.String
	d.EmiratesID = eID.String
	d.DallasID = dID.String
	d.CommunicationLanguage = lang.String
	d.Hierarchy = hierarchy.String
	d.DriverType = dType.String
	if vID.Valid {
		if pid, err := uuid.Parse(vID.String); err == nil {
			d.CurrentVehicleID = uuid.NullUUID{UUID: pid, Valid: true}
		}
	}
	return &d, nil
}

func (r *PostgresDriverRepo) ListByCompany(ctx context.Context, companyID uuid.UUID) ([]models.Driver, error) {
	query := `
		SELECT 
			id, company_id, name, phone, license_number, 
			license_expiry, itc_permit_expiry, visa_expiry,
			first_name, last_name, date_of_birth, nationality, emirates_id,
			date_of_join, dallas_id, communication_language, hierarchy, driver_type,
			current_vehicle_id, created_at
		FROM drivers 
		WHERE company_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.DB.QueryContext(ctx, query, companyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var drivers []models.Driver
	for rows.Next() {
		var d models.Driver
		var phone, license, vID, fName, lName, nation, eID, dID, lang, hierarchy, dType sql.NullString
		if err := rows.Scan(
			&d.ID, &d.CompanyID, &d.Name, &phone, &license,
			&d.LicenseExpiry, &d.ITCPermitExpiry, &d.VisaExpiry,
			&fName, &lName, &d.DateOfBirth, &nation, &eID,
			&d.DateOfJoin, &dID, &lang, &hierarchy, &dType,
			&vID, &d.CreatedAt,
		); err != nil {
			return nil, err
		}
		d.Phone = phone.String
		d.LicenseNumber = license.String
		d.FirstName = fName.String
		d.LastName = lName.String
		d.Nationality = nation.String
		d.EmiratesID = eID.String
		d.DallasID = dID.String
		d.CommunicationLanguage = lang.String
		d.Hierarchy = hierarchy.String
		d.DriverType = dType.String
		if vID.Valid {
			if parsedID, err := uuid.Parse(vID.String); err == nil {
				d.CurrentVehicleID = uuid.NullUUID{UUID: parsedID, Valid: true}
			}
		}
		drivers = append(drivers, d)
	}
	return drivers, nil
}

func (r *PostgresDriverRepo) Update(ctx context.Context, d *models.Driver) error {
	query := `
		UPDATE drivers SET
			name = $2, phone = $3, license_number = $4, 
			license_expiry = $5, itc_permit_expiry = $6, visa_expiry = $7,
			first_name = $8, last_name = $9, date_of_birth = $10, nationality = $11, emirates_id = $12,
			date_of_join = $13, dallas_id = $14, communication_language = $15, hierarchy = $16, driver_type = $17
		WHERE id = $1
	`
	_, err := r.DB.ExecContext(ctx, query,
		d.ID, d.Name, d.Phone, d.LicenseNumber,
		d.LicenseExpiry, d.ITCPermitExpiry, d.VisaExpiry,
		d.FirstName, d.LastName, d.DateOfBirth, d.Nationality, d.EmiratesID,
		d.DateOfJoin, d.DallasID, d.CommunicationLanguage, d.Hierarchy, d.DriverType,
	)
	return err
}

func (r *PostgresDriverRepo) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM drivers WHERE id = $1`
	_, err := r.DB.ExecContext(ctx, query, id)
	return err
}
