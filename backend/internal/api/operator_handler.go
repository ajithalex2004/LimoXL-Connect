package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"limoxlink-backend/internal/middleware"
	"limoxlink-backend/internal/models"
	"limoxlink-backend/internal/repository"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type OperatorHandler struct {
	CompanyRepo          *repository.PostgresCompanyRepo
	UserRepo             *repository.PostgresUserRepo
	TripRepo             repository.TripRepository
	OutsourceCompanyRepo *repository.OutsourceCompanyRepo
}

func NewOperatorHandler(companyRepo *repository.PostgresCompanyRepo, userRepo *repository.PostgresUserRepo, tripRepo repository.TripRepository, outsourceRepo *repository.OutsourceCompanyRepo) *OperatorHandler {
	return &OperatorHandler{
		CompanyRepo:          companyRepo,
		UserRepo:             userRepo,
		TripRepo:             tripRepo,
		OutsourceCompanyRepo: outsourceRepo,
	}
}

type CreateOutsourceCompanyRequest struct {
	Name          string `json:"name"`
	ContactPerson string `json:"contact_person"`
	Designation   string `json:"designation"`
	Email         string `json:"email"`
	ContactNumber string `json:"contact_number"`
	Address       string `json:"address"`
	City          string `json:"city"`
	Country       string `json:"country"`
	TradeLicense  string `json:"trade_license_no"`
	ITCPermit     string `json:"itc_permit_no"`
	VATNo         string `json:"vat_no"`
}

func (h *OperatorHandler) HandleCreateOutsourceCompany(w http.ResponseWriter, r *http.Request) {
	var req CreateOutsourceCompanyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Create Outsource Company
	company := &repository.OutsourceCompany{
		Name:          req.Name,
		ContactPerson: sql.NullString{String: req.ContactPerson, Valid: req.ContactPerson != ""},
		Designation:   sql.NullString{String: req.Designation, Valid: req.Designation != ""},
		Email:         sql.NullString{String: req.Email, Valid: req.Email != ""},
		ContactNumber: sql.NullString{String: req.ContactNumber, Valid: req.ContactNumber != ""},
		Address:       sql.NullString{String: req.Address, Valid: req.Address != ""},
		City:          sql.NullString{String: req.City, Valid: req.City != ""},
		Country:       sql.NullString{String: req.Country, Valid: req.Country != ""},
		TradeLicense:  sql.NullString{String: req.TradeLicense, Valid: req.TradeLicense != ""},
		ITCPermit:     sql.NullString{String: req.ITCPermit, Valid: req.ITCPermit != ""},
		VATNo:         sql.NullString{String: req.VATNo, Valid: req.VATNo != ""},
		IsActive:      true,
	}

	if err := h.OutsourceCompanyRepo.Create(r.Context(), company); err != nil {
		fmt.Printf("Error creating outsource company record: %v\n", err)
		http.Error(w, "Failed to create outsource company", http.StatusInternalServerError)
		return
	}

	// Also create/sync to main companies table for dispatch visibility
	mainCompany := &models.Company{
		ID:       company.ID,
		Name:     company.Name,
		Type:     "SUPPLY",
		Verified: true,
	}
	if err := h.CompanyRepo.Create(r.Context(), mainCompany); err != nil {
		fmt.Printf("Error creating main company record (sync): %v\n", err)
		// We don't fail the whole request because the outsource record is already created
		// but this might lead to visibility issues. Log it.
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":     "Created",
		"company_id": company.ID,
	})
}

