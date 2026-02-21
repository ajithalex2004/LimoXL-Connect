package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"limoxlink-backend/internal/middleware"
	"limoxlink-backend/internal/models"
	"limoxlink-backend/internal/repository"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type PartnerHandler struct {
	Repo repository.TripRepository
}

func NewPartnerHandler(repo repository.TripRepository) *PartnerHandler {
	return &PartnerHandler{Repo: repo}
}

func (h *PartnerHandler) getPartnerID(r *http.Request) (uuid.UUID, error) {
	claims, ok := r.Context().Value("user").(*middleware.Claims)
	if !ok {
		// Fallback for dev/demo if auth middleware is not applied strictly
		// or if we want to allow non-authed access (not recommended for production)
		// return uuid.MustParse("a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11"), nil
		return uuid.Nil, errors.New("unauthorized")
	}
	return uuid.Parse(claims.CompanyID)
}

func (h *PartnerHandler) ListRFQs(w http.ResponseWriter, r *http.Request) {
	partnerID, err := h.getPartnerID(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	trips, err := h.Repo.ListOpenRFQs(r.Context(), partnerID)
	if err != nil {
		http.Error(w, "Error listing RFQs", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(trips)
}

func (h *PartnerHandler) HandleListRFQHistory(w http.ResponseWriter, r *http.Request) {
	partnerID, err := h.getPartnerID(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	trips, err := h.Repo.ListRFQHistory(r.Context(), partnerID)
	if err != nil {
		http.Error(w, "Error listing RFQ history", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(trips)
}

func (h *PartnerHandler) ListAssignedTrips(w http.ResponseWriter, r *http.Request) {
	partnerID, err := h.getPartnerID(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	trips, err := h.Repo.ListPartnerTrips(r.Context(), partnerID)
	if err != nil {
		http.Error(w, "Error listing assigned trips", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(trips)
}

func (h *PartnerHandler) HandleSubmitQuote(w http.ResponseWriter, r *http.Request) {
	var req models.SubmitQuoteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	supplierID, err := h.getPartnerID(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if err := h.Repo.SubmitQuote(r.Context(), req, supplierID); err != nil {
		http.Error(w, "Error submitting quote", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"status": "Quote Submitted"})
}

type AssignDriverRequest struct {
	TripID       string `json:"trip_id"`
	DriverName   string `json:"driver_name"`
	DriverPhone  string `json:"driver_phone"`
	VehicleModel string `json:"vehicle_model"`
	VehiclePlate string `json:"vehicle_plate"`
}

func (h *PartnerHandler) HandleAssignDriver(w http.ResponseWriter, r *http.Request) {
	var req AssignDriverRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	tripID, err := uuid.Parse(req.TripID)
	if err != nil {
		http.Error(w, "Invalid Trip ID", http.StatusBadRequest)
		return
	}

	// Ideally we check if this partner owns this trip before assigning
	// But skipping for now

	if err := h.Repo.AssignDriver(r.Context(), tripID, req.DriverName, req.DriverPhone, req.VehicleModel, req.VehiclePlate); err != nil {
		http.Error(w, "Error assigning driver", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "Driver Assigned", "message": "Link generated successfully"})
}

func (h *PartnerHandler) HandleAcceptRFQ(w http.ResponseWriter, r *http.Request) {
	type AcceptRequest struct {
		TripID string `json:"trip_id"`
	}
	var req AcceptRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	tripID, err := uuid.Parse(req.TripID)
	if err != nil {
		http.Error(w, "Invalid Trip ID", http.StatusBadRequest)
		return
	}

	partnerID, err := h.getPartnerID(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if err := h.Repo.AcceptRFQ(r.Context(), tripID, partnerID); err != nil {
		http.Error(w, "Error accepting RFQ", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "Accepted"})
}

func (h *PartnerHandler) HandleRejectRFQ(w http.ResponseWriter, r *http.Request) {
	type RejectRequest struct {
		TripID string `json:"trip_id"`
	}
	var req RejectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	tripID, err := uuid.Parse(req.TripID)
	if err != nil {
		http.Error(w, "Invalid Trip ID", http.StatusBadRequest)
		return
	}

	partnerID, err := h.getPartnerID(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if err := h.Repo.RejectRFQ(r.Context(), tripID, partnerID); err != nil {
		http.Error(w, "Error rejecting RFQ", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "Rejected"})
}

func (h *PartnerHandler) HandleListUninvoicedTrips(w http.ResponseWriter, r *http.Request) {
	supplierID, err := h.getPartnerID(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	trips, err := h.Repo.ListUninvoicedTrips(r.Context(), supplierID)
	if err != nil {
		http.Error(w, "Error listing trips", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(trips)
}

func (h *PartnerHandler) HandleListInvoices(w http.ResponseWriter, r *http.Request) {
	supplierID, err := h.getPartnerID(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	invoices, err := h.Repo.ListInvoices(r.Context(), supplierID)
	if err != nil {
		http.Error(w, "Error listing invoices", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(invoices)
}

func (h *PartnerHandler) HandleSubmitInvoice(w http.ResponseWriter, r *http.Request) {
	var req models.SubmitInvoiceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		fmt.Printf("ERROR: Failed to decode request body: %v\n", err)
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	tripID, err := uuid.Parse(req.TripID)
	if err != nil {
		fmt.Printf("ERROR: Invalid trip ID '%s': %v\n", req.TripID, err)
		http.Error(w, "Invalid trip ID", http.StatusBadRequest)
		return
	}

	supplierID, err := h.getPartnerID(r)
	if err != nil {
		fmt.Printf("ERROR: Failed to get partner ID: %v\n", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	invoice := models.Invoice{
		TripID:            tripID,
		SupplierCompanyID: supplierID,
		InvoiceNumber:     req.InvoiceNumber,
		Amount:            req.Amount,
		Status:            models.InvoiceStatusPending,
	}

	fmt.Printf("DEBUG: Submitting invoice: TripID=%s, SupplierID=%s, InvoiceNumber=%s, Amount=%.2f\n",
		tripID, supplierID, req.InvoiceNumber, req.Amount)

	if err := h.Repo.SubmitInvoice(r.Context(), invoice); err != nil {
		fmt.Printf("ERROR: Failed to submit invoice to database: %v\n", err)
		// Check if it's a duplicate key error
		if strings.Contains(err.Error(), "duplicate key") || strings.Contains(err.Error(), "unique constraint") {
			http.Error(w, fmt.Sprintf("Invoice number '%s' already exists. Please use a different invoice number.", req.InvoiceNumber), http.StatusConflict)
			return
		}
		http.Error(w, "Error submitting invoice", http.StatusInternalServerError)
		return
	}

	fmt.Printf("SUCCESS: Invoice submitted successfully\n")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"status": "Invoice Submitted"})
}

func (h *PartnerHandler) HandleUpdateInvoice(w http.ResponseWriter, r *http.Request) {
	invoiceIDStr := chi.URLParam(r, "id")
	invoiceID, err := uuid.Parse(invoiceIDStr)
	if err != nil {
		http.Error(w, "Invalid invoice ID", http.StatusBadRequest)
		return
	}

	var req models.SubmitInvoiceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if err := h.Repo.UpdateInvoice(r.Context(), invoiceID, req.InvoiceNumber, req.Amount); err != nil {
		fmt.Printf("ERROR: Failed to update invoice: %v\n", err)
		http.Error(w, "Error updating invoice", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "Invoice Updated"})
}

func (h *PartnerHandler) HandleCloseInvoice(w http.ResponseWriter, r *http.Request) {
	invoiceIDStr := chi.URLParam(r, "id")
	invoiceID, err := uuid.Parse(invoiceIDStr)
	if err != nil {
		http.Error(w, "Invalid invoice ID", http.StatusBadRequest)
		return
	}

	if err := h.Repo.CloseInvoice(r.Context(), invoiceID); err != nil {
		fmt.Printf("ERROR: Failed to close invoice: %v\n", err)
		http.Error(w, "Error closing invoice", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "Invoice Closed"})
}
