package api

import (
	"encoding/json"
	"limoxlink-backend/internal/models"
	"limoxlink-backend/internal/repository"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type NUIMasterHandler struct {
	MasterRepo repository.NUIMasterRepository
}

func (h *NUIMasterHandler) Create(w http.ResponseWriter, r *http.Request) {
	companyID, ok := getCompanyIDFromClaims(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var master models.NUIMaster
	if err := json.NewDecoder(r.Body).Decode(&master); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	master.CompanyID = companyID
	if err := h.MasterRepo.Create(r.Context(), &master); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(master)
}

func (h *NUIMasterHandler) List(w http.ResponseWriter, r *http.Request) {
	companyID, ok := getCompanyIDFromClaims(r)
	if !ok {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]interface{}{})
		return
	}

	category := r.URL.Query().Get("category")
	masters, err := h.MasterRepo.ListByCompany(r.Context(), companyID, category)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if masters == nil {
		masters = []models.NUIMaster{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(masters)
}

func (h *NUIMasterHandler) Update(w http.ResponseWriter, r *http.Request) {
	companyID, ok := getCompanyIDFromClaims(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var master models.NUIMaster
	if err := json.NewDecoder(r.Body).Decode(&master); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	master.ID = id
	master.CompanyID = companyID
	if err := h.MasterRepo.Update(r.Context(), &master); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(master)
}

func (h *NUIMasterHandler) Delete(w http.ResponseWriter, r *http.Request) {
	_, ok := getCompanyIDFromClaims(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if err := h.MasterRepo.Delete(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
