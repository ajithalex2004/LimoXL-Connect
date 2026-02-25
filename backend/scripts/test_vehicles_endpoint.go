package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	// Test the vehicles endpoint
	resp, err := http.Get("http://localhost:8080/api/operator/vehicles")
	if err != nil {
		log.Fatal("Request failed:", err)
	}

	body, _ := ioutil.ReadAll(resp.Body)

	fmt.Printf("Status Code: %d\n", resp.StatusCode)
	fmt.Printf("Response Body: %s\n", string(body))

	if resp.StatusCode != 200 {
		fmt.Println("\n❌ Error detected")
	} else {
		var vehicles []map[string]interface{}
		json.Unmarshal(body, &vehicles)
		fmt.Printf("\n✓ Found %d vehicles\n", len(vehicles))
	}
}
