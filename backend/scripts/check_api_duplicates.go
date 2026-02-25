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
	// Login
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

	// Get RFQs
	client := &http.Client{}
	req, _ := http.NewRequest("GET", "http://localhost:8080/api/partner/rfqs", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err = client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	body, _ = ioutil.ReadAll(resp.Body)

	var trips []map[string]interface{}
	json.Unmarshal(body, &trips)

	fmt.Printf("Total trips returned: %d\n\n", len(trips))

	// Check for duplicate IDs
	seen := make(map[string]int)
	for _, trip := range trips {
		id := trip["id"].(string)
		seen[id]++
	}

	fmt.Println("=== Checking for duplicate IDs in API response ===")
	foundDupes := false
	for id, count := range seen {
		if count > 1 {
			fmt.Printf("❌ ID %s appears %d times\n", id, count)
			foundDupes = true
		}
	}

	if !foundDupes {
		fmt.Println("✓ No duplicate IDs in API response")
	}

	// Show all trip IDs and references
	fmt.Println("\n=== All trips in response ===")
	for i, trip := range trips {
		fmt.Printf("%d. ID: %s, Ref: %v\n", i+1, trip["id"].(string), trip["booking_reference"])
	}
}
