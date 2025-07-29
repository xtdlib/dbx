package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
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

	// Listen for notifications on channel "test_channel"
	if _, err := dbx.Exec(ctx, "LISTEN test_channel"); err != nil {
		log.Fatalf("Failed to LISTEN: %v", err)
	}
	fmt.Println("Listening for notifications on channel 'test_channel'...")

	// Start a goroutine to send notifications every 2 seconds
	go func() {
		for i := 1; ; i++ {
			time.Sleep(2 * time.Second)
			payload := fmt.Sprintf("Hello from notification #%d", i)
			query := fmt.Sprintf("NOTIFY test_channel, '%s'", payload)
			if _, err := dbx.Exec(ctx, query); err != nil {
				log.Printf("Failed to send notification: %v", err)
			} else {
				fmt.Printf("Sent: %s\n", payload)
			}
		}
	}()

	// Handle Ctrl+C gracefully
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)

	// Main loop to receive notifications
	for {
		select {
		case <-sigChan:
			fmt.Println("\nStopping listener...")
			return
		default:
			// Wait for notification with timeout
			ctx, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
			notification, err := waitForNotification(ctx)
			cancel()

			if err == nil && notification != nil {
				fmt.Printf("Received notification: Channel=%s, Payload=%s\n", 
					notification.Channel, notification.Payload)
			}
		}
	}
}

// waitForNotification waits for a notification from PostgreSQL
func waitForNotification(ctx context.Context) (*pgconn.Notification, error) {
	// Get the underlying pgx connection
	conn := dbx.GetConn()
	if conn == nil {
		return nil, fmt.Errorf("no connection available")
	}

	// Wait for notification
	return conn.WaitForNotification(ctx)
}