package repository

import (
	"context"
	"database/sql"
	"limoxlink-backend/internal/models"
	"time"

	"github.com/google/uuid"
)

type PostgresUserRepo struct {
	DB *sql.DB
}

func NewPostgresUserRepo(db *sql.DB) *PostgresUserRepo {
	return &PostgresUserRepo{DB: db}
}

func (r *PostgresUserRepo) Create(ctx context.Context, user *models.User) error {
	query := `
		INSERT INTO users (id, company_id, role, email, password_hash, name, is_super_admin, password_change_required, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	user.ID = uuid.New()
	user.CreatedAt = time.Now()

	_, err := r.DB.ExecContext(ctx, query,
		user.ID,
		user.CompanyID,
		user.Role,
		user.Email,
		user.PasswordHash,
		user.Name,
		user.IsSuperAdmin,
		user.PasswordChangeRequired,
		user.CreatedAt,
	)
	return err
}

func (r *PostgresUserRepo) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `SELECT id, company_id, role, email, password_hash, name, is_super_admin, password_change_required, created_at FROM users WHERE email = $1`

	row := r.DB.QueryRowContext(ctx, query, email)
	var u models.User
	err := row.Scan(
		&u.ID, &u.CompanyID, &u.Role, &u.Email, &u.PasswordHash, &u.Name, &u.IsSuperAdmin, &u.PasswordChangeRequired, &u.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *PostgresUserRepo) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	query := `SELECT id, company_id, role, email, password_hash, name, is_super_admin, password_change_required, created_at FROM users WHERE id = $1`

	row := r.DB.QueryRowContext(ctx, query, id)
	var u models.User
	err := row.Scan(
		&u.ID, &u.CompanyID, &u.Role, &u.Email, &u.PasswordHash, &u.Name, &u.IsSuperAdmin, &u.PasswordChangeRequired, &u.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *PostgresUserRepo) UpdatePassword(ctx context.Context, id uuid.UUID, hash string) error {
	query := `UPDATE users SET password_hash = $1, password_change_required = $2 WHERE id = $3`
	_, err := r.DB.ExecContext(ctx, query, hash, false, id)
	return err
}

func (r *PostgresUserRepo) ListByCompany(ctx context.Context, companyID uuid.UUID) ([]*models.User, error) {
	query := `SELECT id, company_id, role, email, password_hash, name, is_super_admin, password_change_required, created_at FROM users WHERE company_id = $1 ORDER BY created_at DESC`

	rows, err := r.DB.QueryContext(ctx, query, companyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		var u models.User
		err := rows.Scan(
			&u.ID, &u.CompanyID, &u.Role, &u.Email, &u.PasswordHash, &u.Name, &u.IsSuperAdmin, &u.PasswordChangeRequired, &u.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, &u)
	}
	return users, nil
}
