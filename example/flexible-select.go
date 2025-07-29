package main

import (
	"context"
	"log"
	"time"

	"github.com/xtdlib/dbx"
)

// Holdings struct with Amount field commented out to demonstrate flexible mapping
type Holdings struct {
	Ts       time.Time `db:"ts"`
	Loc      string    `db:"loc"`
	Currency string    `db:"currency"`
	// Amount   pgtype.Numeric `db:"amount"`  // Intentionally omitted
	Notes *string `db:"notes"`
}

func main() {
	ctx := context.Background()
	dbx.MustConnect(ctx, "postgres://postgres:postgres@oci-aca-001:3030/postgres?sslmode=disable")
	defer dbx.Close(ctx)

	log.Println("Testing flexible field mapping - Amount field is omitted from struct")
	
	// Get a single row (even though table has Amount column)
	holding, err := dbx.Get[Holdings](ctx, "SELECT * FROM holdings LIMIT 1")
	if err != nil {
		panic(err)
	}
	log.Printf("Single holding (no Amount): %+v\n", holding)

	// Select multiple rows
	holdings, err := dbx.Select[Holdings](ctx, "SELECT * FROM holdings ORDER BY ts DESC LIMIT 3")
	if err != nil {
		panic(err)
	}

	log.Printf("\nFound %d holdings (Amount column ignored):\n", len(holdings))
	for i, h := range holdings {
		log.Printf("%d. %s: %s at %s\n", i+1, h.Loc, h.Currency, h.Ts.Format("2006-01-02 15:04:05"))
	}
}
