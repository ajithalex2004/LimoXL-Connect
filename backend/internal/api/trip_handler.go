package api

import (
	"encoding/json"
	"limoxlink-backend/internal/models"
	"limoxlink-backend/internal/repository"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type TripHandler struct {
	Repo repository.TripRepository
}

func NewTripHandler(repo repository.TripRepository) *TripHandler {
	return &TripHandler{Repo: repo}
}

// HandleSearchVehicles handles POST /api/marketplace/search-vehicles
func (h *TripHandler) HandleSearchVehicles(w http.ResponseWriter, r *http.Request) {
	var req models.VehicleSearchRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	results, err := h.Repo.SearchAvailableVehicles(r.Context(), req)
	if err != nil {
		http.Error(w, "Failed to search vehicles", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

// HandleBookVehicle handles POST /api/marketplace/book-vehicle
func (h *TripHandler) HandleBookVehicle(w http.ResponseWriter, r *http.Request) {
	var req models.BookVehicleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	response, err := h.Repo.BookVehicle(r.Context(), req)
	if err != nil {
		http.Error(w, "Failed to book vehicle", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// HandleListTrips handles GET /api/trips
func (h *TripHandler) HandleListTrips(w http.ResponseWriter, r *http.Request) {
	trips, err := h.Repo.ListTrips(r.Context())
	if err != nil {
		http.Error(w, "Failed to list trips", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(trips)
}

// HandleGetTrip handles GET /api/trips/{id}
func (h *TripHandler) HandleGetTrip(w http.ResponseWriter, r *http.Request) {
	tripIDStr := chi.URLParam(r, "id")
	tripID, err := uuid.Parse(tripIDStr)
	if err != nil {
		http.Error(w, "Invalid Trip ID", http.StatusBadRequest)
		return
	}

	trip, err := h.Repo.GetTrip(r.Context(), tripID)
	if err != nil {
		http.Error(w, "Trip not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(trip)
}
