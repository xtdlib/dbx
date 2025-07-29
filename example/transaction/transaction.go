package main

import (
	"context"
	"log"

	"github.com/xtdlib/dbx"
)

func main() {
	ctx := context.Background()

	// Connect to database
	dbx.MustConnect(ctx, "postgres://postgres:password@localhost:5432/testdb?sslmode=disable")
	defer dbx.Close(ctx)

	log.Println("=== Simple Transaction Example ===")

	// Begin transaction
	tx, err := dbx.Begin(ctx)
	if err != nil {
		log.Fatalf("Failed to begin transaction: %v", err)
	}

	// Insert two records in the same transaction
	_, err = tx.Exec(ctx, "INSERT INTO holdings (loc, currency, amount) VALUES ($1, $2, $3)", "test-loc", "btc", 1.0)
	if err != nil {
		tx.Rollback(ctx)
		log.Fatalf("Failed to insert first record: %v", err)
	}
	log.Println("Inserted first record")




	err = tx.Rollback(ctx)
	if err != nil {
		panic(err)
	}

	// _, err = tx.Exec(ctx, "INSERT INTO holdings (loc, currency, amount) VALUES ($1, $2, $3)", "test-loc", "eth", 2.0)
	// if err != nil {
	// 	tx.Rollback(ctx)
	// 	log.Fatalf("Failed to insert second record: %v", err)
	// }
	// log.Println("Inserted second record")

	// Commit transaction
	err = tx.Commit(ctx)
	if err != nil {
		log.Fatalf("Failed to commit transaction: %v", err)
	}

	log.Println("Transaction committed successfully")
}

