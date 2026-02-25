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
	"strings"

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

	// Cleanup potential leftovers
	db.Exec("DELETE FROM trips WHERE reference_no IN ('TEST-REJECT-2', 'TEST-ACCEPT-2')")
	db.Exec("DELETE FROM trip_offers WHERE trip_id IN (SELECT id FROM trips WHERE reference_no IN ('TEST-REJECT-2', 'TEST-ACCEPT-2'))")
	// Also need to cleanup offers by trip ID if we don't know the ID?
	// The offers are linked to trip_id. The DELETE FROM trips might cascade if configured, but safe to delete manually if not.
	// Actually trip_offers has FK to trips. If no Cascade, we must delete offers first.
	// But since we generate new UUIDs each run, we can't delete offers by ID easily unless we query trips first.
	// Let's do a robust cleanup.

	// 1. Create Mock Trips directly in DB
	trip1 := uuid.New()
	trip2 := uuid.New()

	_, err = db.Exec(`
		INSERT INTO trips (id, reference_no, status, visibility, created_at, updated_at, rfq_number, pickup_time, pickup_location, dropoff_location, vehicle_type_requested, passenger_name, pickup_zone, dropoff_zone)
		VALUES 
		($1, 'TEST-REJECT-2', 'MARKETPLACE_SEARCH', 'PUBLIC', NOW(), NOW(), 'RFQ-TEST-1', NOW(), 'POINT(55.27 25.20)'::geography, 'POINT(55.30 25.25)'::geography, 'Sedan', 'Test Passenger', 'Zone A', 'Zone B'),
		($2, 'TEST-ACCEPT-2', 'MARKETPLACE_SEARCH', 'PUBLIC', NOW(), NOW(), 'RFQ-TEST-2', NOW(), 'POINT(55.27 25.20)'::geography, 'POINT(55.30 25.25)'::geography, 'Sedan', 'Test Passenger', 'Zone A', 'Zone B')
	`, trip1, trip2)
	if err != nil {
		log.Fatal("Error creating test trips:", err)
	}
	fmt.Println("✓ Created 2 mock trips: TEST-REJECT and TEST-ACCEPT")

	// 2. Login
	loginBody := map[string]string{
		"email":    "ahmed@rider.ae",
		"password": "Password123!",
	}
	jsonBody, _ := json.Marshal(loginBody)
	resp, _ := http.Post("http://localhost:8080/api/auth/login", "application/json", bytes.NewBuffer(jsonBody))
	body, _ := ioutil.ReadAll(resp.Body)
	var loginResp struct {
		Token string `json:"token"`
	}
	json.Unmarshal(body, &loginResp)
	token := loginResp.Token
	fmt.Println("✓ Logged in as Rider Arabia")

	client := &http.Client{}

	// 3. Reject TEST-REJECT
	fmt.Println("\nTesting REJECT...")
	rejectBody := map[string]string{"trip_id": trip1.String()}
	b, _ := json.Marshal(rejectBody)
	req, _ := http.NewRequest("POST", "http://localhost:8080/api/partner/reject", bytes.NewBuffer(b))
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err = client.Do(req)
	if err != nil || resp.StatusCode != 200 {
		fmt.Printf("❌ Reject Failed: %v %v\n", err, resp.Status)
	} else {
		fmt.Println("✓ Reject API Call Successful")
	}

	// 4. Accept TEST-ACCEPT
	fmt.Println("\nTesting ACCEPT...")
	acceptBody := map[string]string{"trip_id": trip2.String()}
	b, _ = json.Marshal(acceptBody)
	req, _ = http.NewRequest("POST", "http://localhost:8080/api/partner/accept", bytes.NewBuffer(b))
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err = client.Do(req)
	if err != nil || resp.StatusCode != 200 {
		fmt.Printf("❌ Accept Failed: %v %v\n", err, resp.Status)
	} else {
		fmt.Println("✓ Accept API Call Successful")
	}

	// 5. Verify Lists
	fmt.Println("\nVerifying Visibility...")

	// List Open RFQs - Should contain neither
	req, _ = http.NewRequest("GET", "http://localhost:8080/api/partner/rfqs", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, _ = client.Do(req)
	body, _ = ioutil.ReadAll(resp.Body)
	listStr := string(body)

	if strings.Contains(listStr, "TEST-REJECT-2") {
		fmt.Println("❌ TEST-REJECT-2 still visible in RFQs!")
	} else {
		fmt.Println("✓ TEST-REJECT-2 is GONE from RFQs")
	}

	if strings.Contains(listStr, "TEST-ACCEPT-2") {
		fmt.Println("❌ TEST-ACCEPT-2 still visible in RFQs!")
	} else {
		fmt.Println("✓ TEST-ACCEPT-2 is GONE from RFQs")
	}

	// List Assigned Trips - Should contain TEST-ACCEPT-2
	req, _ = http.NewRequest("GET", "http://localhost:8080/api/partner/trips", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, _ = client.Do(req)
	body, _ = ioutil.ReadAll(resp.Body)
	listStr = string(body)

	if strings.Contains(listStr, "TEST-ACCEPT-2") {
		fmt.Println("✓ TEST-ACCEPT-2 is VISIBLE in Assigned Trips")
	} else {
		fmt.Println("❌ TEST-ACCEPT-2 NOT FOUND in Assigned Trips!")
		fmt.Println("Response Preview:", listStr[:min(len(listStr), 500)])
	}

	// Clean up
	db.Exec("DELETE FROM trips WHERE id IN ($1, $2)", trip1, trip2)
	db.Exec("DELETE FROM trip_offers WHERE trip_id IN ($1, $2)", trip1, trip2)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
