package main

import (
	"context"
	"fmt"
	"log"

	"github.com/xtdlib/dbx"
)

func main() {
	ctx := context.Background()

	// Connect to database
	connStr := "postgres://postgres:password@localhost:5432/testdb?sslmode=disable"
	if err := dbx.Connect(ctx, connStr); err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer dbx.Close(ctx)

	// Listen for notifications
	if _, err := dbx.Exec(ctx, "LISTEN test_channel"); err != nil {
		log.Fatalf("Failed to LISTEN: %v", err)
	}
	
	fmt.Println("Listening on 'test_channel'...")
	fmt.Println("\nTo send a notification, run this in psql:")
	fmt.Println("  NOTIFY test_channel, 'Hello World';")
	
	// Wait for notifications
	conn := dbx.GetConn()
	for {
		notification, err := conn.WaitForNotification(ctx)
		if err != nil {
			log.Printf("Error waiting for notification: %v", err)
			continue
		}
		
		fmt.Printf("\nReceived: %s\n", notification.Payload)
	}
}