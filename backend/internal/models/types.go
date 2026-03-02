package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// Enums
type CompanyType string

const (
	CompanyTypeDemand  CompanyType = "DEMAND"
	CompanyTypeSupply  CompanyType = "SUPPLY"
	CompanyTypeBoth    CompanyType = "BOTH"
	CompanyTypePartner CompanyType = "PARTNER" // New
)

type UserRole string

const (
	RoleAdmin         UserRole = "ADMIN"
	RoleOps           UserRole = "OPS"
	RoleSupplierAdmin UserRole = "SUPPLIER_ADMIN"
	RoleDispatcher    UserRole = "DISPATCHER"
	RolePartner       UserRole = "PARTNER"
	RoleSuperAdmin    UserRole = "SUPER_ADMIN"
)

type VehicleStatus string

const (
	VehicleStatusIdle        VehicleStatus = "IDLE"
	VehicleStatusOnTrip      VehicleStatus = "ON_TRIP"
	VehicleStatusOffline     VehicleStatus = "OFFLINE"
	VehicleStatusMaintenance VehicleStatus = "MAINTENANCE"
)

type TripStatus string

const (
	TripStatusCreated           TripStatus = "CREATED"
	TripStatusMarketplaceSearch TripStatus = "MARKETPLACE_SEARCH"
	TripStatusOffered           TripStatus = "OFFERED"
	TripStatusAccepted          TripStatus = "ACCEPTED"
	TripStatusOutsourceAssigned TripStatus = "OUTSOURCE_ASSIGNED"
	TripStatusDriverAssigned    TripStatus = "DRIVER_ASSIGNED"
	TripStatusEnRoute           TripStatus = "EN_ROUTE"
	TripStatusInTrip            TripStatus = "IN_TRIP"
	TripStatusCompleted         TripStatus = "COMPLETED"
	TripStatusCancelled         TripStatus = "CANCELLED"
	TripStatusFailed            TripStatus = "FAILED"
)

type TripServiceType string

const (
	ServiceTypeOneWay    TripServiceType = "ONE_WAY"
	ServiceTypeRoundTrip TripServiceType = "ROUND_TRIP"
	ServiceTypeHourly    TripServiceType = "HOURLY"
	ServiceTypePackage   TripServiceType = "PACKAGE"
)

type TripOfferStatus string

const (
	OfferStatusPending  TripOfferStatus = "PENDING"
	OfferStatusAccepted TripOfferStatus = "ACCEPTED"
	OfferStatusRejected TripOfferStatus = "REJECTED"
	OfferStatusTimeout  TripOfferStatus = "TIMEOUT"
)

// Structs

type Company struct {
	ID                 uuid.UUID       `db:"id" json:"id"`
	Name               string          `db:"name" json:"name"`
	Type               CompanyType     `db:"type" json:"type"`
	ContactInfo        json.RawMessage `db:"contact_info" json:"contact_info"` // JSONB
	Verified           bool            `db:"verified" json:"verified"`
	TradeLicenseNo     string          `db:"trade_license_no" json:"trade_license_no"`
	ITCPermitNo        string          `db:"itc_permit_no" json:"itc_permit_no"`
	VATNo              string          `db:"vat_no" json:"vat_no"`
	SLAScore           float64         `db:"sla_score" json:"sla_score"`
	Rating             float64         `db:"rating" json:"rating"`
	VerificationStatus string          `db:"verification_status" json:"verification_status"`
	Settings           json.RawMessage `db:"settings" json:"settings"` // JSONB
	CreatedAt          time.Time       `db:"created_at" json:"created_at"`
	UpdatedAt          time.Time       `db:"updated_at" json:"updated_at"`
}

type User struct {
	ID                     uuid.UUID     `db:"id" json:"id"`
	CompanyID              uuid.NullUUID `db:"company_id" json:"company_id"`
	Role                   UserRole      `db:"role" json:"role"`
	Email                  string        `db:"email" json:"email"`
	PasswordHash           string        `db:"password_hash" json:"-"`
	Name                   string        `db:"name" json:"name"`
	IsSuperAdmin           bool          `db:"is_super_admin" json:"is_super_admin"`
	PasswordChangeRequired bool          `db:"password_change_required" json:"password_change_required"`
	CreatedAt              time.Time     `db:"created_at" json:"created_at"`
}

