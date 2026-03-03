package repository

import (
	"context"
	"database/sql"
	"limoxlink-backend/internal/models"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
)

type TripRepository interface {
	SearchAvailableVehicles(ctx context.Context, req models.VehicleSearchRequest) ([]models.VehicleSearchResult, error)
	BookVehicle(ctx context.Context, req models.BookVehicleRequest) (*models.BookVehicleResponse, error)
	ListTrips(ctx context.Context) ([]models.Trip, error)
	ListOpenRFQs(ctx context.Context, partnerID uuid.UUID) ([]models.Trip, error)
	ListRFQHistory(ctx context.Context, partnerID uuid.UUID) ([]models.Trip, error)
	ListPartnerTrips(ctx context.Context, partnerID uuid.UUID) ([]models.Trip, error)
	SubmitQuote(ctx context.Context, req models.SubmitQuoteRequest, supplierID uuid.UUID) error
	ListTripOffers(ctx context.Context, operatorCompanyID uuid.UUID) ([]models.TripOffer, error)
	GetTripByToken(ctx context.Context, token string) (*models.Trip, error)
	UpdateTripStatus(ctx context.Context, tripID uuid.UUID, status models.TripStatus) error
	AssignDriver(ctx context.Context, tripID uuid.UUID, driverName, driverPhone, vehicleModel, vehiclePlate string) error
	AcceptRFQ(ctx context.Context, tripID uuid.UUID, partnerID uuid.UUID) error
	RejectRFQ(ctx context.Context, tripID uuid.UUID, partnerID uuid.UUID) error
	GetTrip(ctx context.Context, id uuid.UUID) (*models.Trip, error)

	AcceptTripOffer(ctx context.Context, offerID uuid.UUID) error
	RejectTripOffer(ctx context.Context, offerID uuid.UUID) error
	SubmitInvoice(ctx context.Context, invoice models.Invoice) error
	ListUninvoicedTrips(ctx context.Context, partnerID uuid.UUID) ([]models.Trip, error)
	UpdateInvoice(ctx context.Context, invoiceID uuid.UUID, invoiceNumber string, amount float64) error
	CloseInvoice(ctx context.Context, invoiceID uuid.UUID) error
	ListInvoices(ctx context.Context, partnerID uuid.UUID) ([]models.Invoice, error)
	ListOperatorTrips(ctx context.Context, operatorID uuid.UUID) ([]models.Trip, error)
	AssignOutsource(ctx context.Context, tripID uuid.UUID, partnerIDs []uuid.UUID) error
	CreateTrip(ctx context.Context, trip *models.Trip) error
	AssignInternalDispatch(ctx context.Context, tripID uuid.UUID, driverID uuid.UUID, vehicleID uuid.UUID) error
}

type PostgresTripRepo struct {
	DB *sql.DB
}

func NewPostgresTripRepo(db *sql.DB) *PostgresTripRepo {
	return &PostgresTripRepo{DB: db}
}

