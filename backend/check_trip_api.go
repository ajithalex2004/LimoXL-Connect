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
	// Login
	loginBody := map[string]string{
		"email":    "ahmed@rider.ae",
		"password": "Password123!",
	}
	jsonBody, _ := json.Marshal(loginBody)

	resp, err := http.Post("http://localhost:8080/api/auth/login", "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	var loginResp struct {
		Token string `json:"token"`
	}
	json.Unmarshal(body, &loginResp)

	// Fetch RFQs
	req, _ := http.NewRequest("GET", "http://localhost:8080/api/partner/rfqs", nil)
	req.Header.Set("Authorization", "Bearer "+loginResp.Token)

	client := &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, _ = ioutil.ReadAll(resp.Body)

	// Check for TR260001
	if strings.Contains(string(body), "TR260001") {
		fmt.Println("✓ FOUND TR260001 in API Response!")
		// Try to parse and show details
		var trips []map[string]interface{}
		json.Unmarshal(body, &trips)
		for _, t := range trips {
			if ref, ok := t["booking_reference"].(string); ok && ref == "TR260001" {
				fmt.Printf("Details: %+v\n", t)
			}
		}
	} else {
		fmt.Println("❌ TR260001 NOT FOUND in API Response.")
		fmt.Printf("Response Body Length: %d\n", len(body))
		// fmt.Println("Response Preview:", string(body)[:500])
	}
}
