package api

import (
	"encoding/json"
	"limoxlink-backend/internal/models"
	"limoxlink-backend/internal/repository"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type SecureLinkHandler struct {
	Repo repository.TripRepository
}

func NewSecureLinkHandler(repo repository.TripRepository) *SecureLinkHandler {
	return &SecureLinkHandler{Repo: repo}
}

func (h *SecureLinkHandler) GetTripStatus(w http.ResponseWriter, r *http.Request) {
	token := chi.URLParam(r, "token")
	if token == "" {
		http.Error(w, "Token required", http.StatusBadRequest)
		return
	}

	trip, err := h.Repo.GetTripByToken(r.Context(), token)
	if err != nil {
		http.Error(w, "Invalid link or expired", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(trip)
}

func (h *SecureLinkHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	token := chi.URLParam(r, "token")
	if token == "" {
		http.Error(w, "Token required", http.StatusBadRequest)
		return
	}

	// First validate token
	trip, err := h.Repo.GetTripByToken(r.Context(), token)
	if err != nil {
		http.Error(w, "Invalid link", http.StatusForbidden)
		return
	}

	var req struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}

	newStatus := models.TripStatus(req.Status)
	// Validate status transition here ideally

	if err := h.Repo.UpdateTripStatus(r.Context(), trip.ID, newStatus); err != nil {
		http.Error(w, "Failed to update status", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status":     "updated",
		"new_status": string(newStatus),
	})
}
