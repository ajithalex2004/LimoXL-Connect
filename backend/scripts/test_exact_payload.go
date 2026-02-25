package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	// Login as Rider Arabia
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

	// Test with EXACT same invoice number as browser
	invoiceData := map[string]interface{}{
		"trip_id":        "f857f1ae-bb78-4ef9-8909-6290d1c96ac3",
		"invoice_number": "INV-2026-RIDER-001",
		"amount":         300,
	}
	jsonBody, _ = json.Marshal(invoiceData)

	client := &http.Client{}
	req, _ := http.NewRequest("POST", "http://localhost:8080/api/partner/invoices", bytes.NewBuffer(jsonBody))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	fmt.Println("\n=== Testing with browser's exact payload ===")
	fmt.Printf("Payload: %s\n", string(jsonBody))

	resp, err = client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	body, _ = ioutil.ReadAll(resp.Body)

	fmt.Printf("\nStatus: %d\n", resp.StatusCode)
	fmt.Printf("Response: %s\n", string(body))

	if resp.StatusCode != 201 {
		fmt.Println("\n❌ FAILED - This might be a duplicate invoice number issue")
		fmt.Println("Check if invoice 'INV-2026-RIDER-001' already exists in the database")
	} else {
		fmt.Println("\n✅ SUCCESS")
	}
}
