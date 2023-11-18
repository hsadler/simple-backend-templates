#!/bin/bash

# Set environment variables
export DATABASE_URL="postgres://$DB_USER:$DB_PASSWORD@$DB_HOST:$DB_PORT/$DB_NAME?sslmode=disable"

# Run database migrations
migrate -path ./migrations -database "$DATABASE_URL" up

# Start your server
./app