// SearchAvailableVehicles returns a list of mock available vehicles based on proximity
func (r *PostgresTripRepo) SearchAvailableVehicles(ctx context.Context, req models.VehicleSearchRequest) ([]models.VehicleSearchResult, error) {
	// Mock Database of Vehicles in the Area
	allVehicles := []models.VehicleSearchResult{
		{
			VehicleID:    uuid.New().String(),
			Make:         "GMC",
			Model:        "Yukon",
			Year:         2024,
			Color:        "Black",
			VehicleType:  "SUV",
			VehicleClass: "Premium",
			SupplierName: "Elite Fleet",
			Location:     "Dubai Marina",
			Distance:     "2 km away",
		},
		{
			VehicleID:    uuid.New().String(),
			Make:         "Mercedes",
			Model:        "S-Class",
			Year:         2023,
			Color:        "Black",
			VehicleType:  "Limo",
			VehicleClass: "Luxury",
			SupplierName: "Royal Rides",
			Location:     "Downtown Dubai",
			Distance:     "5 km away",
		},
		{
			VehicleID:    uuid.New().String(),
			Make:         "Toyota",
			Model:        "Camry",
			Year:         2023,
			Color:        "White",
			VehicleType:  "Sedan",
			VehicleClass: "Standard",
			SupplierName: "City Transport",
			Location:     "Business Bay",
			Distance:     "3 km away",
		},
		{
			VehicleID:    uuid.New().String(),
			Make:         "Mercedes",
			Model:        "V-Class",
			Year:         2024,
			Color:        "Black",
			VehicleType:  "Van",
			VehicleClass: "Premium",
			SupplierName: "Elite Fleet",
			Location:     "Palm Jumeirah",
			Distance:     "8 km away",
		},
		{
			VehicleID:    uuid.New().String(),
			Make:         "Lexus",
			Model:        "ES350",
			Year:         2023,
			Color:        "Silver",
			VehicleType:  "Sedan",
			VehicleClass: "Business",
			SupplierName: "Prestige Limo",
			Location:     "JLT",
			Distance:     "4 km away",
		},
		{
			VehicleID:    uuid.New().String(),
			Make:         "Cadillac",
			Model:        "Escalade",
			Year:         2024,
			Color:        "Black",
			VehicleType:  "SUV",
			VehicleClass: "Luxury",
			SupplierName: "VIP Chauffeur",
			Location:     "Dubai Mall",
			Distance:     "6 km away",
		},
		{
			VehicleID:    uuid.New().String(),
			Make:         "BMW",
			Model:        "7 Series",
			Year:         2024,
			Color:        "Dark Blue",
			VehicleType:  "Limo",
			VehicleClass: "Luxury",
			SupplierName: "Executive Rides",
			Location:     "DIFC",
			Distance:     "4.5 km away",
		},
		{
			VehicleID:    uuid.New().String(),
			Make:         "Tesla",
			Model:        "Model S",
			Year:         2023,
			Color:        "White",
			VehicleType:  "Sedan",
			VehicleClass: "Eco",
			SupplierName: "Green Mobility",
			Location:     "Media City",
			Distance:     "7 km away",
		},
	}

	// Filter Logic
	var filtered []models.VehicleSearchResult
	for _, v := range allVehicles {
		if req.VehicleClass != "" && req.VehicleClass != "Any" && v.VehicleClass != req.VehicleClass {
			continue
		}
		if req.VehicleType != "" && req.VehicleType != "Any" && v.VehicleType != req.VehicleType {
			continue
		}
		// In real app, check proximity to req.PickupZone
		filtered = append(filtered, v)
	}

	return filtered, nil
}

// Helper for case-insensitive contains
func contains(s, substr string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}

func (r *PostgresTripRepo) BookVehicle(ctx context.Context, req models.BookVehicleRequest) (*models.BookVehicleResponse, error) {
	// In a real implementation, this would:
	// 1. Create a Trip record in the DB
	// 2. Link it to the external supplier/vehicle
	// 3. Send a notification to the supplier

	// Mock Success Response
	return &models.BookVehicleResponse{
		TripID:  uuid.New().String(),
		Status:  "Farmed-Out",
		Message: "Booking request sent to supplier successfully.",
	}, nil
}

func (r *PostgresTripRepo) ListTrips(ctx context.Context) ([]models.Trip, error) {
	// For MVP, assuming the logged in operator is the "Limo Operator Co"
	operatorIDStr := "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a00"

	query := `
		SELECT 
			id, reference_no, pickup_zone, dropoff_zone, pickup_time::timestamp,
			status, passenger_name, passenger_phone, vehicle_type_requested,
			created_at::timestamp, rfq_number, pickup_landmark, dropoff_landmark, service_type
		FROM trips
		WHERE requesting_company_id = $1
		ORDER BY created_at DESC
	`
	rows, err := r.DB.QueryContext(ctx, query, operatorIDStr)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var trips []models.Trip
	for rows.Next() {
		var t models.Trip
		var pZone, dZone, pName, pPhone, vType, rfqNo, pLandmark, dLandmark, sType sql.NullString
		err := rows.Scan(
			&t.ID, &t.BookingReference, &pZone, &dZone, &t.PickupTime,
			&t.Status, &pName, &pPhone, &vType, &t.CreatedAt, &rfqNo,
			&pLandmark, &dLandmark, &sType,
		)
		if err != nil {
			return nil, err
		}
		t.PickupZone = pZone.String
		t.DropoffZone = dZone.String
		t.PassengerName = pName.String
		t.PassengerPhone = pPhone.String
		t.RequestedVehicleType = vType.String
		t.RFQNumber = rfqNo.String
		t.PickupLandmark = pLandmark.String
		t.DropoffLandmark = dLandmark.String
		t.ServiceType = sType.String

		// Fill mock data for now for missing fields in query but needed for UI
		t.SupplierName = "N/A"

		trips = append(trips, t)
	}
	return trips, nil
}

