#!/bin/bash

# This script creates a new Goose migration file
# Usage: ./scripts/create_migration.sh my_migration_name

set -e

# Check if goose is installed
if ! command -v goose &> /dev/null; then
    echo "Error: goose is not installed"
    echo "Install it with: go install github.com/pressly/goose/v3/cmd/goose@latest"
    exit 1
fi

# Check if a migration name was provided
if [ -z "$1" ]; then
    echo "Error: No migration name provided"
    echo "Usage: ./scripts/create_migration.sh my_migration_name"
    exit 1
fi

# Create the migration
migration_name=$1
migrations_dir="./migrations"

# Make sure migrations directory exists
mkdir -p "$migrations_dir"

# Create the migration file
goose -dir "$migrations_dir" create "$migration_name" sql

echo "Migration created successfully"