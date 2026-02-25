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

	client := &http.Client{}

	// Test all three endpoints
	endpoints := []struct {
		name string
		url  string
	}{
		{"New Requests (ListRFQs)", "http://localhost:8080/api/partner/rfqs"},
		{"Submitted Quotes (ListRFQHistory)", "http://localhost:8080/api/partner/rfqs/history"},
		{"Assigned Trips", "http://localhost:8080/api/partner/trips"},
	}

	for _, ep := range endpoints {
		req, _ := http.NewRequest("GET", ep.url, nil)
		req.Header.Set("Authorization", "Bearer "+token)
		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("❌ %s failed: %v\n", ep.name, err)
			continue
		}
		body, _ := ioutil.ReadAll(resp.Body)

		var trips []map[string]interface{}
		json.Unmarshal(body, &trips)

		fmt.Printf("\n%s:\n", ep.name)
		fmt.Printf("  Status: %d\n", resp.StatusCode)
		fmt.Printf("  Count: %d trips\n", len(trips))
		if len(trips) > 0 {
			fmt.Printf("  Sample: %v\n", trips[0]["reference_no"])
		}
	}
}
