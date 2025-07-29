package dbx

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
)

var conn *pgx.Conn

// Type aliases for commonly used pgx types
type Row = pgx.Row
type Rows = pgx.Rows
type Tx = pgx.Tx
type Batch = pgx.Batch
type BatchResults = pgx.BatchResults
type CommandTag = pgconn.CommandTag
type CopyFromSource = pgx.CopyFromSource
type Identifier = pgx.Identifier

// Connect initializes the package-level connection
func Connect(ctx context.Context, connString string) error {
	var err error
	conn, err = pgx.Connect(ctx, connString)
	return err
}

func MustConnect(ctx context.Context, connString string) {
	if err := Connect(ctx, connString); err != nil {
		panic(err)
	}
}

// ConnectConfig initializes the package-level connection with config
func ConnectConfig(ctx context.Context, config *pgx.ConnConfig) error {
	var err error
	conn, err = pgx.ConnectConfig(ctx, config)
	return err
}

// Close closes the database connection
func Close(ctx context.Context) error {
	if conn != nil {
		return conn.Close(ctx)
	}
	return nil
}

// Query executes a query that returns rows
func Query(ctx context.Context, sql string, args ...any) (Rows, error) {
	return conn.Query(ctx, sql, args...)
}

// QueryRow executes a query that is expected to return at most one row
func QueryRow(ctx context.Context, sql string, args ...any) Row {
	return conn.QueryRow(ctx, sql, args...)
}

// Exec executes a query without returning any rows
func Exec(ctx context.Context, sql string, args ...any) (CommandTag, error) {
	return conn.Exec(ctx, sql, args...)
}

// Ping verifies the connection to the database is still alive
func Ping(ctx context.Context) error {
	return conn.Ping(ctx)
}

// Begin starts a transaction
func Begin(ctx context.Context) (Tx, error) {
	return conn.Begin(ctx)
}

// BeginTx starts a transaction with options
func BeginTx(ctx context.Context, txOptions pgx.TxOptions) (Tx, error) {
	return conn.BeginTx(ctx, txOptions)
}

// CopyFrom performs a copy from operation
func CopyFrom(ctx context.Context, tableName Identifier, columnNames []string, rowSrc CopyFromSource) (int64, error) {
	return conn.CopyFrom(ctx, tableName, columnNames, rowSrc)
}

// SendBatch sends a batch of queries
func SendBatch(ctx context.Context, b *Batch) BatchResults {
	return conn.SendBatch(ctx, b)
}

// Config returns the current connection config
func Config() *pgx.ConnConfig {
	if conn != nil {
		return conn.Config()
	}
	return nil
}

// IsClosed reports whether the connection is closed
func IsClosed() bool {
	if conn == nil {
		return true
	}
	return conn.IsClosed()
}

// Prepare creates a prepared statement
func Prepare(ctx context.Context, name, sql string) (*pgconn.StatementDescription, error) {
	return conn.Prepare(ctx, name, sql)
}

// Deallocate deallocates a prepared statement
func Deallocate(ctx context.Context, name string) error {
	return conn.Deallocate(ctx, name)
}

// LoadType loads a composite type definition
func LoadType(ctx context.Context, typeName string) (*pgtype.Type, error) {
	return conn.LoadType(ctx, typeName)
}

// TypeMap returns the connection's type map
func TypeMap() *pgtype.Map {
	if conn != nil {
		return conn.TypeMap()
	}
	return nil
}

// PgConn returns the underlying pgconn.Conn
func PgConn() *pgconn.PgConn {
	if conn != nil {
		return conn.PgConn()
	}
	return nil
}

// GetConn returns the underlying pgx.Conn for advanced operations like notifications
func GetConn() *pgx.Conn {
	return conn
}
