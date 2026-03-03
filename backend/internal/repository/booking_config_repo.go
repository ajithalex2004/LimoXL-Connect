package repository

import (
	"context"
	"database/sql"
	"limoxlink-backend/internal/models"

	"github.com/google/uuid"
)

type PostgresBookingConfigRepo struct {
	db *sql.DB
}

func NewPostgresBookingConfigRepo(db *sql.DB) *PostgresBookingConfigRepo {
	return &PostgresBookingConfigRepo{db: db}
}

func (r *PostgresBookingConfigRepo) Create(ctx context.Context, config *models.BookingConfig) error {
	query := `
		INSERT INTO booking_configs (
			company_id, name, booking_type, request_type, priority, sort_order,
			vehicle_classes, vehicle_groups, vehicle_usages,
			pickup_buffer, auto_dispatch_buffer, pricing_source,
			approval_workflow_required, epod_required
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
		RETURNING id, created_at, updated_at`

	return r.db.QueryRowContext(ctx, query,
		config.CompanyID, config.Name, config.BookingType, config.RequestType, config.Priority, config.SortOrder,
		config.VehicleClasses, config.VehicleGroups, config.VehicleUsages,
		config.PickupBuffer, config.AutoDispatchBuffer, config.PricingSource,
		config.ApprovalWorkflowRequired, config.EPODRequired,
	).Scan(&config.ID, &config.CreatedAt, &config.UpdatedAt)
}

func (r *PostgresBookingConfigRepo) GetByID(ctx context.Context, id uuid.UUID) (*models.BookingConfig, error) {
	query := `SELECT id, company_id, name, booking_type, request_type, priority, sort_order,
	                 vehicle_classes, vehicle_groups, vehicle_usages,
	                 pickup_buffer, auto_dispatch_buffer, pricing_source,
	                 approval_workflow_required, epod_required, is_active, created_at, updated_at
	          FROM booking_configs WHERE id = $1`

	var config models.BookingConfig
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&config.ID, &config.CompanyID, &config.Name, &config.BookingType, &config.RequestType, &config.Priority, &config.SortOrder,
		&config.VehicleClasses, &config.VehicleGroups, &config.VehicleUsages,
		&config.PickupBuffer, &config.AutoDispatchBuffer, &config.PricingSource,
		&config.ApprovalWorkflowRequired, &config.EPODRequired, &config.IsActive, &config.CreatedAt, &config.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func (r *PostgresBookingConfigRepo) ListByCompany(ctx context.Context, companyID uuid.UUID) ([]models.BookingConfig, error) {
	query := `SELECT id, company_id, name, booking_type, request_type, priority, sort_order,
	                 vehicle_classes, vehicle_groups, vehicle_usages,
	                 pickup_buffer, auto_dispatch_buffer, pricing_source,
	                 approval_workflow_required, epod_required, is_active, created_at, updated_at
	          FROM booking_configs WHERE company_id = $1 ORDER BY sort_order ASC`

	rows, err := r.db.QueryContext(ctx, query, companyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var configs []models.BookingConfig
	for rows.Next() {
		var config models.BookingConfig
		err := rows.Scan(
			&config.ID, &config.CompanyID, &config.Name, &config.BookingType, &config.RequestType, &config.Priority, &config.SortOrder,
			&config.VehicleClasses, &config.VehicleGroups, &config.VehicleUsages,
			&config.PickupBuffer, &config.AutoDispatchBuffer, &config.PricingSource,
			&config.ApprovalWorkflowRequired, &config.EPODRequired, &config.IsActive, &config.CreatedAt, &config.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		configs = append(configs, config)
	}
	return configs, nil
}

func (r *PostgresBookingConfigRepo) Update(ctx context.Context, config *models.BookingConfig) error {
	query := `
		UPDATE booking_configs SET
			name = $1, booking_type = $2, request_type = $3, priority = $4, sort_order = $5,
			vehicle_classes = $6, vehicle_groups = $7, vehicle_usages = $8,
			pickup_buffer = $9, auto_dispatch_buffer = $10, pricing_source = $11,
			approval_workflow_required = $12, epod_required = $13, is_active = $14,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $15 AND company_id = $16`

	_, err := r.db.ExecContext(ctx, query,
		config.Name, config.BookingType, config.RequestType, config.Priority, config.SortOrder,
		config.VehicleClasses, config.VehicleGroups, config.VehicleUsages,
		config.PickupBuffer, config.AutoDispatchBuffer, config.PricingSource,
		config.ApprovalWorkflowRequired, config.EPODRequired, config.IsActive,
		config.ID, config.CompanyID,
	)
	return err
}

func (r *PostgresBookingConfigRepo) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM booking_configs WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
