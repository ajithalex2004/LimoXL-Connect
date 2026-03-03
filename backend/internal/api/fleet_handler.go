package api

import (
	"encoding/json"
	"fmt"
	"limoxlink-backend/internal/middleware"
	"limoxlink-backend/internal/models"
	"limoxlink-backend/internal/repository"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type FleetHandler struct {
	VehicleRepo    repository.VehicleRepository
	DriverRepo     repository.DriverRepository
	AttachmentRepo repository.FleetAttachmentRepository
}

func NewFleetHandler(vRepo repository.VehicleRepository, dRepo repository.DriverRepository, aRepo repository.FleetAttachmentRepository) *FleetHandler {
	return &FleetHandler{
		VehicleRepo:    vRepo,
		DriverRepo:     dRepo,
		AttachmentRepo: aRepo,
	}
}

// getCompanyIDFromClaims extracts the company UUID from the authenticated JWT claims.
func getCompanyIDFromClaims(r *http.Request) (uuid.UUID, bool) {
	claims, ok := r.Context().Value(middleware.ClaimsKey).(*middleware.Claims)
	if !ok || claims == nil {
		return uuid.Nil, false
	}

	if claims.IsSuperAdmin {
		return uuid.Nil, true // SuperAdmin has no specific company, but is authorized
	}

	if claims.CompanyID == "" {
		return uuid.Nil, false
	}

	id, err := uuid.Parse(claims.CompanyID)
	if err != nil {
		return uuid.Nil, false
	}
	return id, true
}

func (h *FleetHandler) ListVehicles(w http.ResponseWriter, r *http.Request) {
	companyID, ok := getCompanyIDFromClaims(r)
	if !ok {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]interface{}{})
		return
	}

	vehicles, err := h.VehicleRepo.ListByCompany(r.Context(), companyID)
	if err != nil {
		http.Error(w, "Error listing vehicles", http.StatusInternalServerError)
		return
	}
	if vehicles == nil {
		vehicles = []models.Vehicle{}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(vehicles)
}

func (h *FleetHandler) CreateVehicle(w http.ResponseWriter, r *http.Request) {
	companyID, ok := getCompanyIDFromClaims(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req models.Vehicle
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Backward Compatibility: Map VehicleGroup to Type (DB Enum)
	if req.Type == "" {
		switch req.VehicleGroup {
		case "Sedan":
			req.Type = "SEDAN"
		case "SUV":
			req.Type = "SUV"
		case "Van":
			req.Type = "LUXURY_VAN"
		case "Bus":
			req.Type = "BUS"
		default:
			req.Type = "SEDAN"
		}
	}

	fmt.Printf("DEBUG: Creating Vehicle: %+v\n", req)

	req.CompanyID = companyID
	req.Status = models.VehicleStatusOffline

	if err := h.VehicleRepo.Create(r.Context(), &req); err != nil {
		fmt.Printf("DEBUG: DB Error: %v\n", err)
		if strings.Contains(err.Error(), "duplicate key") {
			http.Error(w, "Vehicle with this license plate already exists", http.StatusConflict)
			return
		}
		http.Error(w, "Error creating vehicle: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(req)
}

func (h *FleetHandler) UpdateVehicle(w http.ResponseWriter, r *http.Request) {
	companyID, ok := getCompanyIDFromClaims(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var v models.Vehicle
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid vehicle ID", http.StatusBadRequest)
		return
	}
	v.ID = id
	v.CompanyID = companyID

	if err := h.VehicleRepo.Update(r.Context(), &v); err != nil {
		http.Error(w, "Error updating vehicle: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}

func (h *FleetHandler) DeleteVehicle(w http.ResponseWriter, r *http.Request) {
	_, ok := getCompanyIDFromClaims(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid vehicle ID", http.StatusBadRequest)
		return
	}

	if err := h.VehicleRepo.Delete(r.Context(), id); err != nil {
		http.Error(w, "Error deleting vehicle", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *FleetHandler) ListDrivers(w http.ResponseWriter, r *http.Request) {
	companyID, ok := getCompanyIDFromClaims(r)
	if !ok {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]interface{}{})
		return
	}

	drivers, err := h.DriverRepo.ListByCompany(r.Context(), companyID)
	if err != nil {
		http.Error(w, "Error listing drivers", http.StatusInternalServerError)
		return
	}
	if drivers == nil {
		drivers = []models.Driver{}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(drivers)
}

func (h *FleetHandler) CreateDriver(w http.ResponseWriter, r *http.Request) {
	companyID, ok := getCompanyIDFromClaims(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req models.Driver
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	req.CompanyID = companyID

	if err := h.DriverRepo.Create(r.Context(), &req); err != nil {
		http.Error(w, "Error creating driver: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(req)
}

func (h *FleetHandler) UpdateDriver(w http.ResponseWriter, r *http.Request) {
	companyID, ok := getCompanyIDFromClaims(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var d models.Driver
	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid driver ID", http.StatusBadRequest)
		return
	}
	d.ID = id
	d.CompanyID = companyID

	if err := h.DriverRepo.Update(r.Context(), &d); err != nil {
		http.Error(w, "Error updating driver: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(d)
}

func (h *FleetHandler) DeleteDriver(w http.ResponseWriter, r *http.Request) {
	_, ok := getCompanyIDFromClaims(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid driver ID", http.StatusBadRequest)
		return
	}

	if err := h.DriverRepo.Delete(r.Context(), id); err != nil {
		http.Error(w, "Error deleting driver", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *FleetHandler) CreateAttachment(w http.ResponseWriter, r *http.Request) {
	var att models.FleetAttachment
	if err := json.NewDecoder(r.Body).Decode(&att); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.AttachmentRepo.Create(r.Context(), &att); err != nil {
		http.Error(w, "Error creating attachment", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(att)
}

func (h *FleetHandler) ListAttachments(w http.ResponseWriter, r *http.Request) {
	entityIDStr := r.URL.Query().Get("entity_id")
	entityType := r.URL.Query().Get("entity_type")

	if entityIDStr == "" || entityType == "" {
		http.Error(w, "entity_id and entity_type are required", http.StatusBadRequest)
		return
	}

	entityID, err := uuid.Parse(entityIDStr)
	if err != nil {
		http.Error(w, "Invalid entity ID", http.StatusBadRequest)
		return
	}

	attachments, err := h.AttachmentRepo.ListByEntity(r.Context(), entityID, entityType)
	if err != nil {
		http.Error(w, "Error listing attachments", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(attachments)
}

func (h *FleetHandler) DeleteAttachment(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid attachment ID", http.StatusBadRequest)
		return
	}

	if err := h.AttachmentRepo.Delete(r.Context(), id); err != nil {
		http.Error(w, "Error deleting attachment", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
