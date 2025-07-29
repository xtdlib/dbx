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
	if err := dbx.Connect(ctx, "postgres://postgres:password@localhost:5432/testdb?sslmode=disable"); err != nil {
		log.Fatal("Failed to connect:", err)
	}
	defer dbx.Close(ctx)

	// Setup LISTEN on a channel
	channel := "events"
	if _, err := dbx.Exec(ctx, "LISTEN "+channel); err != nil {
		log.Fatalf("Failed to LISTEN: %v", err)
	}
	fmt.Printf("Listening on channel '%s'...\n", channel)

	// Example 1: Send notification from same connection
	fmt.Println("\n--- Example 1: Basic NOTIFY ---")
	if _, err := dbx.Exec(ctx, "NOTIFY events, 'Hello from dbx!'"); err != nil {
		log.Printf("Failed to NOTIFY: %v", err)
	}
	fmt.Println("Notification sent!")

	// Example 2: Send notification with JSON payload
	fmt.Println("\n--- Example 2: NOTIFY with JSON payload ---")
	payload := `{"id": 123, "action": "update", "timestamp": "2025-01-29T10:00:00Z"}`
	query := fmt.Sprintf("NOTIFY %s, '%s'", channel, payload)
	if _, err := dbx.Exec(ctx, query); err != nil {
		log.Printf("Failed to NOTIFY with JSON: %v", err)
	}
	fmt.Println("JSON notification sent!")

	// Example 3: Using a trigger (create table and trigger for demo)
	fmt.Println("\n--- Example 3: Trigger-based NOTIFY ---")
	
	// Create a demo table
	_, err := dbx.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS notify_demo (
			id SERIAL PRIMARY KEY,
			message TEXT,
			created_at TIMESTAMP DEFAULT NOW()
		)
	`)
	if err != nil {
		log.Printf("Failed to create table: %v", err)
	}

	// Create trigger function
	_, err = dbx.Exec(ctx, `
		CREATE OR REPLACE FUNCTION notify_trigger() RETURNS trigger AS $$
		BEGIN
			PERFORM pg_notify('events', 
				json_build_object(
					'table', TG_TABLE_NAME,
					'action', TG_OP,
					'id', NEW.id,
					'message', NEW.message
				)::text
			);
			RETURN NEW;
		END;
		$$ LANGUAGE plpgsql
	`)
	if err != nil {
		log.Printf("Failed to create trigger function: %v", err)
	}

	// Create trigger
	_, err = dbx.Exec(ctx, `
		DROP TRIGGER IF EXISTS notify_demo_trigger ON notify_demo;
		CREATE TRIGGER notify_demo_trigger
		AFTER INSERT ON notify_demo
		FOR EACH ROW EXECUTE FUNCTION notify_trigger()
	`)
	if err != nil {
		log.Printf("Failed to create trigger: %v", err)
	}

	// Insert a row to trigger notification
	fmt.Println("Inserting row to trigger notification...")
	_, err = dbx.Exec(ctx, "INSERT INTO notify_demo (message) VALUES ($1)", "Triggered notification!")
	if err != nil {
		log.Printf("Failed to insert: %v", err)
	}

	// Example 4: Multiple channels
	fmt.Println("\n--- Example 4: Multiple channels ---")
	channels := []string{"alerts", "updates"}
	for _, ch := range channels {
		if _, err := dbx.Exec(ctx, "LISTEN "+ch); err != nil {
			log.Printf("Failed to LISTEN on %s: %v", ch, err)
		} else {
			fmt.Printf("Now listening on channel '%s'\n", ch)
		}
	}

	// Send to different channels
	dbx.Exec(ctx, "NOTIFY alerts, 'System alert!'")
	dbx.Exec(ctx, "NOTIFY updates, 'Data updated!'")

	// Clean up
	fmt.Println("\n--- Cleanup ---")
	dbx.Exec(ctx, "DROP TABLE IF EXISTS notify_demo CASCADE")
	dbx.Exec(ctx, "DROP FUNCTION IF EXISTS notify_trigger() CASCADE")
	
	// Unlisten
	dbx.Exec(ctx, "UNLISTEN events")
	for _, ch := range channels {
		dbx.Exec(ctx, "UNLISTEN "+ch)
	}

	fmt.Println("\nDemo completed!")
	fmt.Println("\nNote: To actually receive notifications in real-time, you would need to:")
	fmt.Println("1. Use pgx.Conn directly with conn.WaitForNotification()")
	fmt.Println("2. Or implement a notification handler using channels")
	fmt.Println("3. Have separate connections for LISTEN and NOTIFY (notifications from same connection are delivered immediately)")
}