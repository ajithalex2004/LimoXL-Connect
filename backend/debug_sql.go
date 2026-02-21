package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func main() {
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		connStr = "postgres://postgres:password@localhost:5432/limoxlink?sslmode=disable"
	}
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	partnerID := "24a8f9b1-5bdd-459e-826c-586946f176e6"

	fmt.Printf("Querying for Partner ID: %s\n", partnerID)

	query := `
		SELECT
			t.id, t.reference_no, t.status, t.visibility, ta.company_id
		FROM trips t
		LEFT JOIN trip_access ta ON t.id = ta.trip_id
		WHERE t.reference_no = 'TR260001'
	`
	// Running the query WITHOUT the WHERE filter first to see what rows exist for this trip

	rows, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	fmt.Println("\nRaw Rows for TR260001:")
	for rows.Next() {
		var id, ref, status, vis string
		var compID sql.NullString
		rows.Scan(&id, &ref, &status, &vis, &compID)

		fmt.Printf("Trip: %s | Ref: %s | Status: '%s' | Vis: '%s' | AccessCompID: '%s'\n",
			id, ref, status, vis, compID.String)

		// Check the specific condition logic manually
		matchStatus := (status == "MARKETPLACE_SEARCH" || status == "OFFERED")
		matchVis := (vis == "PUBLIC" || compID.String == partnerID)
		fmt.Printf("  -> Match Status? %v\n", matchStatus)
		fmt.Printf("  -> Match Vis/Access? %v (Public: %v, AccessMatch: %v)\n", matchVis, vis == "PUBLIC", compID.String == partnerID)
		fmt.Printf("  -> Final Result (AND): %v\n", matchStatus && matchVis)
	}
}
