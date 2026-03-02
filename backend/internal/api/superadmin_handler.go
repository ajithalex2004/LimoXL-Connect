package api

import (
	"encoding/json"
	"limoxlink-backend/internal/models"
	"limoxlink-backend/internal/repository"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
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

// HandleCreateTenant creates a new tenant
func (h *SuperAdminHandler) HandleCreateTenant(w http.ResponseWriter, r *http.Request) {
	var req models.Tenant
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}

	if err := h.TenantRepo.Create(r.Context(), &req); err != nil {
		http.Error(w, "Failed to create tenant", http.StatusInternalServerError)
		return
	}

	// Initialize with all features disabled by default
	for _, key := range models.AllFeatures {
		h.TenantRepo.UpdateFeature(r.Context(), req.ID, key, false)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(req)
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
