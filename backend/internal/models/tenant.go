package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// TenantStatus represents the status of a tenant
type TenantStatus string

const (
	TenantStatusActive    TenantStatus = "ACTIVE"
	TenantStatusSuspended TenantStatus = "SUSPENDED"
	TenantStatusTrial     TenantStatus = "TRIAL"
)

// TenantPlan represents the subscription plan of a tenant
type TenantPlan string

const (
	TenantPlanStarter      TenantPlan = "STARTER"
	TenantPlanProfessional TenantPlan = "PROFESSIONAL"
	TenantPlanEnterprise   TenantPlan = "ENTERPRISE"
)

// Available feature keys for tenant microservice gating
const (
	FeatureDispatch             = "dispatch"
	FeatureOutsourceMarketplace = "outsource_marketplace"
	FeatureFleetManagement      = "fleet_management"
	FeatureTeamManagement       = "team_management"
	FeatureInvoicing            = "invoicing"
	FeaturePartnerPortal        = "partner_portal"
	FeatureAnalytics            = "analytics"
)

// AllFeatures lists every feature key that can be assigned per tenant
var AllFeatures = []string{
	FeatureDispatch,
	FeatureOutsourceMarketplace,
	FeatureFleetManagement,
	FeatureTeamManagement,
	FeatureInvoicing,
	FeaturePartnerPortal,
	FeatureAnalytics,
}

// Tenant represents an operator company as a tenant
type Tenant struct {
	ID        uuid.UUID    `db:"id" json:"id"`
	CompanyID *uuid.UUID   `db:"company_id" json:"company_id"`
	Name      string       `db:"name" json:"name"`
	Slug      string       `db:"slug" json:"slug"`
	Status    TenantStatus `db:"status" json:"status"`
	Plan      TenantPlan   `db:"plan" json:"plan"`
	MaxUsers  int          `db:"max_users" json:"max_users"`
	CreatedAt time.Time    `db:"created_at" json:"created_at"`
	UpdatedAt time.Time    `db:"updated_at" json:"updated_at"`
}

// TenantFeature represents a single feature flag for a tenant
type TenantFeature struct {
	ID         uuid.UUID       `db:"id" json:"id"`
	TenantID   uuid.UUID       `db:"tenant_id" json:"tenant_id"`
	FeatureKey string          `db:"feature_key" json:"feature_key"`
	IsEnabled  bool            `db:"is_enabled" json:"is_enabled"`
	Config     json.RawMessage `db:"config" json:"config"`
	CreatedAt  time.Time       `db:"created_at" json:"created_at"`
	UpdatedAt  time.Time       `db:"updated_at" json:"updated_at"`
}

// TenantWithFeatures combines a tenant with its feature flags and stats
type TenantWithFeatures struct {
	Tenant
	Features  []TenantFeature `json:"features"`
	UserCount int             `json:"user_count"`
	TripCount int             `json:"trip_count"`

	// Company name for display
	CompanyName string `json:"company_name,omitempty"`
}

// FeatureMap returns a map of feature_key -> is_enabled for quick lookup
func (t *TenantWithFeatures) FeatureMap() map[string]bool {
	m := make(map[string]bool)
	for _, f := range t.Features {
		m[f.FeatureKey] = f.IsEnabled
	}
	return m
}
