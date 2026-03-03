package repository

import (
	"context"
	"database/sql"
	"fmt"
	"limoxlink-backend/internal/models"

	"github.com/google/uuid"
)

type PostgresNUIMasterRepo struct {
	DB *sql.DB
}

func NewPostgresNUIMasterRepo(db *sql.DB) *PostgresNUIMasterRepo {
	return &PostgresNUIMasterRepo{DB: db}
}

func (r *PostgresNUIMasterRepo) Create(ctx context.Context, master *models.NUIMaster) error {
	query := `
		INSERT INTO nui_masters (company_id, category, name, description, is_active)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at`
	
	err := r.DB.QueryRowContext(ctx, query, 
		master.CompanyID, master.Category, master.Name, master.Description, master.IsActive,
	).Scan(&master.ID, &master.CreatedAt, &master.UpdatedAt)
	
	if err != nil {
		return fmt.Errorf("error creating nui master: %w", err)
	}
	return nil
}

func (r *PostgresNUIMasterRepo) GetByID(ctx context.Context, id uuid.UUID) (*models.NUIMaster, error) {
	query := `SELECT id, company_id, category, name, description, is_active, created_at, updated_at FROM nui_masters WHERE id = $1`
	master := &models.NUIMaster{}
	err := r.DB.QueryRowContext(ctx, query, id).Scan(
		&master.ID, &master.CompanyID, &master.Category, &master.Name, &master.Description, &master.IsActive, &master.CreatedAt, &master.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("error getting nui master: %w", err)
	}
	return master, nil
}

func (r *PostgresNUIMasterRepo) ListByCompany(ctx context.Context, companyID uuid.UUID, category string) ([]models.NUIMaster, error) {
	query := `SELECT id, company_id, category, name, description, is_active, created_at, updated_at FROM nui_masters WHERE company_id = $1`
	args := []interface{}{companyID}
	
	if category != "" {
		query += " AND category = $2"
		args = append(args, category)
	}
	query += " ORDER BY name ASC"

	rows, err := r.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("error listing nui masters: %w", err)
	}
	defer rows.Close()

	var masters []models.NUIMaster
	for rows.Next() {
		var m models.NUIMaster
		if err := rows.Scan(&m.ID, &m.CompanyID, &m.Category, &m.Name, &m.Description, &m.IsActive, &m.CreatedAt, &m.UpdatedAt); err != nil {
			return nil, fmt.Errorf("error scanning nui master: %w", err)
		}
		masters = append(masters, m)
	}
	return masters, nil
}

func (r *PostgresNUIMasterRepo) Update(ctx context.Context, master *models.NUIMaster) error {
	query := `
		UPDATE nui_masters 
		SET category = $1, name = $2, description = $3, is_active = $4, updated_at = CURRENT_TIMESTAMP
		WHERE id = $5 AND company_id = $6`
	
	res, err := r.DB.ExecContext(ctx, query, 
		master.Category, master.Name, master.Description, master.IsActive, master.ID, master.CompanyID,
	)
	if err != nil {
		return fmt.Errorf("error updating nui master: %w", err)
	}
	
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("nui master not found")
	}
	return nil
}

func (r *PostgresNUIMasterRepo) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM nui_masters WHERE id = $1`
	_, err := r.DB.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("error deleting nui master: %w", err)
	}
	return nil
}
