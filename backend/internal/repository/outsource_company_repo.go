package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type OutsourceCompany struct {
	ID            uuid.UUID       `json:"id"`
	Name          string          `json:"name"`
	ContactPerson sql.NullString  `json:"contact_person"`
	Designation   sql.NullString  `json:"designation"`
	Email         sql.NullString  `json:"email"`
	ContactNumber sql.NullString  `json:"contact_number"`
	Address       sql.NullString  `json:"address"`
	City          sql.NullString  `json:"city"`
	Country       sql.NullString  `json:"country"`
	TradeLicense  sql.NullString  `json:"trade_license_no"`
	ITCPermit     sql.NullString  `json:"itc_permit_no"`
	VATNo         sql.NullString  `json:"vat_no"`
	SLAScore      float64         `json:"sla_score"`
	Rating        sql.NullFloat64 `json:"rating"`
	Status        string          `json:"verification_status"`
	IsActive      bool            `json:"is_active"`
	Notes         sql.NullString  `json:"notes"`
	CreatedAt     time.Time       `json:"created_at"`
	UpdatedAt     time.Time       `json:"updated_at"`
	DeletedAt     sql.NullTime    `json:"deleted_at,omitempty"`
}

type OutsourceCompanyRepo struct {
	db *sql.DB
}

func NewOutsourceCompanyRepo(db *sql.DB) *OutsourceCompanyRepo {
	return &OutsourceCompanyRepo{db: db}
}

func (r *OutsourceCompanyRepo) Create(ctx context.Context, company *OutsourceCompany) error {
	query := `
		INSERT INTO outsource_companies (
			id, name, contact_person, designation, email, contact_number,
			address, city, country, trade_license_no, itc_permit_no, vat_no,
			rating, is_active, notes
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
		RETURNING created_at, updated_at
	`

	if company.ID == uuid.Nil {
		company.ID = uuid.New()
	}

	return r.db.QueryRowContext(ctx, query,
		company.ID,
		company.Name,
		company.ContactPerson,
		company.Designation,
		company.Email,
		company.ContactNumber,
		company.Address,
		company.City,
		company.Country,
		company.TradeLicense,
		company.ITCPermit,
		company.VATNo,
		company.Rating,
		company.IsActive,
		company.Notes,
	).Scan(&company.CreatedAt, &company.UpdatedAt)
}

func (r *OutsourceCompanyRepo) GetByID(ctx context.Context, id uuid.UUID) (*OutsourceCompany, error) {
	query := `
		SELECT id, name, contact_person, designation, email, contact_number,
		       address, city, country, specialties, rating, is_active, notes,
		       created_at, updated_at, deleted_at
		FROM outsource_companies
		WHERE id = $1 AND deleted_at IS NULL
	`

	company := &OutsourceCompany{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&company.ID,
		&company.Name,
		&company.ContactPerson,
		&company.Designation,
		&company.Email,
		&company.ContactNumber,
		&company.Address,
		&company.City,
		&company.Country,
		&company.TradeLicense,
		&company.ITCPermit,
		&company.VATNo,
		&company.SLAScore,
		&company.Rating,
		&company.Status,
		&company.IsActive,
		&company.Notes,
		&company.CreatedAt,
		&company.UpdatedAt,
		&company.DeletedAt,
	)

	if err != nil {
		return nil, err
	}

	return company, nil
}

func (r *OutsourceCompanyRepo) List(ctx context.Context) ([]*OutsourceCompany, error) {
	query := `
		SELECT id, name, contact_person, designation, email, contact_number,
		       address, city, country, trade_license_no, itc_permit_no, vat_no,
		       rating, is_active, notes,
		       created_at, updated_at, deleted_at
		FROM outsource_companies
		WHERE deleted_at IS NULL
		ORDER BY name
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	companies := []*OutsourceCompany{}
	for rows.Next() {
		company := &OutsourceCompany{}
		err := rows.Scan(
			&company.ID,
			&company.Name,
			&company.ContactPerson,
			&company.Designation,
			&company.Email,
			&company.ContactNumber,
			&company.Address,
			&company.City,
			&company.Country,
			&company.TradeLicense,
			&company.ITCPermit,
			&company.VATNo,
			&company.Rating,
			&company.IsActive,
			&company.Notes,
			&company.CreatedAt,
			&company.UpdatedAt,
			&company.DeletedAt,
		)
		if err != nil {
			return nil, err
		}
		companies = append(companies, company)
	}

	return companies, rows.Err()
}

func (r *OutsourceCompanyRepo) Update(ctx context.Context, company *OutsourceCompany) error {
	query := `
		UPDATE outsource_companies
		SET name = $2, contact_person = $3, designation = $4, email = $5,
		    contact_number = $6, address = $7, city = $8, country = $9,
		    trade_license_no = $10, itc_permit_no = $11, vat_no = $12,
		    rating = $13, is_active = $14, notes = $15,
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = $1 AND deleted_at IS NULL
	`

	_, err := r.db.ExecContext(ctx, query,
		company.ID,
		company.Name,
		company.ContactPerson,
		company.Designation,
		company.Email,
		company.ContactNumber,
		company.Address,
		company.City,
		company.Country,
		company.TradeLicense,
		company.ITCPermit,
		company.VATNo,
		company.Rating,
		company.IsActive,
		company.Notes,
	)

	return err
}

func (r *OutsourceCompanyRepo) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE outsource_companies
		SET deleted_at = CURRENT_TIMESTAMP
		WHERE id = $1
	`

	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
