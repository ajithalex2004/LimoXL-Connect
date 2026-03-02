package api

import (
	"encoding/json"
	"fmt"
	"limoxlink-backend/internal/middleware"
	"limoxlink-backend/internal/models"
	"limoxlink-backend/internal/repository"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

type FleetHandler struct {
	VehicleRepo repository.VehicleRepository
	DriverRepo  repository.DriverRepository
}

func NewFleetHandler(vRepo repository.VehicleRepository, dRepo repository.DriverRepository) *FleetHandler {
	return &FleetHandler{
		VehicleRepo: vRepo,
		DriverRepo:  dRepo,
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
