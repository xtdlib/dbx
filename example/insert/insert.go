package main

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/xtdlib/dbx"
)

type Holdings struct {
	Ts       time.Time      `db:"ts"`
	Loc      string         `db:"loc"`
	Currency string         `db:"currency"`
	Amount   pgtype.Numeric `db:"amount"`
	Notes    *string        `db:"notes"`
}

func main() {
	ctx := context.Background()
	
	// Connect to database
	dbx.MustConnect(ctx, "postgres://postgres:password@localhost:5432/testdb?sslmode=disable")
	defer dbx.Close(ctx)

	// Create a new holding record
	notes := "Example holding inserted via dbx"
	amount := pgtype.Numeric{}
	amount.Scan("1000.50") // 1000.50 units

	// Insert using simple Exec
	result, err := dbx.Exec(ctx, 
		`INSERT INTO holdings (ts, loc, currency, amount, notes) 
		 VALUES ($1, $2, $3, $4, $5)`,
		time.Now(),
		"binance",
		"btc",
		amount,
		notes,
	)
	if err != nil {
		log.Fatalf("Failed to insert: %v", err)
	}

	log.Printf("Inserted %d row(s)\n", result.RowsAffected())

	// Insert another record using RETURNING to get the inserted data back
	var insertedHolding Holdings
	err = dbx.QueryRow(ctx,
		`INSERT INTO holdings (ts, loc, currency, amount, notes) 
		 VALUES ($1, $2, $3, $4, $5)
		 RETURNING *`,
		time.Now(),
		"coinbase",
		"eth",
		amount,
		nil, // NULL notes
	).Scan(
		&insertedHolding.Ts,
		&insertedHolding.Loc,
		&insertedHolding.Currency,
		&insertedHolding.Amount,
		&insertedHolding.Notes,
	)
	if err != nil {
		log.Fatalf("Failed to insert with RETURNING: %v", err)
	}

	log.Printf("Inserted holding: %+v\n", insertedHolding)

	// Example using the new InsertStruct method
	log.Println("\n--- Using InsertStruct method ---")
	
	// Create a new holding using struct
	newHolding := Holdings{
		Ts:       time.Now(),
		Loc:      "kraken",
		Currency: "sol",
		Amount:   amount,
		Notes:    nil,
	}
	
	// Insert always returns the inserted row
	inserted1, err := dbx.InsertStruct(ctx, "holdings", newHolding)
	if err != nil {
		log.Fatalf("Failed to InsertStruct: %v", err)
	}
	log.Printf("Inserted holding: %+v\n", inserted1)
	
	// Insert another record
	anotherNotes := "Another holding"
	anotherHolding := Holdings{
		Ts:       time.Now(),
		Loc:      "okx",
		Currency: "matic",
		Amount:   amount,
		Notes:    &anotherNotes,
	}
	
	inserted2, err := dbx.InsertStruct(ctx, "holdings", anotherHolding)
	if err != nil {
		log.Fatalf("Failed to InsertStruct: %v", err)
	}
	log.Printf("Inserted another holding: %+v\n", inserted2)

	// Verify by selecting the records we just inserted
	rows, err := dbx.Query(ctx, 
		`SELECT * FROM holdings 
		 WHERE loc IN ('binance', 'coinbase', 'kraken', 'okx') 
		 ORDER BY ts DESC 
		 LIMIT 4`)
	if err != nil {
		log.Fatalf("Failed to query: %v", err)
	}
	defer rows.Close()

	log.Println("\nRecently inserted holdings:")
	for rows.Next() {
		var h Holdings
		err := rows.Scan(&h.Ts, &h.Loc, &h.Currency, &h.Amount, &h.Notes)
		if err != nil {
			log.Fatalf("Failed to scan row: %v", err)
		}
		log.Printf("  %+v\n", h)
	}

	if err := rows.Err(); err != nil {
		log.Fatalf("Row iteration error: %v", err)
	}
}