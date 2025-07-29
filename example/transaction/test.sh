#!/bin/bash
docker stop test-postgres && docker rm test-postgres
docker run --name test-postgres -e POSTGRES_PASSWORD=password -e POSTGRES_DB=testdb -p 5432:5432 -d postgres:15 
PGPASSWORD=password psql -h localhost -U postgres -d testdb -c "CREATE TABLE holdings (ts TIMESTAMP WITH TIME ZONE DEFAULT NOW(), loc VARCHAR(50) NOT NULL, currency VARCHAR(10) NOT NULL, amount NUMERIC(20,8) NOT NULL, notes TEXT);" 

# go run transaction.go && docker stop test-postgres && docker rm test-postgres