type Vehicle struct {
	ID           uuid.UUID     `db:"id" json:"id"`
	CompanyID    uuid.UUID     `db:"company_id" json:"company_id"`
	PlateNumber  string        `db:"license_plate" json:"plate_number"` // Remote DB uses license_plate
	Type         string        `db:"type" json:"type"`                  // Deprecated in favor of Group/Model but kept for compat
	VehicleClass string        `db:"vehicle_class" json:"vehicle_class"`
	VehicleGroup string        `db:"vehicle_group" json:"vehicle_group"`
	Model        string        `db:"model" json:"model"`
	Capacity     int           `db:"capacity" json:"capacity"`
	Status       VehicleStatus `db:"status" json:"status"`
	// Basic Lat/Lng for now. PostGIS interaction will happen in Repo.
	CurrentLat      *float64   `db:"current_lat" json:"lat,omitempty"` // Mapped via query usually
	CurrentLng      *float64   `db:"current_lng" json:"lng,omitempty"`
	PermitExpiry    *time.Time `db:"permit_expiry" json:"permit_expiry"`
	InsuranceExpiry *time.Time `db:"insurance_expiry" json:"insurance_expiry"`
	LastHeartbeat   *time.Time `db:"last_heartbeat" json:"last_heartbeat"`
	CreatedAt       time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt       time.Time  `db:"updated_at" json:"updated_at"`
}

type Driver struct {
	ID               uuid.UUID     `db:"id" json:"id"`
	CompanyID        uuid.UUID     `db:"company_id" json:"company_id"`
	Name             string        `db:"name" json:"name"`
	Phone            string        `db:"phone" json:"phone,omitempty"`
	LicenseNumber    string        `db:"license_number" json:"license_number,omitempty"`
	CurrentVehicleID uuid.NullUUID `db:"current_vehicle_id" json:"current_vehicle_id,omitempty"`
	LicenseExpiry    *time.Time    `db:"license_expiry" json:"license_expiry"`
	ITCPermitExpiry  *time.Time    `db:"itc_permit_expiry" json:"itc_permit_expiry"`
	VisaExpiry       *time.Time    `db:"visa_expiry" json:"visa_expiry"`
	CreatedAt        time.Time     `db:"created_at" json:"created_at"`
}

type Trip struct {
	ID                    uuid.UUID  `db:"id" json:"id"`
	RequestingCompanyID   uuid.UUID  `db:"requesting_company_id" json:"requesting_company_id"`
	BookingReference      string     `db:"booking_reference" json:"booking_reference"`
	PickupZone            string     `db:"pickup_zone" json:"pickup_zone"`
	DropoffZone           string     `db:"dropoff_zone" json:"dropoff_zone"`
	PickupTime            time.Time  `db:"pickup_time" json:"pickup_time"`
	Status                TripStatus `db:"status" json:"status"`
	PassengerName         string     `db:"passenger_name" json:"passenger_name"`
	PassengerPhone        string     `db:"passenger_phone" json:"passenger_phone"`
	Price                 float64    `db:"price" json:"price"`
	PaymentMethod         string     `db:"payment_method" json:"payment_method"`
	SupplierName          string     `db:"supplier_name" json:"supplier_name"`
	VehicleID             string     `db:"vehicle_id" json:"vehicle_id"`                           // Simplified for mock
	DriverLinkToken       string     `db:"driver_link_token" json:"driver_link_token,omitempty"`   // New Secure Link Token
	ServiceType           string     `db:"service_type" json:"service_type"`                       // One Way, Hourly, etc.
	RequestedVehicleType  string     `db:"requested_vehicle_type" json:"requested_vehicle_type"`   // e.g. Sedan, SUV
	RequestedVehicleClass string     `db:"requested_vehicle_class" json:"requested_vehicle_class"` // e.g. Standard, Premium
	RequestedVehicleGroup string     `db:"requested_vehicle_group" json:"requested_vehicle_group"` // e.g. Sedan, SUV
	RFQNumber             string     `db:"rfq_number" json:"rfq_number,omitempty"`                 // New Field for RFQ ID
	Visibility            string     `db:"visibility" json:"visibility"`                           // PUBLIC or DIRECT
	PickupLandmark        string     `db:"pickup_landmark" json:"pickup_landmark,omitempty"`
	DropoffLandmark       string     `db:"dropoff_landmark" json:"dropoff_landmark,omitempty"`
	CreatedAt             time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt             time.Time  `db:"updated_at" json:"updated_at"`
}

