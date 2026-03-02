package api

import (
	"encoding/json"
	"limoxlink-backend/internal/middleware"
	"limoxlink-backend/internal/models"
	"limoxlink-backend/internal/repository"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type SuperAdminHandler struct {
	TenantRepo *repository.TenantRepository
	UserRepo   *repository.PostgresUserRepo
}

func NewSuperAdminHandler(tenantRepo *repository.TenantRepository, userRepo *repository.PostgresUserRepo) *SuperAdminHandler {
	return &SuperAdminHandler{
		TenantRepo: tenantRepo,
		UserRepo:   userRepo,
	}
}

// HandleListTenants returns all tenants with their stats
func (h *SuperAdminHandler) HandleListTenants(w http.ResponseWriter, r *http.Request) {
	tenants, err := h.TenantRepo.ListAll(r.Context())
	if err != nil {
		http.Error(w, "Failed to list tenants", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tenants)
}

// HandleCreateTenant creates a new tenant and an associated company
func (h *SuperAdminHandler) HandleCreateTenant(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name string `json:"name"`
		Slug string `json:"slug"`
		Plan string `json:"plan"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}

	// 1. Create a Company for this tenant first
	var companyID uuid.UUID
	err := h.TenantRepo.DB.QueryRowContext(r.Context(), `
		INSERT INTO companies (name, type, verified)
		VALUES ($1, 'DEMAND', true)
		RETURNING id
	`, req.Name).Scan(&companyID)
	if err != nil {
		http.Error(w, "Failed to create company for tenant: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 2. Create the Tenant
	tenant := &models.Tenant{
		CompanyID: &companyID,
		Name:      req.Name,
		Slug:      req.Slug,
		Status:    "ACTIVE",
		Plan:      models.TenantPlan(req.Plan),
		MaxUsers:  10,
	}

	if err := h.TenantRepo.Create(r.Context(), tenant); err != nil {
		http.Error(w, "Failed to create tenant: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 3. Initialize features
	for _, key := range models.AllFeatures {
		h.TenantRepo.UpdateFeature(r.Context(), tenant.ID, key, true)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(tenant)
}

// HandleUpdateTenant updates tenant details
func (h *SuperAdminHandler) HandleUpdateTenant(w http.ResponseWriter, r *http.Request) {
	tenantIDStr := chi.URLParam(r, "id")
	tenantID, err := uuid.Parse(tenantIDStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var req models.Tenant
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}
	req.ID = tenantID

	if err := h.TenantRepo.Update(r.Context(), &req); err != nil {
		http.Error(w, "Failed to update tenant", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(req)
}

// HandleToggleFeature enables/disables a feature for a tenant
func (h *SuperAdminHandler) HandleToggleFeature(w http.ResponseWriter, r *http.Request) {
	tenantIDStr := chi.URLParam(r, "id")
	tenantID, err := uuid.Parse(tenantIDStr)
	if err != nil {
		http.Error(w, "Invalid Tenant ID", http.StatusBadRequest)
		return
	}

	var req struct {
		FeatureKey string `json:"feature_key"`
		IsEnabled  bool   `json:"is_enabled"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}

	if err := h.TenantRepo.UpdateFeature(r.Context(), tenantID, req.FeatureKey, req.IsEnabled); err != nil {
		http.Error(w, "Failed to update feature", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "Feature Updated"})
}
// HandleSwitchTenant generates a new JWT token impersonating a tenant's context
func (h *SuperAdminHandler) HandleSwitchTenant(w http.ResponseWriter, r *http.Request) {
	tenantIDStr := chi.URLParam(r, "id")
	tenantID, err := uuid.Parse(tenantIDStr)
	if err != nil {
		http.Error(w, "Invalid Tenant ID", http.StatusBadRequest)
		return
	}

	// Verify tenant exists and get company ID
	var companyID string
	var role string = "ADMIN" // Default to Admin role within the tenant
	err = h.TenantRepo.DB.QueryRowContext(r.Context(), "SELECT company_id FROM tenants WHERE id = $1", tenantID).Scan(&companyID)
	if err != nil {
		http.Error(w, "Tenant not found", http.StatusNotFound)
		return
	}

	claims, _ := r.Context().Value(middleware.ClaimsKey).(*middleware.Claims)

	// Generate new token with TenantID and CompanyID
	token, err := middleware.GenerateToken(claims.UserID, companyID, tenantID.String(), role, true)
	if err != nil {
		http.Error(w, "Failed to generate impersonation token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"token": token,
	})
}

// HandleDeleteTenant removes a tenant
func (h *SuperAdminHandler) HandleDeleteTenant(w http.ResponseWriter, r *http.Request) {
	tenantIDStr := chi.URLParam(r, "id")
	tenantID, err := uuid.Parse(tenantIDStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if err := h.TenantRepo.Delete(r.Context(), tenantID); err != nil {
		http.Error(w, "Failed to delete tenant", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// HandleCreateTenantAdmin creates an admin user for a specific tenant
func (h *SuperAdminHandler) HandleCreateTenantAdmin(w http.ResponseWriter, r *http.Request) {
	tenantIDStr := chi.URLParam(r, "id")
	tenantID, err := uuid.Parse(tenantIDStr)
	if err != nil {
		http.Error(w, "Invalid Tenant ID", http.StatusBadRequest)
		return
	}

	var req struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}

	// 1. Get CompanyID from Tenant
	var companyID uuid.UUID
	err = h.TenantRepo.DB.QueryRowContext(r.Context(), "SELECT company_id FROM tenants WHERE id = $1", tenantID).Scan(&companyID)
	if err != nil {
		http.Error(w, "Tenant not found or has no linked company", http.StatusNotFound)
		return
	}

	// 2. Hash Password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	// 3. Create User
	user := &models.User{
		CompanyID:              uuid.NullUUID{UUID: companyID, Valid: true},
		Name:                   req.Name,
		Email:                  req.Email,
		PasswordHash:           string(hashedPassword),
		Role:                   models.RoleAdmin,
		IsSuperAdmin:           false,
		PasswordChangeRequired: true,
	}

	if err := h.UserRepo.Create(r.Context(), user); err != nil {
		http.Error(w, "Failed to create admin user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}
