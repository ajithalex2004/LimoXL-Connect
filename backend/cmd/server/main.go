package main

import (
	"encoding/json"
	"fmt"
	"limoxlink-backend/db"
	"limoxlink-backend/internal/api"
	"limoxlink-backend/internal/middleware"
	"limoxlink-backend/internal/repository"
	"log"
	"net/http"
	"path/filepath"

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

	// Seed default admin user
	db.SeedAdmin()

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

	// Add a Custom Request Debugger
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-Backend-Version", "2.1.0")
			log.Printf("DEBUG [%s]: Incoming %s %s", r.RemoteAddr, r.Method, r.URL.Path)
			next.ServeHTTP(w, r)
		})
	})

	// CORS Config
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"}, // Temporarily widen for debug
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Public / Health Routes
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	// Temp: DB diagnostics endpoint
	r.Get("/api/debug-db", func(w http.ResponseWriter, r *http.Request) {
		dbURL := os.Getenv("DATABASE_URL")
		// Mask password for safety
		masked := dbURL
		if len(masked) > 20 {
			masked = masked[:40] + "***MASKED***" + masked[len(masked)-30:]
		}
		var userCount int
		var adminEmail string
		db.DB.QueryRow("SELECT COUNT(*) FROM users").Scan(&userCount)
		db.DB.QueryRow("SELECT COALESCE(email,'none') FROM users WHERE role='ADMIN' LIMIT 1").Scan(&adminEmail)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"db_url_masked": masked,
			"user_count":    userCount,
			"admin_email":   adminEmail,
		})
	})

	// API Routes
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

		// Operator Fleet (using FleetHandler directly under /api/operator)
		r.Get("/operator/vehicles", fleetHandler.ListVehicles)
		r.Post("/operator/vehicles", fleetHandler.CreateVehicle)
		r.Get("/operator/drivers", fleetHandler.ListDrivers)
		r.Post("/operator/drivers", fleetHandler.CreateDriver)

		// Companies
		r.Post("/companies", companyHandler.CreateCompany)
		r.Get("/companies", companyHandler.ListCompanies)
		r.Get("/companies/{id}", companyHandler.GetCompany)

		// Users
		r.Group(func(r chi.Router) {
			r.Use(middleware.AuthMiddleware)
			r.Get("/users", userHandler.HandleListUsers)
		})
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

	// Serve Frontend Static Files & SPA Fallback (only if no route matched)
	workDir, _ := os.Getwd()
	distPath := filepath.Join(workDir, "..", "dist")

	if _, err := os.Stat(filepath.Join(distPath, "index.html")); err != nil {
		log.Printf("Warning: Frontend dist/index.html not found at %s\n", filepath.Join(distPath, "index.html"))
	}

	filesDir := http.Dir(distPath)
	fileServer := http.FileServer(filesDir)

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("DEBUG: Route NOT FOUND: %s %s (Falling back to dist/index.html or static file)", r.Method, r.URL.Path)

		// 1. If it's an API route that didn't match, return JSON error
		if strings.HasPrefix(r.URL.Path, "/api") {
			log.Printf("DEBUG: Route is an /api path, returning 404 JSON error")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{
				"error":  "API route not found",
				"path":   r.URL.Path,
				"method": r.Method,
			})
			return
		}

		// 2. Check if it's a physical file in dist
		path := filepath.Join(distPath, r.URL.Path)
		info, err := os.Stat(path)
		if err == nil && !info.IsDir() {
			fileServer.ServeHTTP(w, r)
			return
		}

		// 3. Otherwise serve index.html for React SPA
		log.Printf("DEBUG: Serving fallback index.html for %s", r.URL.Path)
		http.ServeFile(w, r, filepath.Join(distPath, "index.html"))
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
