package repository

import (
	"context"
	"limoxlink-backend/internal/models"

	"github.com/google/uuid"
)

type CompanyRepository interface {
	Create(ctx context.Context, company *models.Company) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Company, error)
	List(ctx context.Context) ([]models.Company, error)
	// Add Update/Delete as needed
}

type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	ListByCompany(ctx context.Context, companyID uuid.UUID) ([]*models.User, error)
	UpdatePassword(ctx context.Context, id uuid.UUID, hash string) error
}

type VehicleRepository interface {
	Create(ctx context.Context, vehicle *models.Vehicle) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Vehicle, error)
	ListByCompany(ctx context.Context, companyID uuid.UUID) ([]models.Vehicle, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status models.VehicleStatus) error
}
