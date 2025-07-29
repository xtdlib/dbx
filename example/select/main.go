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
	dbx.MustConnect(ctx, "postgres://postgres:password@localhost:5432/testdb?sslmode=disable")

	// err := dbx.QueryRow(context.Background(), "select * from holdings").Scan(&holdings)
	// if err != nil {
	// 	fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
	// 	os.Exit(1)
	// }

	// Get a single row
	holding, err := dbx.Get[Holdings](context.Background(), "select * from holdings limit 1;")
	if err != nil {
		panic(err)
	}
	log.Printf("Single holding: %+v\n", holding)

	// Select multiple rows
	holdings, err := dbx.Select[Holdings](context.Background(), "select * from holdings where currency in ('btc', 'eth', 'sol') order by ts desc limit 5;")
	if err != nil {
		panic(err)
	}

	log.Printf("\nFound %d holdings:\n", len(holdings))
	for i, h := range holdings {
		amount, _ := h.Amount.Float64Value()
		log.Printf("%d. %s: %s (%.2f) at %s\n", i+1, h.Loc, h.Currency, amount.Float64, h.Ts.Format("2006-01-02 15:04:05"))
	}
}