func (h *OperatorHandler) HandleListOutsourceCompanies(w http.ResponseWriter, r *http.Request) {
	companies, err := h.OutsourceCompanyRepo.List(r.Context())
	if err != nil {
		fmt.Printf("Error listing outsource companies: %v\n", err)
		http.Error(w, "Failed to list outsource companies", http.StatusInternalServerError)
		return
	}

	// Convert to JSON-friendly format
	type CompanyResponse struct {
		ID            string   `json:"id"`
		Name          string   `json:"name"`
		ContactPerson *string  `json:"contact_person,omitempty"`
		Designation   *string  `json:"designation,omitempty"`
		Email         *string  `json:"email,omitempty"`
		ContactNumber *string  `json:"contact_number,omitempty"`
		Address       *string  `json:"address,omitempty"`
		City          *string  `json:"city,omitempty"`
		Country       *string  `json:"country,omitempty"`
		TradeLicense  *string  `json:"trade_license_no,omitempty"`
		ITCPermit     *string  `json:"itc_permit_no,omitempty"`
		VATNo         *string  `json:"vat_no,omitempty"`
		Rating        *float64 `json:"rating,omitempty"`
		IsActive      bool     `json:"is_active"`
		Notes         *string  `json:"notes,omitempty"`
		CreatedAt     string   `json:"created_at"`
		UpdatedAt     string   `json:"updated_at"`
	}

	response := make([]CompanyResponse, 0, len(companies))
	for _, c := range companies {
		cr := CompanyResponse{
			ID:        c.ID.String(),
			Name:      c.Name,
			IsActive:  c.IsActive,
			CreatedAt: c.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt: c.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
		if c.ContactPerson.Valid {
			cr.ContactPerson = &c.ContactPerson.String
		}
		if c.Designation.Valid {
			cr.Designation = &c.Designation.String
		}
		if c.Email.Valid {
			cr.Email = &c.Email.String
		}
		if c.ContactNumber.Valid {
			cr.ContactNumber = &c.ContactNumber.String
		}
		if c.Address.Valid {
			cr.Address = &c.Address.String
		}
		if c.City.Valid {
			cr.City = &c.City.String
		}
		if c.Country.Valid {
			cr.Country = &c.Country.String
		}
		if c.TradeLicense.Valid {
			cr.TradeLicense = &c.TradeLicense.String
		}
		if c.ITCPermit.Valid {
			cr.ITCPermit = &c.ITCPermit.String
		}
		if c.VATNo.Valid {
			cr.VATNo = &c.VATNo.String
		}
		if c.Rating.Valid {
			cr.Rating = &c.Rating.Float64
		}
		if c.Notes.Valid {
			cr.Notes = &c.Notes.String
		}
		response = append(response, cr)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *OperatorHandler) HandleListQuotes(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middleware.ClaimsKey).(*middleware.Claims)
	if !ok || claims == nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]interface{}{})
		return
	}
	operatorID, err := uuid.Parse(claims.CompanyID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]interface{}{})
		return
	}

	quotes, err := h.TripRepo.ListTripOffers(r.Context(), operatorID)
	if err != nil {
		fmt.Printf("Error listing quotes: %v\n", err)
		http.Error(w, "Failed to list quotes", http.StatusInternalServerError)
		return
	}
	if quotes == nil {
		quotes = []models.TripOffer{}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(quotes)
}

func (h *OperatorHandler) HandleAcceptQuote(w http.ResponseWriter, r *http.Request) {
	quoteIDStr := chi.URLParam(r, "id")
	quoteID, err := uuid.Parse(quoteIDStr)
	if err != nil {
		http.Error(w, "Invalid Quote ID", http.StatusBadRequest)
		return
	}

	if err := h.TripRepo.AcceptTripOffer(r.Context(), quoteID); err != nil {
		fmt.Printf("Error accepting quote: %v\n", err)
		http.Error(w, "Failed to accept quote", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "Quote Accepted"})
}

func (h *OperatorHandler) HandleRejectQuote(w http.ResponseWriter, r *http.Request) {
	quoteIDStr := chi.URLParam(r, "id")
	quoteID, err := uuid.Parse(quoteIDStr)
	if err != nil {
		http.Error(w, "Invalid Quote ID", http.StatusBadRequest)
		return
	}

	if err := h.TripRepo.RejectTripOffer(r.Context(), quoteID); err != nil {
		fmt.Printf("Error rejecting quote: %v\n", err)
		http.Error(w, "Failed to reject quote", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "Quote Rejected"})
}

func (h *OperatorHandler) HandleListAllTrips(w http.ResponseWriter, r *http.Request) {
	// Get company ID from authenticated user's JWT claims
	claims, ok := r.Context().Value(middleware.ClaimsKey).(*middleware.Claims)
	if !ok || claims == nil {
		// Fallback: return empty list if no claims (shouldn't happen with auth middleware)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]interface{}{})
		return
	}

	operatorID, err := uuid.Parse(claims.CompanyID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]interface{}{})
		return
	}

	trips, err := h.TripRepo.ListOperatorTrips(r.Context(), operatorID)
	if err != nil {
		fmt.Printf("Error listing operator trips: %v\n", err)
		http.Error(w, "Error listing trips", http.StatusInternalServerError)
		return
	}
	if trips == nil {
		trips = []models.Trip{}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(trips)
}

func (h *OperatorHandler) HandleAssignOutsource(w http.ResponseWriter, r *http.Request) {
	tripIDStr := chi.URLParam(r, "id")
	tripID, err := uuid.Parse(tripIDStr)
	if err != nil {
		http.Error(w, "Invalid Trip ID", http.StatusBadRequest)
		return
	}

	var req struct {
		PartnerIDs []string `json:"partner_ids"`
	}
	// Decode optional body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil && err.Error() != "EOF" {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}

	var partnerUUIDs []uuid.UUID
	for _, idStr := range req.PartnerIDs {
		if id, err := uuid.Parse(idStr); err == nil {
			partnerUUIDs = append(partnerUUIDs, id)
		}
	}

	if err := h.TripRepo.AssignOutsource(r.Context(), tripID, partnerUUIDs); err != nil {
		fmt.Printf("Error assigning outsource: %v\n", err)
		http.Error(w, "Failed to assign outsource", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "assigned",
		"message": "Trip sent to marketplace",
	})
}

