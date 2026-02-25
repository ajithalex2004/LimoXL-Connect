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
	// 1. Login
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
	fmt.Println("Logged in as Rider Arabia. Token len:", len(token))

	// 2. List RFQs
	client := &http.Client{}
	req, _ := http.NewRequest("GET", "http://localhost:8080/api/partner/rfqs", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err = client.Do(req)
	if err != nil {
		log.Fatal("List RFQs failed:", err)
	}
	body, _ = ioutil.ReadAll(resp.Body)

	fmt.Printf("\nAPI Response Code: %d\n", resp.StatusCode)

	// Print raw body snippet
	fmt.Printf("Body Snippet: %s\n", string(body))

	var trips []map[string]interface{}
	json.Unmarshal(body, &trips)
	fmt.Printf("\nFound %d trips in API response:\n", len(trips))
	for _, t := range trips {
		fmt.Printf("- Ref: %v | ID: %v | Status: %v\n", t["reference_no"], t["id"], t["status"])
	}
}
