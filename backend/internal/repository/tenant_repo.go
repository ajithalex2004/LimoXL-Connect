package repository

import (
	"context"
	"database/sql"
	"fmt"
	"limoxlink-backend/internal/models"

	"github.com/google/uuid"
)

type TenantRepository struct {
	DB *sql.DB
}

func NewTenantRepository(db *sql.DB) *TenantRepository {
	return &TenantRepository{DB: db}
}

// GetByCompanyID returns the tenant associated with a company
func (r *TenantRepository) GetByCompanyID(ctx context.Context, companyID uuid.UUID) (*models.Tenant, error) {
	query := `SELECT id, company_id, name, slug, status, plan, max_users, created_at, updated_at 
	          FROM tenants WHERE company_id = $1`
	
	t := &models.Tenant{}
	var companyIDNull uuid.NullUUID
	err := r.DB.QueryRowContext(ctx, query, companyID).Scan(
		&t.ID, &companyIDNull, &t.Name, &t.Slug, &t.Status, &t.Plan, &t.MaxUsers, &t.CreatedAt, &t.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	if companyIDNull.Valid {
		t.CompanyID = &companyIDNull.UUID
	}
	return t, nil
}

// ListAll returns all tenants for SuperAdmin
func (r *TenantRepository) ListAll(ctx context.Context) ([]models.TenantWithFeatures, error) {
	query := `
		SELECT t.id, t.company_id, t.name, t.slug, t.status, t.plan, t.max_users, t.created_at, t.updated_at,
		       c.name as company_name,
		       (SELECT COUNT(*) FROM users u WHERE u.company_id = t.company_id) as user_count,
		       (SELECT COUNT(*) FROM trips tr WHERE tr.requesting_company_id = t.company_id) as trip_count
		FROM tenants t
		LEFT JOIN companies c ON t.company_id = c.id
		ORDER BY t.created_at DESC
	`
	rows, err := r.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tenants []models.TenantWithFeatures
	for rows.Next() {
		var twf models.TenantWithFeatures
		var companyIDNull uuid.NullUUID
		err := rows.Scan(
			&twf.ID, &companyIDNull, &twf.Name, &twf.Slug, &twf.Status, &twf.Plan, &twf.MaxUsers, &twf.CreatedAt, &twf.UpdatedAt,
			&twf.CompanyName, &twf.UserCount, &twf.TripCount,
		)
		if err != nil {
			return nil, err
		}
		if companyIDNull.Valid {
			twf.CompanyID = &companyIDNull.UUID
		}
		
		// Fetch features
		features, err := r.GetFeatures(ctx, twf.ID)
		if err == nil {
			twf.Features = features
		}
		
		tenants = append(tenants, twf)
	}
	return tenants, nil
}

// GetFeatures returns feature flags for a tenant
func (r *TenantRepository) GetFeatures(ctx context.Context, tenantID uuid.UUID) ([]models.TenantFeature, error) {
	query := `SELECT id, tenant_id, feature_key, is_enabled, config, created_at, updated_at 
	          FROM tenant_features WHERE tenant_id = $1`
	
	rows, err := r.DB.QueryContext(ctx, query, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var features []models.TenantFeature
	for rows.Next() {
		var f models.TenantFeature
		err := rows.Scan(&f.ID, &f.TenantID, &f.FeatureKey, &f.IsEnabled, &f.Config, &f.CreatedAt, &f.UpdatedAt)
		if err != nil {
			continue
		}
		features = append(features, f)
	}
	return features, nil
}

// IsFeatureEnabled checks if a specific feature is enabled for a tenant
func (r *TenantRepository) IsFeatureEnabled(ctx context.Context, tenantID uuid.UUID, featureKey string) (bool, error) {
	var enabled bool
	query := `SELECT is_enabled FROM tenant_features WHERE tenant_id = $1 AND feature_key = $2`
	err := r.DB.QueryRowContext(ctx, query, tenantID, featureKey).Scan(&enabled)
	if err == sql.ErrNoRows {
		return false, nil
	}
	return enabled, err
}

// UpdateFeature toggles a feature for a tenant
func (r *TenantRepository) UpdateFeature(ctx context.Context, tenantID uuid.UUID, featureKey string, enabled bool) error {
	query := `
		INSERT INTO tenant_features (tenant_id, feature_key, is_enabled)
		VALUES ($1, $2, $3)
		ON CONFLICT (tenant_id, feature_key) 
		DO UPDATE SET is_enabled = $3, updated_at = NOW()
	`
	_, err := r.DB.ExecContext(ctx, query, tenantID, featureKey, enabled)
	return err
}

// Create creates a new tenant
func (r *TenantRepository) Create(ctx context.Context, t *models.Tenant) error {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	query := `
		INSERT INTO tenants (id, company_id, name, slug, status, plan, max_users)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING created_at, updated_at
	`
	return r.DB.QueryRowContext(ctx, query, t.ID, t.CompanyID, t.Name, t.Slug, t.Status, t.Plan, t.MaxUsers).Scan(&t.CreatedAt, &t.UpdatedAt)
}

// Update updates tenant details
func (r *TenantRepository) Update(ctx context.Context, t *models.Tenant) error {
	query := `
		UPDATE tenants 
		SET name = $2, slug = $3, status = $4, plan = $5, max_users = $6, updated_at = NOW()
		WHERE id = $1
	`
	_, err := r.DB.ExecContext(ctx, query, t.ID, t.Name, t.Slug, t.Status, t.Plan, t.MaxUsers)
	return err
}
