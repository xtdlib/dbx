# PostgreSQL LISTEN/NOTIFY Example

This example demonstrates how to use PostgreSQL's LISTEN/NOTIFY feature with the dbx library.

## Features Demonstrated

1. **Basic NOTIFY** - Send simple text notifications
2. **JSON Payloads** - Send structured data as JSON
3. **Trigger-based Notifications** - Automatic notifications on database changes
4. **Multiple Channels** - Listen and notify on multiple channels

## Prerequisites

Start PostgreSQL (if not already running):
```bash
docker run --name test-postgres -e POSTGRES_PASSWORD=password -e POSTGRES_DB=testdb -p 5432:5432 -d postgres:15
```

## Running the Example

```bash
go run main.go
```

## Example Output

```
Listening on channel 'events'...

--- Example 1: Basic NOTIFY ---
Notification sent!

--- Example 2: NOTIFY with JSON payload ---
JSON notification sent!

--- Example 3: Trigger-based NOTIFY ---
Inserting row to trigger notification...

--- Example 4: Multiple channels ---
Now listening on channel 'alerts'
Now listening on channel 'updates'

--- Cleanup ---

Demo completed!

Note: To actually receive notifications in real-time, you would need to:
1. Use pgx.Conn directly with conn.WaitForNotification()
2. Or implement a notification handler using channels
3. Have separate connections for LISTEN and NOTIFY (notifications from same connection are delivered immediately)
```

## Real-time Notification Handling

The dbx library provides access to the underlying pgx connection via `dbx.PgConn()`. For real-time notification handling, you would typically:

1. Use separate connections for LISTEN and NOTIFY
2. Implement a goroutine that calls `conn.WaitForNotification()` on the pgx connection
3. Handle notifications through channels

This example focuses on the LISTEN/NOTIFY SQL commands which are the foundation for PostgreSQL's pub/sub system.