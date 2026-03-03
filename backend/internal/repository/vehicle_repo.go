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

func (r *PostgresVehicleRepo) Create(ctx context.Context, v *models.Vehicle) error {
	query := `
		INSERT INTO vehicles (
			id, company_id, license_plate, type, vehicle_class, vehicle_group, 
			model, capacity, status, permit_expiry, insurance_expiry, 
			chassis_no, vin, year_of_manufacture, color, registration_number,
			plate_code, plate_category, emirate, hierarchy, vehicle_usage,
			created_at, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23)
	`
	v.ID = uuid.New()
	v.CreatedAt = time.Now()
	v.UpdatedAt = time.Now()

	_, err := r.DB.ExecContext(ctx, query,
		v.ID, v.CompanyID, v.PlateNumber, v.Type, v.VehicleClass, v.VehicleGroup,
		v.Model, v.Capacity, v.Status, v.PermitExpiry, v.InsuranceExpiry,
		v.ChassisNo, v.VIN, v.YearOfManufacture, v.Color, v.RegistrationNumber,
		v.PlateCode, v.PlateCategory, v.Emirate, v.Hierarchy, v.VehicleUsage,
		v.CreatedAt, v.UpdatedAt,
	)
	return err
}

func (r *PostgresVehicleRepo) GetByID(ctx context.Context, id uuid.UUID) (*models.Vehicle, error) {
	query := `
		SELECT 
			id, company_id, license_plate, type, vehicle_class, vehicle_group, model, capacity, status, 
			permit_expiry, insurance_expiry, last_heartbeat, created_at, updated_at,
			chassis_no, vin, year_of_manufacture, color, registration_number,
			plate_code, plate_category, emirate, hierarchy, vehicle_usage,
			ST_Y(current_location::geometry) as lat, ST_X(current_location::geometry) as lng
		FROM vehicles WHERE id = $1
	`

	row := r.DB.QueryRowContext(ctx, query, id)
	var v models.Vehicle
	var vClass, vGroup, vModel, chassis, vin, color, reg, pCode, pCat, emirate, hierarchy, usage sql.NullString
	var lat, lng sql.NullFloat64

	err := row.Scan(
		&v.ID, &v.CompanyID, &v.PlateNumber, &v.Type, &vClass, &vGroup, &vModel, &v.Capacity, &v.Status,
		&v.PermitExpiry, &v.InsuranceExpiry, &v.LastHeartbeat, &v.CreatedAt, &v.UpdatedAt,
		&chassis, &vin, &v.YearOfManufacture, &color, &reg, &pCode, &pCat, &emirate, &hierarchy, &usage,
		&lat, &lng,
	)
	if err != nil {
		return nil, err
	}

	v.VehicleClass = vClass.String
	v.VehicleGroup = vGroup.String
	v.Model = vModel.String
	v.ChassisNo = chassis.String
	v.VIN = vin.String
	v.Color = color.String
	v.RegistrationNumber = reg.String
	v.PlateCode = pCode.String
	v.PlateCategory = pCat.String
	v.Emirate = emirate.String
	v.Hierarchy = hierarchy.String
	v.VehicleUsage = usage.String

	if lat.Valid { v.CurrentLat = &lat.Float64 }
	if lng.Valid { v.CurrentLng = &lng.Float64 }

	return &v, nil
}

func (r *PostgresVehicleRepo) ListByCompany(ctx context.Context, companyID uuid.UUID) ([]models.Vehicle, error) {
	query := `
		SELECT 
			id, company_id, license_plate, type, vehicle_class, vehicle_group, model, capacity, status, 
			permit_expiry, insurance_expiry, last_heartbeat, created_at, updated_at,
			chassis_no, vin, year_of_manufacture, color, registration_number,
			plate_code, plate_category, emirate, hierarchy, vehicle_usage,
			ST_Y(current_location::geometry) as lat, ST_X(current_location::geometry) as lng
		FROM vehicles WHERE company_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.DB.QueryContext(ctx, query, companyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var vehicles []models.Vehicle
	for rows.Next() {
		var v models.Vehicle
		var vClass, vGroup, vModel, chassis, vin, color, reg, pCode, pCat, emirate, hierarchy, usage sql.NullString
		var lat, lng sql.NullFloat64
		err := rows.Scan(
			&v.ID, &v.CompanyID, &v.PlateNumber, &v.Type, &vClass, &vGroup, &vModel, &v.Capacity, &v.Status,
			&v.PermitExpiry, &v.InsuranceExpiry, &v.LastHeartbeat, &v.CreatedAt, &v.UpdatedAt,
			&chassis, &vin, &v.YearOfManufacture, &color, &reg, &pCode, &pCat, &emirate, &hierarchy, &usage,
			&lat, &lng,
		)
		if err != nil {
			return nil, err
		}
		v.VehicleClass = vClass.String
		v.VehicleGroup = vGroup.String
		v.Model = vModel.String
		v.ChassisNo = chassis.String
		v.VIN = vin.String
		v.Color = color.String
		v.RegistrationNumber = reg.String
		v.PlateCode = pCode.String
		v.PlateCategory = pCat.String
		v.Emirate = emirate.String
		v.Hierarchy = hierarchy.String
		v.VehicleUsage = usage.String
		if lat.Valid { v.CurrentLat = &lat.Float64 }
		if lng.Valid { v.CurrentLng = &lng.Float64 }
		vehicles = append(vehicles, v)
	}
	return vehicles, nil
}

func (r *PostgresVehicleRepo) Update(ctx context.Context, v *models.Vehicle) error {
	query := `
		UPDATE vehicles SET
			license_plate = $2, type = $3, vehicle_class = $4, vehicle_group = $5, 
			model = $6, capacity = $7, status = $8, permit_expiry = $9, insurance_expiry = $10, 
			chassis_no = $11, vin = $12, year_of_manufacture = $13, color = $14, registration_number = $15,
			plate_code = $16, plate_category = $17, emirate = $18, hierarchy = $19, vehicle_usage = $20,
			updated_at = $21
		WHERE id = $1
	`
	v.UpdatedAt = time.Now()
	_, err := r.DB.ExecContext(ctx, query,
		v.ID, v.PlateNumber, v.Type, v.VehicleClass, v.VehicleGroup,
		v.Model, v.Capacity, v.Status, v.PermitExpiry, v.InsuranceExpiry,
		v.ChassisNo, v.VIN, v.YearOfManufacture, v.Color, v.RegistrationNumber,
		v.PlateCode, v.PlateCategory, v.Emirate, v.Hierarchy, v.VehicleUsage,
		v.UpdatedAt,
	)
	return err
}

func (r *PostgresVehicleRepo) UpdateStatus(ctx context.Context, id uuid.UUID, status models.VehicleStatus) error {
	query := `UPDATE vehicles SET status = $1, updated_at = $2 WHERE id = $3`
	_, err := r.DB.ExecContext(ctx, query, status, time.Now(), id)
	return err
}

func (r *PostgresVehicleRepo) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM vehicles WHERE id = $1`
	_, err := r.DB.ExecContext(ctx, query, id)
	return err
}
