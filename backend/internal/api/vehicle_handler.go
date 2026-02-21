package api

import (
	"encoding/json"
	"fmt"
	"limoxlink-backend/internal/models"
	"limoxlink-backend/internal/repository"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type VehicleHandler struct {
	Repo repository.VehicleRepository
}

func NewVehicleHandler(repo repository.VehicleRepository) *VehicleHandler {
	return &VehicleHandler{Repo: repo}
}

func (h *VehicleHandler) CreateVehicle(w http.ResponseWriter, r *http.Request) {
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

	if err := h.Repo.Create(r.Context(), &req); err != nil {
		http.Error(w, "Error creating vehicle: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(req)
}

func (h *VehicleHandler) GetVehicle(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid vehicle ID", http.StatusBadRequest)
		return
	}

	vehicle, err := h.Repo.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, "Vehicle not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(vehicle)
}

func (h *VehicleHandler) ListVehicles(w http.ResponseWriter, r *http.Request) {
	// For MVP, if no company_id is passed, maybe return all?
	// Or require company_id. Let's look for query param.
	companyIDStr := r.URL.Query().Get("company_id")
	if companyIDStr == "" {
		http.Error(w, "company_id query param is required", http.StatusBadRequest)
		return
	}

	companyID, err := uuid.Parse(companyIDStr)
	if err != nil {
		http.Error(w, "Invalid company ID", http.StatusBadRequest)
		return
	}

	vehicles, err := h.Repo.ListByCompany(r.Context(), companyID)
	if err != nil {
		fmt.Printf("DEBUG: ListVehicles error: %v\n", err)
		http.Error(w, "DB_ERR: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(vehicles)
}
