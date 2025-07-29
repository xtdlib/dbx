package main

import (
	"context"
	"log"

	"github.com/xtdlib/dbx"
)

type Tag string

func main() {
	ctx := context.Background()
	dbx.MustConnect(ctx, "postgres://postgres:postgres@oci-aca-001:3030/postgres?sslmode=disable")

	var tags []string

	// result, err := dbx.Query(ctx, "select loc from holdings limit 2")
	err := dbx.QueryRow(ctx, "select array['foo', 'bar']").Scan(&tags)
	if err != nil {
		panic(err)
	}

	log.Println(tags)
}
