package api

import (
	"encoding/json"
	"fmt"
	"limoxlink-backend/internal/models"
	"limoxlink-backend/internal/repository"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

// ... (existing code)

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

func (h *FleetHandler) getCompanyID(r *http.Request) uuid.UUID {
	// Simple path-based logic for demo.
	// In real app, this comes from Middleware/JWT context.
	if r.URL.Path == "/api/operator/vehicles" || r.URL.Path == "/api/operator/drivers" {
		return uuid.MustParse("a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a00") // Operator ID
	}
	return uuid.MustParse("a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11") // Default to Partner ID
}

func (h *FleetHandler) ListVehicles(w http.ResponseWriter, r *http.Request) {
	companyID := h.getCompanyID(r)

	vehicles, err := h.VehicleRepo.ListByCompany(r.Context(), companyID)
	if err != nil {
		http.Error(w, "Error listing vehicles", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(vehicles)
}

func (h *FleetHandler) CreateVehicle(w http.ResponseWriter, r *http.Request) {
	companyID := h.getCompanyID(r)

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
			req.Type = "SEDAN" // Default fallback
		}
	}

	fmt.Printf("DEBUG: Creating Vehicle: %+v\n", req) // Debug Log

	req.CompanyID = companyID
	req.Status = models.VehicleStatusOffline // Default status

	if err := h.VehicleRepo.Create(r.Context(), &req); err != nil {
		fmt.Printf("DEBUG: DB Error: %v\n", err) // Debug Log
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
	companyID := h.getCompanyID(r)

	drivers, err := h.DriverRepo.ListByCompany(r.Context(), companyID)
	if err != nil {
		http.Error(w, "Error listing drivers", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(drivers)
}

func (h *FleetHandler) CreateDriver(w http.ResponseWriter, r *http.Request) {
	companyID := h.getCompanyID(r)

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