type TripOffer struct {
	ID                uuid.UUID       `db:"id" json:"id"`
	TripID            uuid.UUID       `db:"trip_id" json:"trip_id"`
	SupplierCompanyID uuid.UUID       `db:"supplier_company_id" json:"supplier_company_id"`
	SupplierName      string          `db:"supplier_name" json:"supplier_name"` // Joined field
	Status            TripOfferStatus `db:"status" json:"status"`
	Score             float64         `db:"score" json:"score"`
	Price             float64         `db:"price" json:"price"` // Not in schema init? Need to check.
	Notes             string          `db:"notes" json:"notes"` // Not in schema init?
	CreatedAt         time.Time       `db:"created_at" json:"created_at"`
	ExpiresAt         time.Time       `db:"expires_at" json:"expires_at"`
	BookingReference  string          `db:"booking_reference" json:"booking_reference"` // Joined field
	RFQNumber         string          `db:"rfq_number" json:"rfq_number,omitempty"`     // Joined field
}

// Search/Quote Types

type QuoteRequest struct {
	PickupLat  float64   `json:"pickup_lat"`
	PickupLng  float64   `json:"pickup_lng"`
	DropoffLat float64   `json:"dropoff_lat"`
	DropoffLng float64   `json:"dropoff_lng"`
	PickupTime time.Time `json:"pickup_time"`
}

type SubmitQuoteRequest struct {
	TripID string  `json:"trip_id"`
	Price  float64 `json:"price"`
	Notes  string  `json:"notes"`
}

// VehicleSearchRequest params for finding B2B farm-out vehicles
type VehicleSearchRequest struct {
	PickupZone   string    `json:"pickup_zone"`
	DropoffZone  string    `json:"dropoff_zone"`
	PickupTime   time.Time `json:"pickup_time"`
	VehicleClass string    `json:"vehicle_class,omitempty"` // Optional filter
	VehicleType  string    `json:"vehicle_type,omitempty"`  // Optional filter
}

// VehicleSearchResult represents an available vehicle from a supplier
type VehicleSearchResult struct {
	VehicleID    string `json:"vehicle_id"`
	Make         string `json:"make"`
	Model        string `json:"model"`
	Year         int    `json:"year"`
	Color        string `json:"color"`
	VehicleType  string `json:"vehicle_type"`  // Sedan, SUV...
	VehicleClass string `json:"vehicle_class"` // Standard, Premium...
	SupplierName string `json:"supplier_name"`
	Location     string `json:"location"` // Current Zone
	Distance     string `json:"distance"` // e.g. "5 km away"
	ImageURL     string `json:"image_url,omitempty"`
}

type BookVehicleRequest struct {
	VehicleID        string    `json:"vehicle_id"`
	SupplierName     string    `json:"supplier_name"`
	BookingReference string    `json:"booking_reference,omitempty"` // New Field
	PickupZone       string    `json:"pickup_zone"`
	DropoffZone      string    `json:"dropoff_zone"`
	PickupTime       time.Time `json:"pickup_time"`
	PassengerName    string    `json:"passenger_name"`
	PassengerPhone   string    `json:"passenger_phone"`
	PaymentMethod    string    `json:"payment_method"` // New Field: Cash, Credit, On Account
	Price            float64   `json:"price"`
	Notes            string    `json:"notes"`
}

type BookVehicleResponse struct {
	TripID  string `json:"trip_id"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

type InvoiceStatus string

const (
	InvoiceStatusPending  InvoiceStatus = "PENDING"
	InvoiceStatusPaid     InvoiceStatus = "PAID"
	InvoiceStatusRejected InvoiceStatus = "REJECTED"
)

type Invoice struct {
	ID                uuid.UUID     `db:"id" json:"id"`
	TripID            uuid.UUID     `db:"trip_id" json:"trip_id"`
	SupplierCompanyID uuid.UUID     `db:"supplier_company_id" json:"supplier_company_id"`
	InvoiceNumber     string        `db:"invoice_number" json:"invoice_number"`
	Amount            float64       `db:"amount" json:"amount"`
	PlatformFee       float64       `db:"platform_fee" json:"platform_fee"`
	VATAmount         float64       `db:"vat_amount" json:"vat_amount"`
	NetPayout         float64       `db:"net_payout" json:"net_payout"`
	Status            InvoiceStatus `db:"status" json:"status"`
	CreatedAt         time.Time     `db:"created_at" json:"created_at"`
	UpdatedAt         time.Time     `db:"updated_at" json:"updated_at"`

	// Join Fields
	BookingReference string `db:"booking_reference" json:"booking_reference,omitempty"`
}

type SubmitInvoiceRequest struct {
	TripID        string  `json:"trip_id"`
	InvoiceNumber string  `json:"invoice_number"`
	Amount        float64 `json:"amount"`
}
