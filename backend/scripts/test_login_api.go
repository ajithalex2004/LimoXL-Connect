package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	url := "http://localhost:8080/api/auth/login"
	payload := map[string]string{
		"email":    "admin@limoxlink.com",
		"password": "admin123",
	}
	body, _ := json.Marshal(payload)

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		fmt.Printf("Login request failed: %v\n", err)
		return
	}
	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)
	fmt.Printf("DEBUG_LOGIN_RESPONSE: %s\n", string(respBody))
}
