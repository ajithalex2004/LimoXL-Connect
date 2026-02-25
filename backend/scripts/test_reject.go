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
	// 1. Login as Rider Arabia
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

	// 2. List RFQs to get a trip ID
	client := &http.Client{}
	req, _ := http.NewRequest("GET", "http://localhost:8080/api/partner/rfqs", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err = client.Do(req)
	if err != nil {
		log.Fatal("List RFQs failed:", err)
	}
	body, _ = ioutil.ReadAll(resp.Body)

	var trips []map[string]interface{}
	json.Unmarshal(body, &trips)

	if len(trips) == 0 {
		fmt.Println("❌ No trips available to test reject")
		return
	}

	tripID := trips[0]["id"].(string)
	fmt.Printf("✓ Found Trip ID: %s\n", tripID)

	// 3. Test Reject
	rejectReq := map[string]string{"trip_id": tripID}
	jsonReject, _ := json.Marshal(rejectReq)
	req, _ = http.NewRequest("POST", "http://localhost:8080/api/partner/reject", bytes.NewBuffer(jsonReject))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err = client.Do(req)
	if err != nil {
		log.Fatal("Reject API failed:", err)
	}

	body, _ = ioutil.ReadAll(resp.Body)
	fmt.Printf("Reject Response: %d\nBody: %s\n", resp.StatusCode, string(body))

	if resp.StatusCode == 200 {
		fmt.Println("✓ Reject API Call Successful")
	} else {
		fmt.Printf("❌ Reject Failed with status %d\n", resp.StatusCode)
	}
}
