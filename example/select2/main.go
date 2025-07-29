package main

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/xtdlib/dbx"
)

type Holding struct {
	Ts       time.Time      `db:"ts"`
	Loc      string         `db:"loc"`
	Currency string         `db:"currency"`
	Amount   pgtype.Numeric `db:"amount"`
	Notes    *string        `db:"notes"`
}

type Tag string

func main() {
	ctx := context.Background()
	dbx.MustConnect(ctx, "postgres://postgres:postgres@oci-aca-001:3030/postgres?sslmode=disable")

	var tags []Tag
	_ = tags

	// result, err := dbx.Query(ctx, "select loc from holdings limit 2")
	result, err := dbx.Query(ctx, "select array['foo', 'bar']")
	err = result.Scan(&tags)
	if err != nil {
		panic(err)
	}

	// // Get a single row
	// holding, err := dbx.Get[Holdings](context.Background(), "select * from holdings limit 1;")
	// if err != nil {
	// 	panic(err)
	// }
	// log.Printf("Single holding: %+v\n", holding)
	//
	// // Select multiple rows
	// holdings, err := dbx.Select[Holdings](context.Background(), "select * from holdings where currency in ('btc', 'eth', 'sol') order by ts desc limit 5;")
	// if err != nil {
	// 	panic(err)
	// }
}
