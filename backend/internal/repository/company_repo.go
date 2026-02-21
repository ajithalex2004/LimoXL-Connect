package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"limoxlink-backend/internal/models"
	"time"

	"github.com/google/uuid"
)

type PostgresCompanyRepo struct {
	DB *sql.DB
}

func NewPostgresCompanyRepo(db *sql.DB) *PostgresCompanyRepo {
	return &PostgresCompanyRepo{DB: db}
}

func (r *PostgresCompanyRepo) Create(ctx context.Context, company *models.Company) error {
	query := `
		INSERT INTO companies (id, name, type, contact_info, verified, settings, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	company.ID = uuid.New()
	company.CreatedAt = time.Now()
	company.UpdatedAt = time.Now()

	// Default empty JSON if nil
	if company.ContactInfo == nil {
		company.ContactInfo = json.RawMessage("{}")
	}
	if company.Settings == nil {
		company.Settings = json.RawMessage("{}")
	}

	_, err := r.DB.ExecContext(ctx, query,
		company.ID,
		company.Name,
		company.Type,
		company.ContactInfo,
		company.Verified,
		company.Settings,
		company.CreatedAt,
		company.UpdatedAt,
	)
	return err
}

func (r *PostgresCompanyRepo) GetByID(ctx context.Context, id uuid.UUID) (*models.Company, error) {
	query := `SELECT id, name, type, contact_info, verified, settings, created_at, updated_at FROM companies WHERE id = $1`

	row := r.DB.QueryRowContext(ctx, query, id)
	var c models.Company
	err := row.Scan(
		&c.ID, &c.Name, &c.Type, &c.ContactInfo, &c.Verified, &c.Settings, &c.CreatedAt, &c.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *PostgresCompanyRepo) List(ctx context.Context) ([]models.Company, error) {
	query := `SELECT id, name, type, contact_info, verified, settings, created_at, updated_at FROM companies`

	rows, err := r.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var companies []models.Company
	for rows.Next() {
		var c models.Company
		if err := rows.Scan(
			&c.ID, &c.Name, &c.Type, &c.ContactInfo, &c.Verified, &c.Settings, &c.CreatedAt, &c.UpdatedAt,
		); err != nil {
			return nil, err
		}
		companies = append(companies, c)
	}
	return companies, nil
}
