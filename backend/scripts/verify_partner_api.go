package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	// 1. Reset Password first to ensure we can login
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		connStr = "postgres://postgres:password@localhost:5432/limoxlink?sslmode=disable"
	}
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	password := "Password123!"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	// Update Rider Arabia user
	_, err = db.Exec(`
		UPDATE users 
		SET password_hash = $1 
		WHERE email = 'ahmed@rider.ae'
	`, string(hashedPassword))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("✓ Password reset for ahmed@rider.ae to Password123!")

	// 2. Perform Login via API
	fmt.Println("\nAttempting Login...")
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
	if resp.StatusCode != 200 {
		log.Fatalf("Login Failed: %s\nBody: %s", resp.Status, string(body))
	}

	var loginResp struct {
		Token string `json:"token"`
		User  struct {
			CompanyID string `json:"company_id"`
		} `json:"user"`
	}
	json.Unmarshal(body, &loginResp)
	fmt.Printf("✓ Login Successful! Token obtained.\n")
	fmt.Printf("  User Company ID from Response: %s\n", loginResp.User.CompanyID)

	// 3. Fetch RFQs
	fmt.Println("\nFetching RFQs...")
	req, _ := http.NewRequest("GET", "http://localhost:8080/api/partner/rfqs", nil)
	req.Header.Set("Authorization", "Bearer "+loginResp.Token)

	client := &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, _ = ioutil.ReadAll(resp.Body)
	fmt.Printf("Status: %s\n", resp.Status)
	fmt.Printf("Body: %s\n", string(body))
}