func (r *PostgresTripRepo) ListOpenRFQs(ctx context.Context, partnerID uuid.UUID) ([]models.Trip, error) {
	query := `
		SELECT
			t.id, t.reference_no, t.pickup_zone, t.dropoff_zone, t.pickup_time,
			t.status, t.passenger_name, t.vehicle_type_requested, t.visibility,
			t.created_at, t.rfq_number, t.pickup_landmark, t.dropoff_landmark, t.service_type
		FROM trips t
		LEFT JOIN trip_access ta ON t.id = ta.trip_id
		WHERE t.status IN ('MARKETPLACE_SEARCH', 'OFFERED')
		  AND (t.visibility = 'PUBLIC' OR ta.company_id = $1)
		  AND NOT EXISTS (
			SELECT 1 FROM trip_offers off 
			WHERE off.trip_id = t.id 
			  AND off.supplier_company_id = $1 
			  AND off.status = 'REJECTED'
		  )
	`
	rows, err := r.DB.QueryContext(ctx, query, partnerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var trips []models.Trip
	for rows.Next() {
		var t models.Trip
		var pZone, dZone, vType, rfqNo, pName, visibility, pLandmark, dLandmark, sType sql.NullString // Handle potential nulls
		err := rows.Scan(
			&t.ID, &t.BookingReference, &pZone, &dZone, &t.PickupTime,
			&t.Status, &pName, &vType, &visibility, &t.CreatedAt, &rfqNo,
			&pLandmark, &dLandmark, &sType,
		)
		if err != nil {
			return nil, err
		}
		t.PickupZone = pZone.String
		t.DropoffZone = dZone.String
		t.RequestedVehicleType = vType.String
		t.RFQNumber = rfqNo.String
		t.PassengerName = pName.String
		t.Visibility = visibility.String
		t.PickupLandmark = pLandmark.String
		t.DropoffLandmark = dLandmark.String
		t.ServiceType = sType.String
		// Defaults
		t.RequestedVehicleClass = "Standard"
		t.RequestedVehicleGroup = "Sedan"

		trips = append(trips, t)
	}
	log.Printf("ListOpenRFQs: Found %d trips for partner %s", len(trips), partnerID)
	return trips, nil
}

func (r *PostgresTripRepo) ListRFQHistory(ctx context.Context, partnerID uuid.UUID) ([]models.Trip, error) {
	// History = REJECTED, TIMEOUT, FAILED, COMPLETED (if we want detailed history here, but usually Closed)
	// User said "History Page of RFQ Submitted & Rejected" -> Actually for Partner (Outsource Company),
	// this likely refers to quotes they submitted that were rejected, or RFQs they saw that expired.
	// But strictly speaking, the Partner *sees* Trips.
	// If a Partner submitted a quote and it was REJECTED, the TRIP status might still be MARKETPLACE_SEARCH (if rejected for THIS partner but open for others) - wait multiple offers.
	// OR if the Trip was assigned to someone else, it's irrelevant to this partner unless they were the one rejected?
	// Actually, `trip_offers` status is what matters for the Partner's history of "Submitted & Rejected".
	// The user said "Once it is rejected, it should moved to History Page".
	// So we need to fetch trips where THIS partner has a REJECTED offer, OR trips that are CANCELLED/TIMEOUT from the marketplace.

	// Let's simplify: Fetch Trips where there is an offer from this company with status REJECTED,
	// OR trips where the trip status itself implies history (completed/cancelled) AND this partner had some interaction?
	// User phrasing: "Once [the quote] is rejected, it should moved to History Page".
	// So we are looking for REJECTED Offers.
	// We should return `Trip` objects but maybe enriched with Offer status?
	// Or just Trips.

	query := `
		SELECT
			t.id, t.reference_no, t.pickup_zone, t.dropoff_zone, t.pickup_time,
			t.status, t.passenger_name, t.vehicle_type_requested, t.visibility,
			t.created_at, t.rfq_number, t.service_type
		FROM trips t
		JOIN trip_offers toff ON t.id = toff.trip_id
		WHERE toff.supplier_company_id = $1
		  AND toff.status = 'REJECTED'
		ORDER BY toff.created_at DESC
	`
	rows, err := r.DB.QueryContext(ctx, query, partnerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var trips []models.Trip
	for rows.Next() {
		var t models.Trip
		var pZone, dZone, vType, rfqNo, pName, visibility, sType sql.NullString
		err := rows.Scan(
			&t.ID, &t.BookingReference, &pZone, &dZone, &t.PickupTime,
			&t.Status, &pName, &vType, &visibility, &t.CreatedAt, &rfqNo, &sType,
		)
		if err != nil {
			return nil, err
		}
		t.PickupZone = pZone.String
		t.DropoffZone = dZone.String
		t.RequestedVehicleType = vType.String
		t.RFQNumber = rfqNo.String
		t.PassengerName = pName.String
		t.Visibility = visibility.String
		t.ServiceType = sType.String
		t.RequestedVehicleClass = "Standard"
		t.RequestedVehicleGroup = "Sedan"
		trips = append(trips, t)
	}
	return trips, nil
}

func (r *PostgresTripRepo) ListPartnerTrips(ctx context.Context, partnerID uuid.UUID) ([]models.Trip, error) {
	query := `
		SELECT
			id, reference_no, pickup_zone, dropoff_zone, pickup_time::timestamp,
			status, passenger_name, vehicle_type_requested, driver_link_token,
			created_at::timestamp, rfq_number, pickup_landmark, dropoff_landmark, service_type
		FROM trips
		WHERE fulfillment_company_id = $1
		ORDER BY created_at DESC
	`
	// Note: We normally verify partnerID vs user claims.
	// For this task, we assume we want all trips assigned to "Demo Partner Co"
	// Actually, the seed data doesn't assign trips yet. AcceptRFQ will do that.

	rows, err := r.DB.QueryContext(ctx, query, partnerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var trips []models.Trip
	for rows.Next() {
		var t models.Trip
		var pZone, dZone, vType, dToken, rfqNo, pName, pLandmark, dLandmark, sType sql.NullString
		err := rows.Scan(
			&t.ID, &t.BookingReference, &pZone, &dZone, &t.PickupTime,
			&t.Status, &pName, &vType, &dToken, &t.CreatedAt, &rfqNo,
			&pLandmark, &dLandmark, &sType,
		)
		if err != nil {
			return nil, err
		}
		t.PickupZone = pZone.String
		t.DropoffZone = dZone.String
		t.RequestedVehicleType = vType.String
		t.DriverLinkToken = dToken.String
		t.RFQNumber = rfqNo.String
		t.PassengerName = pName.String
		t.PickupLandmark = pLandmark.String
		t.DropoffLandmark = dLandmark.String
		t.ServiceType = sType.String
		trips = append(trips, t)
	}
	return trips, nil
}

func (r *PostgresTripRepo) SubmitQuote(ctx context.Context, req models.SubmitQuoteRequest, supplierID uuid.UUID) error {
	query := `
		INSERT INTO trip_offers (trip_id, supplier_company_id, price, notes, status, created_at)
		VALUES ($1, $2, $3, $4, 'PENDING', NOW())
	`
	_, err := r.DB.ExecContext(ctx, query, req.TripID, supplierID, req.Price, req.Notes)
	if err != nil {
		return err
	}

	// Also update trip status to OFFERED if it was MARKETPLACE_SEARCH
	updateQuery := `UPDATE trips SET status = 'OFFERED' WHERE id = $1 AND status = 'MARKETPLACE_SEARCH'`
	r.DB.ExecContext(ctx, updateQuery, req.TripID)

	return nil
}

func (r *PostgresTripRepo) ListTripOffers(ctx context.Context, operatorCompanyID uuid.UUID) ([]models.TripOffer, error) {
	query := `
		SELECT 
			toff.id, toff.trip_id, toff.supplier_company_id, c.name, 
			toff.status, toff.score, toff.price, toff.notes, toff.created_at,
			t.reference_no, t.rfq_number
		FROM trip_offers toff
		JOIN companies c ON toff.supplier_company_id = c.id
		JOIN trips t ON toff.trip_id = t.id
		WHERE t.requesting_company_id = $1
		ORDER BY toff.created_at DESC
	`
	rows, err := r.DB.QueryContext(ctx, query, operatorCompanyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	offers := []models.TripOffer{}
	for rows.Next() {
		var o models.TripOffer
		var notes, rfqNo sql.NullString
		err := rows.Scan(
			&o.ID, &o.TripID, &o.SupplierCompanyID, &o.SupplierName,
			&o.Status, &o.Score, &o.Price, &notes, &o.CreatedAt,
			&o.BookingReference, &rfqNo,
		)
		if err != nil {
			return nil, err
		}
		o.Notes = notes.String
		o.RFQNumber = rfqNo.String
		offers = append(offers, o)
	}
	return offers, nil
}

func (r *PostgresTripRepo) GetTripByToken(ctx context.Context, token string) (*models.Trip, error) {
	query := `
		SELECT 
			id, reference_no, pickup_zone, dropoff_zone, 
			pickup_time, status, passenger_name, passenger_phone, driver_link_token, service_type
		FROM trips 
		WHERE driver_link_token = $1
	`
	row := r.DB.QueryRowContext(ctx, query, token)

	var trip models.Trip
	var pZone, dZone, dToken, pName, pPhone, sType sql.NullString // GetTripByToken doesn't fetch RFQNumber yet in query?
	// Actually GetTripByToken query lines 421-425:
	// SELECT id, reference_no, pickup_zone, dropoff_zone, pickup_time, status, passenger_name, passenger_phone, driver_link_token
	// It does NOT have rfq_number. So I should NOT modify it unless I update the query.
	// I'll skip GetTripByToken for now as I didn't verify its query update.
	err := row.Scan(
		&trip.ID, &trip.BookingReference, &pZone, &dZone,
		&trip.PickupTime, &trip.Status, &pName, &pPhone, &dToken, &sType,
	)
	if err != nil {
		return nil, err
	}
	trip.PickupZone = pZone.String
	trip.DropoffZone = dZone.String
	trip.PassengerName = pName.String
	trip.PassengerPhone = pPhone.String
	trip.DriverLinkToken = dToken.String
	trip.ServiceType = sType.String

	return &trip, nil
}

func (r *PostgresTripRepo) UpdateTripStatus(ctx context.Context, tripID uuid.UUID, status models.TripStatus) error {
	query := `UPDATE trips SET status = $1, updated_at = NOW() WHERE id = $2`
	_, err := r.DB.ExecContext(ctx, query, status, tripID)
	return err
}

func (r *PostgresTripRepo) AssignDriver(ctx context.Context, tripID uuid.UUID, driverName, driverPhone, vehicleModel, vehiclePlate string) error {
	// Update trip status and set driver details (mocked for now as we don't have driver/vehicle tables joined strictly yet)
	// But crucially, generate a token
	token := uuid.New().String()
	query := `
		UPDATE trips 
		SET status = 'DRIVER_ASSIGNED', 
		    driver_link_token = $2
		WHERE id = $1
	`
	_, err := r.DB.ExecContext(ctx, query, tripID, token)
	return err
}

func (r *PostgresTripRepo) AcceptRFQ(ctx context.Context, tripID uuid.UUID, partnerID uuid.UUID) error {
	query := `
		UPDATE trips
		SET status = 'ACCEPTED',
		    fulfillment_company_id = $2,
			updated_at = NOW()
		WHERE id = $1
	`
	_, err := r.DB.ExecContext(ctx, query, tripID, partnerID)
	return err
}

func (r *PostgresTripRepo) RejectRFQ(ctx context.Context, tripID uuid.UUID, partnerID uuid.UUID) error {
	// We create a "REJECTED" offer to mark this trip as hidden/rejected for this partner.
	// This does not change the Trip Status itself (so other partners can still see it).
	offerID := uuid.New()
	query := `
		INSERT INTO trip_offers (id, trip_id, supplier_company_id, status, price, score, created_at)
		VALUES ($1, $2, $3, 'REJECTED', 0, 0, NOW())
	`
	_, err := r.DB.ExecContext(ctx, query, offerID, tripID, partnerID)
	return err
}
func (r *PostgresTripRepo) GetTrip(ctx context.Context, id uuid.UUID) (*models.Trip, error) {
	query := `
		SELECT 
			t.id, t.requesting_company_id, t.reference_no, t.pickup_zone, t.dropoff_zone, 
			t.pickup_time::timestamp, t.status, t.vehicle_type_requested, t.passenger_name, t.rfq_number,
			t.pickup_landmark, t.dropoff_landmark, t.service_type
		FROM trips t
		WHERE t.id = $1
	`
	// Simplified select for MVP demo mostly needing passenger/route info
	row := r.DB.QueryRowContext(ctx, query, id)

	var trip models.Trip
	var rfqNo, pLandmark, dLandmark, sType sql.NullString
	err := row.Scan(
		&trip.ID, &trip.RequestingCompanyID, &trip.BookingReference, &trip.PickupZone, &trip.DropoffZone,
		&trip.PickupTime, &trip.Status, &trip.RequestedVehicleType, &trip.PassengerName, &rfqNo,
		&pLandmark, &dLandmark, &sType,
	)
	if err != nil {
		return nil, err
	}
	trip.RFQNumber = rfqNo.String
	trip.PickupLandmark = pLandmark.String
	trip.DropoffLandmark = dLandmark.String
	trip.ServiceType = sType.String
	return &trip, nil
}

func (r *PostgresTripRepo) AcceptTripOffer(ctx context.Context, offerID uuid.UUID) error {
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 1. Get the Offer details to find TripID and SupplierID
	var tripID, supplierID uuid.UUID
	queryGet := `SELECT trip_id, supplier_company_id FROM trip_offers WHERE id = $1`
	err = tx.QueryRowContext(ctx, queryGet, offerID).Scan(&tripID, &supplierID)
	if err != nil {
		return err
	}

	// 2. Update the accepted offer status
	queryUpdateOffer := `UPDATE trip_offers SET status = 'ACCEPTED' WHERE id = $1`
	_, err = tx.ExecContext(ctx, queryUpdateOffer, offerID)
	if err != nil {
		return err
	}

	// 3. Mark other offers for this trip as REJECTED
	queryRejectOthers := `UPDATE trip_offers SET status = 'REJECTED' WHERE trip_id = $1 AND id != $2 AND status = 'PENDING'`
	_, err = tx.ExecContext(ctx, queryRejectOthers, tripID, offerID)
	if err != nil {
		return err
	}

	// 4. Update the Trip status and assign fulfillment company
	queryUpdateTrip := `
		UPDATE trips 
		SET status = 'ACCEPTED', 
		    fulfillment_company_id = $2,
			updated_at = NOW()
		WHERE id = $1
	`
	_, err = tx.ExecContext(ctx, queryUpdateTrip, tripID, supplierID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *PostgresTripRepo) RejectTripOffer(ctx context.Context, offerID uuid.UUID) error {
	query := `UPDATE trip_offers SET status = 'REJECTED' WHERE id = $1`
	_, err := r.DB.ExecContext(ctx, query, offerID)
	return err
}

func (r *PostgresTripRepo) SubmitInvoice(ctx context.Context, invoice models.Invoice) error {
	// VAT is 5% in UAE usually
	vatRate := 0.05
	// Platform fee is 10% for demo
	platformFeeRate := 0.10

	invoice.VATAmount = invoice.Amount * vatRate
	invoice.PlatformFee = invoice.Amount * platformFeeRate
	invoice.NetPayout = invoice.Amount - invoice.PlatformFee + invoice.VATAmount

	query := `
		INSERT INTO trip_invoices (
			trip_id, supplier_company_id, invoice_number, amount, 
			platform_fee, vat_amount, net_payout, status, created_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW())
	`
	_, err := r.DB.ExecContext(ctx, query,
		invoice.TripID, invoice.SupplierCompanyID, invoice.InvoiceNumber, invoice.Amount,
		invoice.PlatformFee, invoice.VATAmount, invoice.NetPayout, invoice.Status,
	)
	return err
}

func (r *PostgresTripRepo) UpdateInvoice(ctx context.Context, invoiceID uuid.UUID, invoiceNumber string, amount float64) error {
	query := `
		UPDATE trip_invoices
		SET invoice_number = $2, amount = $3, updated_at = NOW()
		WHERE id = $1
	`
	_, err := r.DB.ExecContext(ctx, query, invoiceID, invoiceNumber, amount)
	return err
}

func (r *PostgresTripRepo) CloseInvoice(ctx context.Context, invoiceID uuid.UUID) error {
	query := `
		UPDATE trip_invoices
		SET status = 'CLOSED', updated_at = NOW()
		WHERE id = $1
	`
	_, err := r.DB.ExecContext(ctx, query, invoiceID)
	return err
}

func (r *PostgresTripRepo) ListOperatorTrips(ctx context.Context, operatorID uuid.UUID) ([]models.Trip, error) {
	query := `
		SELECT 
			id, reference_no, pickup_zone, dropoff_zone, pickup_time::timestamp,
			status, passenger_name, passenger_phone, vehicle_type_requested,
			created_at::timestamp, rfq_number, pickup_landmark, dropoff_landmark, service_type
		FROM trips
		WHERE requesting_company_id = $1
		ORDER BY created_at DESC
	`
	rows, err := r.DB.QueryContext(ctx, query, operatorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var trips []models.Trip
	for rows.Next() {
		var t models.Trip
		var pZone, dZone, pName, pPhone, vType, rfqNo, pLandmark, dLandmark, sType sql.NullString
		err := rows.Scan(
			&t.ID, &t.BookingReference, &pZone, &dZone, &t.PickupTime,
			&t.Status, &pName, &pPhone, &vType, &t.CreatedAt, &rfqNo,
			&pLandmark, &dLandmark, &sType,
		)
		if err != nil {
			return nil, err
		}
		t.PickupZone = pZone.String
		t.DropoffZone = dZone.String
		t.PassengerName = pName.String
		t.PassengerPhone = pPhone.String
		t.RequestedVehicleType = vType.String
		t.RFQNumber = rfqNo.String
		t.PickupLandmark = pLandmark.String
		t.DropoffLandmark = dLandmark.String
		t.ServiceType = sType.String
		trips = append(trips, t)
	}
	return trips, nil
}

func (r *PostgresTripRepo) AssignOutsource(ctx context.Context, tripID uuid.UUID, partnerIDs []uuid.UUID) error {
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	visibility := "PUBLIC"
	if len(partnerIDs) > 0 {
		visibility = "DIRECT"
	}

	query := `
		UPDATE trips 
		SET status = 'MARKETPLACE_SEARCH', 
		    visibility = $2,
		    rfq_number = COALESCE(rfq_number, 'RFQ' || nextval('rfq_id_seq')),
		    updated_at = NOW()
		WHERE id = $1
	`
	_, err = tx.ExecContext(ctx, query, tripID, visibility)
	if err != nil {
		return err
	}

	if visibility == "DIRECT" {
		stmt, err := tx.PrepareContext(ctx, "INSERT INTO trip_access (trip_id, company_id) VALUES ($1, $2)")
		if err != nil {
			return err
		}
		defer stmt.Close()

		for _, pid := range partnerIDs {
			if _, err := stmt.ExecContext(ctx, tripID, pid); err != nil {
				return err
			}
		}
	}

	return tx.Commit()
}

func (r *PostgresTripRepo) ListUninvoicedTrips(ctx context.Context, partnerID uuid.UUID) ([]models.Trip, error) {
	// Note: The invoices table is for maintenance requests, not trips
	// For now, we return all COMPLETED trips for the partner
	// In the future, we should create a trip_invoices table
	query := `
		SELECT t.id, t.reference_no, t.pickup_zone, t.dropoff_zone, t.pickup_time, t.status, t.price
		FROM trips t
		WHERE t.fulfillment_company_id = $1 
		AND t.status = 'COMPLETED'
		ORDER BY t.pickup_time DESC
	`
	rows, err := r.DB.QueryContext(ctx, query, partnerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var trips []models.Trip
	for rows.Next() {
		var t models.Trip
		var pZone, dZone sql.NullString
		var price sql.NullFloat64
		err := rows.Scan(&t.ID, &t.BookingReference, &pZone, &dZone, &t.PickupTime, &t.Status, &price)
		if err != nil {
			return nil, err
		}
		t.PickupZone = pZone.String
		t.DropoffZone = dZone.String
		if price.Valid {
			t.Price = price.Float64
		}
		trips = append(trips, t)
	}
	return trips, nil
}

func (r *PostgresTripRepo) ListInvoices(ctx context.Context, partnerID uuid.UUID) ([]models.Invoice, error) {
	query := `
		SELECT i.id, i.trip_id, i.invoice_number, i.amount, i.status, i.created_at, t.reference_no
		FROM trip_invoices i
		JOIN trips t ON i.trip_id = t.id
		WHERE i.supplier_company_id = $1
		ORDER BY i.created_at DESC
	`
	rows, err := r.DB.QueryContext(ctx, query, partnerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var invoices []models.Invoice
	for rows.Next() {
		var inv models.Invoice
		var ref string
		err := rows.Scan(
			&inv.ID, &inv.TripID, &inv.InvoiceNumber, &inv.Amount,
			&inv.Status, &inv.CreatedAt, &ref,
		)
		if err != nil {
			return nil, err
		}
		inv.BookingReference = ref
		invoices = append(invoices, inv)
	}
	return invoices, nil
}

func (r *PostgresTripRepo) CreateTrip(ctx context.Context, trip *models.Trip) error {
	trip.ID = uuid.New()
	trip.CreatedAt = time.Now()
	trip.UpdatedAt = time.Now()

	// Default coordinates if not provided (Burj Khalifa)
	defaultLat := 25.1972
	defaultLng := 55.2744

	query := `
		INSERT INTO trips (
			id, requesting_company_id, reference_no, 
			pickup_zone, dropoff_zone, 
			passenger_name, passenger_phone, 
			pickup_time, status, vehicle_type_requested, 
			pickup_location, dropoff_location,
			price, pickup_landmark, dropoff_landmark, service_type,
			created_at, updated_at
		) VALUES (
			$1, $2, 'TR' || nextval('trip_id_seq'), 
			$3, $4, 
			$5, $6, 
			$7, $8, $9,
			ST_SetSRID(ST_MakePoint($10, $11), 4326),
			ST_SetSRID(ST_MakePoint($12, $13), 4326),
			$14, $15, $16, $17, $18, $19
		)
	`
	_, err := r.DB.ExecContext(ctx, query,
		trip.ID, trip.RequestingCompanyID, // $1, $2 (reference_no is generated)
		trip.PickupZone, trip.DropoffZone, // $3, $4
		trip.PassengerName, trip.PassengerPhone, // $5, $6
		trip.PickupTime, trip.Status, trip.RequestedVehicleType, // $7, $8, $9
		defaultLng, defaultLat, // $10, $11
		defaultLng, defaultLat, // $12, $13
		trip.Price, trip.PickupLandmark, trip.DropoffLandmark, trip.ServiceType, // $14, $15, $16, $17
		trip.CreatedAt, trip.UpdatedAt, // $18, $19
	)
	return err
}

func (r *PostgresTripRepo) AssignInternalDispatch(ctx context.Context, tripID uuid.UUID, driverID uuid.UUID, vehicleID uuid.UUID) error {
	query := `
		UPDATE trips
		SET 
			status = 'DRIVER_ASSIGNED',
			assigned_driver_id = $2,
			assigned_vehicle_id = $3,
			fulfillment_company_id = requesting_company_id,
			updated_at = NOW()
		WHERE id = $1
	`
	_, err := r.DB.ExecContext(ctx, query, tripID, driverID, vehicleID)
	return err
}
