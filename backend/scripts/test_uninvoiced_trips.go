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

	// Test Uninvoiced Trips endpoint
	client := &http.Client{}
	req, _ := http.NewRequest("GET", "http://localhost:8080/api/partner/trips/completed", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err = client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	body, _ = ioutil.ReadAll(resp.Body)

	fmt.Printf("\nUninvoiced Trips Response:\n")
	fmt.Printf("Status Code: %d\n", resp.StatusCode)
	fmt.Printf("Body: %s\n", string(body))

	var trips []map[string]interface{}
	json.Unmarshal(body, &trips)
	fmt.Printf("\nTotal Uninvoiced Trips: %d\n", len(trips))

	if len(trips) > 0 {
		fmt.Println("\nTrip Details:")
		for i, trip := range trips {
			fmt.Printf("%d. Ref: %v, Status: %v\n", i+1, trip["booking_reference"], trip["status"])
		}
	}
}
