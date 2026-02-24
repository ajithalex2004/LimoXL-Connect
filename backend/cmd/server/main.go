package main

import (
	"fmt"
	"limoxlink-backend/db"
	"limoxlink-backend/internal/api"
	"limoxlink-backend/internal/middleware"
	"limoxlink-backend/internal/repository"
	"log"
	"net/http"

	"io"
	"os"
	"strings"

	"github.com/go-chi/chi/v5"
	middlewareChi "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func loadEnv() {
	file, err := os.Open(".env")
	if err != nil {
		return
	}
	defer file.Close()

	var data []byte
	data, err = io.ReadAll(file)
	if err != nil {
		return
	}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			if key != "" {
				os.Setenv(key, value)
				log.Printf("Loaded env: %s", key)
			}
		}
	}
}

func main() {
	loadEnv()
	// Initialize Database
	if err := db.InitDB(); err != nil {
		log.Printf("Warning: Failed to initialize database: %v. Ensure Docker is running.\n", err)
	}

	// Auto-migrate for demo convenience
	if err := db.RunMigrations(); err != nil {
		log.Printf("Warning: Failed to run migrations: %v\n", err)
	}

	// Initialize Repositories
	companyRepo := repository.NewPostgresCompanyRepo(db.DB)
	userRepo := repository.NewPostgresUserRepo(db.DB)
	vehicleRepo := repository.NewPostgresVehicleRepo(db.DB)
	driverRepo := repository.NewPostgresDriverRepo(db.DB) // New
	tripRepo := repository.NewPostgresTripRepo(db.DB)
	outsourceCompanyRepo := repository.NewOutsourceCompanyRepo(db.DB) // New

	// Initialize Handlers
	companyHandler := api.NewCompanyHandler(companyRepo)
	userHandler := api.NewUserHandler(userRepo)
	vehicleHandler := api.NewVehicleHandler(vehicleRepo)
	tripHandler := api.NewTripHandler(tripRepo)
	authHandler := api.NewAuthHandler(userRepo)
	operatorHandler := api.NewOperatorHandler(companyRepo, userRepo, tripRepo, outsourceCompanyRepo)
	fleetHandler := api.NewFleetHandler(vehicleRepo, driverRepo) // New

	r := chi.NewRouter()

	// Middleware
	r.Use(middlewareChi.Logger)
	r.Use(middlewareChi.Recoverer)

	// CORS Config
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173", "http://localhost:5174", "http://localhost:5175", "http://localhost:5176", "http://localhost:5177", "http://localhost:5178", "https://limoxl-connect-production.up.railway.app"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Public Routes
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to LimoXLink API"))
	})
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	r.Route("/api", func(r chi.Router) {
		// Auth
		r.Post("/auth/login", authHandler.HandleLogin)
		r.Post("/auth/change-password", authHandler.HandleChangePassword)

		// Operator Routes
		r.Route("/operator", func(r chi.Router) {
			r.Get("/outsource-companies", operatorHandler.HandleListOutsourceCompanies)
			r.Post("/outsource-companies", operatorHandler.HandleCreateOutsourceCompany)
			r.Put("/outsource-companies/{id}", operatorHandler.HandleUpdateOutsourceCompany)
			r.Delete("/outsource-companies/{id}", operatorHandler.HandleDeleteOutsourceCompany)
			r.Get("/quotes", operatorHandler.HandleListQuotes)
			r.Post("/quotes/{id}/accept", operatorHandler.HandleAcceptQuote)
			r.Post("/quotes/{id}/reject", operatorHandler.HandleRejectQuote)
			r.Get("/trips", operatorHandler.HandleListAllTrips)
			r.Post("/trips", operatorHandler.HandleCreateTrip)
			r.Post("/trips/{id}/assign", operatorHandler.HandleAssignOutsource)
			r.Post("/trips/{id}/dispatch", operatorHandler.HandleDispatchTrip)
		}) // New

		// Operator Fleet (Reusing FleetHandler)
		r.Get("/operator/vehicles", fleetHandler.ListVehicles)
		r.Post("/operator/vehicles", fleetHandler.CreateVehicle) // New
		r.Get("/operator/drivers", fleetHandler.ListDrivers)
		r.Post("/operator/drivers", fleetHandler.CreateDriver) // New

		// Companies
		r.Post("/companies", companyHandler.CreateCompany)
		r.Get("/companies", companyHandler.ListCompanies)
		r.Get("/companies/{id}", companyHandler.GetCompany)

		// Users
		r.Post("/users", userHandler.CreateUser)
		r.Get("/users/{id}", userHandler.GetUser)

		// Vehicles
		r.Post("/vehicles", vehicleHandler.CreateVehicle)
		r.Get("/vehicles", vehicleHandler.ListVehicles)
		r.Get("/vehicles/{id}", vehicleHandler.GetVehicle)

		// Marketplace Routes
		r.Post("/marketplace/search-vehicles", tripHandler.HandleSearchVehicles)
		r.Post("/marketplace/book-vehicle", tripHandler.HandleBookVehicle)
		r.Get("/trips", tripHandler.HandleListTrips)
		r.Get("/trips/{id}", tripHandler.HandleGetTrip)

		// Partner Fleet
		r.Get("/partner/vehicles", fleetHandler.ListVehicles)
		r.Post("/partner/vehicles", fleetHandler.CreateVehicle)
		r.Get("/partner/drivers", fleetHandler.ListDrivers)
		r.Post("/partner/drivers", fleetHandler.CreateDriver)

		// Partner Routes (Protected)
		r.Group(func(r chi.Router) {
			r.Use(middleware.AuthMiddleware)
			partnerHandler := api.NewPartnerHandler(tripRepo)
			r.Get("/partner/rfqs", partnerHandler.ListRFQs)
			r.Get("/partner/rfqs/history", partnerHandler.HandleListRFQHistory)
			r.Get("/partner/trips", partnerHandler.ListAssignedTrips)
			r.Post("/partner/quotes", partnerHandler.HandleSubmitQuote)
			r.Post("/partner/assign", partnerHandler.HandleAssignDriver)
			r.Post("/partner/accept", partnerHandler.HandleAcceptRFQ)
			r.Post("/partner/reject", partnerHandler.HandleRejectRFQ)
			r.Get("/partner/trips/completed", partnerHandler.HandleListUninvoicedTrips)
			r.Post("/partner/invoices", partnerHandler.HandleSubmitInvoice)
			r.Put("/partner/invoices/{id}", partnerHandler.HandleUpdateInvoice)
			r.Post("/partner/invoices/{id}/close", partnerHandler.HandleCloseInvoice)
			r.Get("/partner/invoices", partnerHandler.HandleListInvoices)
		})

		// Secure Driver Link Routes
		secureHandler := api.NewSecureLinkHandler(tripRepo)
		r.Get("/status/{token}", secureHandler.GetTripStatus)
		r.Post("/status/{token}", secureHandler.UpdateStatus)

		// Protected Routes
		r.Group(func(r chi.Router) {
			// r.Use(middleware.AuthMiddleware) // Commented out for MVP testing
			r.Get("/protected/check", func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("Authenticated"))
			})
		})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	// Support both ":8080" and "8080" formats
	if !strings.HasPrefix(port, ":") {
		port = ":" + port
	}

	fmt.Printf("Server starting on port %s\n", port)
	if err := http.ListenAndServe(port, r); err != nil {
		log.Fatal(err)
	}
}
