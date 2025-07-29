package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
)

func main() {
	ctx := context.Background()

	// Connect to database (separate connection for inserting)
	connStr := "postgres://postgres:password@localhost:5432/testdb?sslmode=disable"
	conn, err := pgx.Connect(ctx, connStr)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close(ctx)

	// Insert messages that will trigger notifications
	for i := 1; i <= 5; i++ {
		message := fmt.Sprintf("Test message #%d at %s", i, time.Now().Format("15:04:05"))
		
		_, err := conn.Exec(ctx, 
			"INSERT INTO notification_test (message) VALUES ($1)",
			message)
		
		if err != nil {
			log.Printf("Failed to insert: %v", err)
		} else {
			fmt.Printf("Inserted: %s\n", message)
		}
		
		time.Sleep(2 * time.Second)
	}
	
	fmt.Println("Done inserting messages")
}