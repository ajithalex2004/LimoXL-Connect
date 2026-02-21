package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	// Test Accept Quote endpoint
	// First, get a quote ID from the list
	resp, err := http.Get("http://localhost:8080/api/operator/quotes")
	if err != nil {
		log.Fatal("Failed to get quotes:", err)
	}

	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Printf("Quotes Response: %s\n\n", string(body))

	var quotes []map[string]interface{}
	json.Unmarshal(body, &quotes)

	if len(quotes) == 0 {
		fmt.Println("No quotes available to test")
		return
	}

	quoteID := quotes[0]["id"].(string)
	fmt.Printf("Testing Accept Quote for ID: %s\n", quoteID)

	// Test Accept
	url := fmt.Sprintf("http://localhost:8080/api/operator/quotes/%s/accept", quoteID)
	req, _ := http.NewRequest("POST", url, nil)

	client := &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		log.Fatal("Accept request failed:", err)
	}

	body, _ = ioutil.ReadAll(resp.Body)
	fmt.Printf("\nAccept Response Status: %d\n", resp.StatusCode)
	fmt.Printf("Accept Response Body: %s\n", string(body))
}