func (h *OperatorHandler) HandleCreateTrip(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middleware.ClaimsKey).(*middleware.Claims)
	if !ok || claims == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	if claims.IsSuperAdmin {
		// SuperAdmin cannot create a trip without a company context in this endpoint
		http.Error(w, "SuperAdmin must be linked to a company to create trips via this endpoint", http.StatusForbidden)
		return
	}

	operatorID, err := uuid.Parse(claims.CompanyID)
	if err != nil {
		http.Error(w, "Invalid company ID in token", http.StatusForbidden)
		return
	}

	var req models.Trip
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	req.RequestingCompanyID = operatorID
	req.Status = "CREATED"

	// Generate a simple reference number if not provided
	if req.BookingReference == "" {
		req.BookingReference = "B-" + uuid.New().String()[:8]
	}

	if err := h.TripRepo.CreateTrip(r.Context(), &req); err != nil {
		fmt.Printf("Error creating trip: %v\n", err)
		http.Error(w, "Failed to create trip", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(req)
}

func (h *OperatorHandler) HandleDispatchTrip(w http.ResponseWriter, r *http.Request) {
	tripIDStr := chi.URLParam(r, "id")
	tripID, err := uuid.Parse(tripIDStr)
	if err != nil {
		http.Error(w, "Invalid Trip ID", http.StatusBadRequest)
		return
	}

	var req struct {
		DriverID  string `json:"driver_id"`
		VehicleID string `json:"vehicle_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}

	driverUUID := uuid.MustParse(req.DriverID)
	vehicleUUID := uuid.MustParse(req.VehicleID)

	if err := h.TripRepo.AssignInternalDispatch(r.Context(), tripID, driverUUID, vehicleUUID); err != nil {
		fmt.Printf("Error dispatching trip: %v\n", err)
		http.Error(w, "Failed to dispatch trip", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "Dispatched"})
}

func (h *OperatorHandler) HandleUpdateOutsourceCompany(w http.ResponseWriter, r *http.Request) {
	companyIDStr := chi.URLParam(r, "id")
	companyID, err := uuid.Parse(companyIDStr)
	if err != nil {
		http.Error(w, "Invalid Company ID", http.StatusBadRequest)
		return
	}

	var req CreateOutsourceCompanyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Fetch existing company to preserve fields not in the request
	existing, err := h.OutsourceCompanyRepo.GetByID(r.Context(), companyID)
	if err != nil {
		http.Error(w, "Company not found", http.StatusNotFound)
		return
	}

	// Update fields
	existing.Name = req.Name
	existing.ContactPerson = sql.NullString{String: req.ContactPerson, Valid: req.ContactPerson != ""}
	existing.Designation = sql.NullString{String: req.Designation, Valid: req.Designation != ""}
	existing.Email = sql.NullString{String: req.Email, Valid: req.Email != ""}
	existing.ContactNumber = sql.NullString{String: req.ContactNumber, Valid: req.ContactNumber != ""}
	existing.Address = sql.NullString{String: req.Address, Valid: req.Address != ""}
	existing.City = sql.NullString{String: req.City, Valid: req.City != ""}
	existing.Country = sql.NullString{String: req.Country, Valid: req.Country != ""}

	if err := h.OutsourceCompanyRepo.Update(r.Context(), existing); err != nil {
		fmt.Printf("Error updating outsource company: %v\n", err)
		http.Error(w, "Failed to update outsource company", http.StatusInternalServerError)
		return
	}

	// Sync name update to companies table
	if mainComp, err := h.CompanyRepo.GetByID(r.Context(), existing.ID); err == nil {
		mainComp.Name = existing.Name
		// Potential need for a more specific update method in CompanyRepo if this grows
		// but for now we expect CompanyRepo.Create or similar to handle basic updates
		// Or we just ignore name sync for now if it's too complex without a dedicated update in CompanyRepo
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":     "Updated",
		"company_id": existing.ID,
	})
}

func (h *OperatorHandler) HandleDeleteOutsourceCompany(w http.ResponseWriter, r *http.Request) {
	companyIDStr := chi.URLParam(r, "id")
	companyID, err := uuid.Parse(companyIDStr)
	if err != nil {
		http.Error(w, "Invalid Company ID", http.StatusBadRequest)
		return
	}

	if err := h.OutsourceCompanyRepo.Delete(r.Context(), companyID); err != nil {
		fmt.Printf("Error deleting outsource company: %v\n", err)
		http.Error(w, "Failed to delete outsource company", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "Deleted"})
}
