# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is `github.com/xtdlib/dbx`, a simplified PostgreSQL wrapper library for Go that provides a clean API around pgx/v5 with generic support for struct mapping. The library uses a global connection pattern and offers both low-level pgx functions and high-level generic struct operations.

## Core Architecture

### Connection Management
- **Global Connection**: The library uses a single global `*pgx.Conn` variable (`conn`) in `pgx.go`
- **Connection Functions**: `Connect()`, `MustConnect()`, `ConnectConfig()`, and `Close()`
- All database operations go through this single connection

### Two-Layer API Design

#### Low-Level Layer (`pgx.go`)
- Direct wrappers around pgx/v5 functions: `Query()`, `QueryRow()`, `Exec()`, `Begin()`, etc.
- Type aliases for pgx types: `Row`, `Rows`, `Tx`, `Batch`, etc.
- Maintains the same function signatures as pgx but uses the global connection

#### High-Level Generic Layer (`pgxutil.go`) 
- **`Get[T]()`**: Select single row into struct pointer (`*T`)
- **`Select[T]()`**: Select multiple rows into slice of struct pointers (`[]*T`)
- **`InsertStruct[T]()`**: Insert struct and return inserted row with `RETURNING *`

### Struct Mapping
- **Flexible Field Mapping**: Structs can have fewer fields than database columns - unmapped columns are automatically ignored
- Uses custom `scanRowToStruct` function for resilient column-to-field mapping
- Supports `db:"column_name"` struct tags for custom column mapping
- Falls back to lowercase field names if no `db` tag present
- Skip fields with `db:"-"` tag
- Database columns without corresponding struct fields are scanned into dummy variables

## Development Commands

### Setting Up Local PostgreSQL for Testing
```bash
# Start PostgreSQL in Docker
docker run --name test-postgres -e POSTGRES_PASSWORD=password -e POSTGRES_DB=testdb -p 5432:5432 -d postgres:15

# Create test table and data
PGPASSWORD=password psql -h localhost -U postgres -d testdb -c "
CREATE TABLE holdings (
    ts TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    loc VARCHAR(50) NOT NULL,
    currency VARCHAR(10) NOT NULL,
    amount NUMERIC(20,8) NOT NULL,
    notes TEXT
);
INSERT INTO holdings (ts, loc, currency, amount, notes) VALUES 
    ('2025-01-01 10:00:00+00', 'binance', 'btc', 1.5, 'Initial BTC position'),
    ('2025-01-02 11:00:00+00', 'coinbase', 'eth', 10.25, 'ETH purchase'),
    ('2025-01-03 12:00:00+00', 'kraken', 'sol', 100.0, NULL);
"
```

### Running Examples
```bash
# Run select example (demonstrates Get and Select functions)
cd example/select && go run main.go

# Run insert example (demonstrates InsertStruct and manual insert patterns)
cd example/insert && go run insert.go

# Run flexible field mapping example (struct with fewer fields than DB columns)
cd example && go run flexible-select.go
```

### Building and Testing
```bash
# Build the module
go build ./...

# Run with Go 1.24.3+ (uses modern Go generics)
go version  # Should be 1.24.3 or later
```

## Key Dependencies
- `github.com/jackc/pgx/v5` - PostgreSQL driver and toolkit
- `github.com/jackc/pgxutil` - Additional pgx utilities
- `github.com/jackc/pgx/v5/pgtype` - PostgreSQL type support (e.g., `pgtype.Numeric`)

## Example Usage Patterns

### Connection Setup
```go
ctx := context.Background()
dbx.MustConnect(ctx, "postgres://user:pass@host:port/dbname")
defer dbx.Close(ctx)
```

### Struct Definition
```go
type MyTable struct {
    ID       int            `db:"id"`
    Name     string         `db:"name"`
    Amount   pgtype.Numeric `db:"amount"`
    Notes    *string        `db:"notes"`  // nullable field
}
```

### Generic Operations
```go
// Single row
item, err := dbx.Get[MyTable](ctx, "SELECT * FROM my_table WHERE id = $1", 123)

// Multiple rows  
items, err := dbx.Select[MyTable](ctx, "SELECT * FROM my_table WHERE active = $1", true)

// Insert with automatic RETURNING
inserted, err := dbx.InsertStruct(ctx, "my_table", myStruct)
```

## Important Notes

- All generic functions (`Get`, `Select`, `InsertStruct`) automatically use struct field mapping
- `InsertStruct` always uses `RETURNING *` to return the complete inserted row
- The library assumes a PostgreSQL database (uses pgx-specific features)
- Examples use a specific database connection string - update for your environment
- No built-in connection pooling - uses single connection for simplicity