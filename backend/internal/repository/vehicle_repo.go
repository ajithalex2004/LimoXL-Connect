package repository

import (
	"context"
	"database/sql"
	"limoxlink-backend/internal/models"
	"time"

	"github.com/google/uuid"
)

type PostgresVehicleRepo struct {
	DB *sql.DB
}

func NewPostgresVehicleRepo(db *sql.DB) *PostgresVehicleRepo {
	return &PostgresVehicleRepo{DB: db}
}

func (r *PostgresVehicleRepo) Create(ctx context.Context, vehicle *models.Vehicle) error {
	// PostGIS: ST_SetSRID(ST_MakePoint(lng, lat), 4326)
	// For creation, we assume lat/lng might be nil or 0 initially.
	// If provided, we insert. If not, we can leave geo null or 0,0.
	query := `
		INSERT INTO vehicles (id, company_id, license_plate, type, vehicle_class, vehicle_group, model, capacity, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`
	vehicle.ID = uuid.New()
	vehicle.CreatedAt = time.Now()
	vehicle.UpdatedAt = time.Now()

	_, err := r.DB.ExecContext(ctx, query,
		vehicle.ID,
		vehicle.CompanyID,
		vehicle.PlateNumber,
		vehicle.Type,
		vehicle.VehicleClass,
		vehicle.VehicleGroup,
		vehicle.Model,
		vehicle.Capacity,
		vehicle.Status,
		vehicle.CreatedAt,
		vehicle.UpdatedAt,
	)
	return err
}

func (r *PostgresVehicleRepo) GetByID(ctx context.Context, id uuid.UUID) (*models.Vehicle, error) {
	// ST_Y(geo::geometry) as lat, ST_X(geo::geometry) as lng
	query := `
		SELECT 
			id, company_id, license_plate, type, vehicle_class, vehicle_group, model, capacity, status, 
			ST_Y(current_location::geometry) as lat, ST_X(current_location::geometry) as lng, 
			last_heartbeat, created_at, updated_at 
		FROM vehicles WHERE id = $1
	`

	row := r.DB.QueryRowContext(ctx, query, id)
	var v models.Vehicle
	var vClass, vGroup, vModel sql.NullString // Handle nulls for existing records

	var lat, lng sql.NullFloat64

	err := row.Scan(
		&v.ID, &v.CompanyID, &v.PlateNumber, &v.Type, &vClass, &vGroup, &vModel, &v.Capacity, &v.Status,
		&lat, &lng,
		&v.LastHeartbeat, &v.CreatedAt, &v.UpdatedAt,
	)

	v.VehicleClass = vClass.String
	v.VehicleGroup = vGroup.String
	v.Model = vModel.String

	if err != nil {
		return nil, err
	}

	if lat.Valid {
		v.CurrentLat = &lat.Float64
	}
	if lng.Valid {
		v.CurrentLng = &lng.Float64
	}

	return &v, nil
}

func (r *PostgresVehicleRepo) ListByCompany(ctx context.Context, companyID uuid.UUID) ([]models.Vehicle, error) {
	query := `
		SELECT 
			id, company_id, license_plate, type, vehicle_class, vehicle_group, model, capacity, status, 
			ST_Y(current_location::geometry) as lat, ST_X(current_location::geometry) as lng, 
			last_heartbeat, created_at, updated_at 
		FROM vehicles WHERE company_id = $1
	`

	rows, err := r.DB.QueryContext(ctx, query, companyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var vehicles []models.Vehicle
	for rows.Next() {
		var v models.Vehicle
		var vClass, vGroup, vModel sql.NullString
		var lat, lng sql.NullFloat64
		if err := rows.Scan(
			&v.ID, &v.CompanyID, &v.PlateNumber, &v.Type, &vClass, &vGroup, &vModel, &v.Capacity, &v.Status,
			&lat, &lng,
			&v.LastHeartbeat, &v.CreatedAt, &v.UpdatedAt,
		); err != nil {
			return nil, err
		}
		v.VehicleClass = vClass.String
		v.VehicleGroup = vGroup.String
		v.Model = vModel.String
		if lat.Valid {
			v.CurrentLat = &lat.Float64
		}
		if lng.Valid {
			v.CurrentLng = &lng.Float64
		}
		vehicles = append(vehicles, v)
	}
	return vehicles, nil
}

func (r *PostgresVehicleRepo) UpdateStatus(ctx context.Context, id uuid.UUID, status models.VehicleStatus) error {
	query := `UPDATE vehicles SET status = $1, updated_at = $2 WHERE id = $3`
	_, err := r.DB.ExecContext(ctx, query, status, time.Now(), id)
	return err
}
