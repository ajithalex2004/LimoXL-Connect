package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
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
		User  struct {
			CompanyID string `json:"company_id"`
		} `json:"user"`
	}
	json.Unmarshal(body, &loginResp)
	token := loginResp.Token

	// Decode JWT to see company_id
	parts := strings.Split(token, ".")
	if len(parts) > 1 {
		// The payload is the second part (index 1)
		fmt.Printf("JWT Token (truncated): %s...\n", token[:50])
	}

	fmt.Printf("User Company ID from login response: %s\n", loginResp.User.CompanyID)
	fmt.Println("\nExpected fulfillment_company_id: 24a8f9b1-5bdd-459e-826c-586946f176e6")

	if loginResp.User.CompanyID == "24a8f9b1-5bdd-459e-826c-586946f176e6" {
		fmt.Println("✓ IDs MATCH - Query should work!")
	} else {
		fmt.Println("❌ IDs DON'T MATCH - This is the problem!")
	}
}
