package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

func main() {
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		connStr = "postgres://postgres:password@localhost:5432/limoxlink?sslmode=disable"
	}
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// 1. Create a fresh RFQ
	tripID := uuid.New()
	_, err = db.Exec(`
		INSERT INTO trips (
			id, reference_no, status, visibility, created_at, updated_at, 
			pickup_zone, dropoff_zone, pickup_time, passenger_name, passenger_phone,
			requesting_company_id, vehicle_type_requested, price,
			pickup_location, dropoff_location
		) VALUES (
			$1, 'TEST-QUOTE-1', 'MARKETPLACE_SEARCH', 'PUBLIC', NOW(), NOW(),
			'Dubai Marina', 'DXB Airport', NOW() + INTERVAL '1 day', 'Quote Tester', '+971500000000',
			$2, 'Sedan', 0,
			ST_SetSRID(ST_MakePoint(55.27, 25.20), 4326),
			ST_SetSRID(ST_MakePoint(55.36, 25.25), 4326)
		)`,
		tripID, "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a00") // Operator ID
	if err != nil {
		log.Fatal("Failed to create test trip:", err)
	}
	fmt.Println("✓ Created Test Trip:", tripID)

	// 2. Log in as Ratio Arabia
	loginBody := map[string]string{
		"email":    "ahmed@rider.ae",
		"password": "Password123!",
	}
	jsonBody, _ := json.Marshal(loginBody)
	resp, err := http.Post("http://localhost:8080/api/auth/login", "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		log.Fatal("Login failed:", err)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	var loginResp struct {
		Token string `json:"token"`
	}
	json.Unmarshal(body, &loginResp)
	token := loginResp.Token
	fmt.Println("✓ Logged in as Rider Arabia")

	// 3. Submit Quote
	client := &http.Client{Timeout: 10 * time.Second}
	quoteReq := map[string]interface{}{
		"trip_id": tripID.String(),
		"price":   150.50,
		"notes":   "Best price for you",
	}
	jsonQuote, _ := json.Marshal(quoteReq)
	req, _ := http.NewRequest("POST", "http://localhost:8080/api/partner/quotes", bytes.NewBuffer(jsonQuote))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err = client.Do(req)
	if err != nil {
		log.Fatal("Submit Quote failed:", err)
	}

	if resp.StatusCode != 200 {
		body, _ = ioutil.ReadAll(resp.Body)
		log.Fatalf("Submit Quote Failed: %d\nBody: %s", resp.StatusCode, string(body))
	}
	fmt.Println("✓ Quote API Call Successful")

	// 4. Verify DB State
	var status string
	err = db.QueryRow("SELECT status FROM trips WHERE id = $1", tripID).Scan(&status)
	if err != nil {
		log.Fatal("Failed to fetch trip status:", err)
	}
	if status == "OFFERED" {
		fmt.Println("✓ Trip Status updated to OFFERED")
	} else {
		fmt.Printf("❌ Trip Status is %s (Expected OFFERED)\n", status)
	}

	var offerCount int
	err = db.QueryRow("SELECT count(*) FROM trip_offers WHERE trip_id = $1 AND status = 'PENDING'", tripID).Scan(&offerCount)
	if offerCount == 1 {
		fmt.Println("✓ Trip Offer created")
	} else {
		fmt.Printf("❌ Offer count is %d (Expected 1)\n", offerCount)
	}
}
