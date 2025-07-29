#!/bin/bash

# PostgreSQL Test Runner for dbx library
# This script sets up PostgreSQL in Docker and runs all examples

set -e  # Exit on any error

CONTAINER_NAME="dbx-test-postgres"
DB_PASSWORD="password"
DB_NAME="testdb"
DB_PORT="5432"

echo "🚀 Starting PostgreSQL test environment..."

# Function to cleanup
cleanup() {
    echo "🧹 Cleaning up..."
    docker stop $CONTAINER_NAME >/dev/null 2>&1 || true
    docker rm $CONTAINER_NAME >/dev/null 2>&1 || true
}

# Cleanup on script exit
trap cleanup EXIT

# Start PostgreSQL container
echo "📦 Starting PostgreSQL container..."
docker run --name $CONTAINER_NAME \
    -e POSTGRES_PASSWORD=$DB_PASSWORD \
    -e POSTGRES_DB=$DB_NAME \
    -p $DB_PORT:5432 \
    -d postgres:15

# Wait for PostgreSQL to be ready
echo "⏳ Waiting for PostgreSQL to be ready..."
sleep 8

# Create test table
echo "🗄️  Creating test table..."
PGPASSWORD=$DB_PASSWORD psql -h localhost -U postgres -d $DB_NAME -c "
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

echo "✅ Database setup complete!"
echo ""

# Run examples
echo "🧪 Running examples..."
echo ""

# Test select example
echo "--- Testing Select Example ---"
cd ../example/select && go run main.go
echo ""

# Test insert example  
echo "--- Testing Insert Example ---"
cd ../insert && go run insert.go
echo ""

# Test flexible select example
echo "--- Testing Flexible Select Example ---"
cd .. && go run flexible-select.go
echo ""

# Test transaction example
echo "--- Testing Transaction Example ---"
cd transaction && go run transaction.go
echo ""

# Show final database state
echo "📊 Final database state:"
PGPASSWORD=$DB_PASSWORD psql -h localhost -U postgres -d $DB_NAME -c "
SELECT loc, currency, SUM(amount) as total_amount, COUNT(*) as record_count 
FROM holdings 
GROUP BY loc, currency 
ORDER BY loc, currency;
"

echo ""
echo "🎉 All tests completed successfully!"
echo "💡 Database will be cleaned up automatically when script exits"