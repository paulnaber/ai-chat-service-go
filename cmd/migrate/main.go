package main

import (
	"flag"
	"fmt"
	"os"

	"ai-chat-service-go/internal/config"
	"ai-chat-service-go/internal/db"
)

// Command-line tool for database migrations using Goose
func main() {
	// Define command-line flags
	upCmd := flag.NewFlagSet("up", flag.ExitOnError)
	downCmd := flag.NewFlagSet("down", flag.ExitOnError)
	statusCmd := flag.NewFlagSet("status", flag.ExitOnError)
	createCmd := flag.NewFlagSet("create", flag.ExitOnError)

	// Create command requires a name parameter
	createName := createCmd.String("name", "", "Name of the migration to create")

	// Check for command
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Error loading configuration: %v\n", err)
		os.Exit(1)
	}

	// Parse command
	switch os.Args[1] {
	case "up":
		upCmd.Parse(os.Args[2:])
		if err := db.MigrateUp(cfg.Database); err != nil {
			fmt.Printf("Error running migrations up: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Migrations applied successfully")

	case "down":
		downCmd.Parse(os.Args[2:])
		if err := db.MigrateDown(cfg.Database); err != nil {
			fmt.Printf("Error rolling back migrations: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Migration rolled back successfully")

	case "status":
		statusCmd.Parse(os.Args[2:])
		if err := db.MigrateStatus(cfg.Database); err != nil {
			fmt.Printf("Error getting migration status: %v\n", err)
			os.Exit(1)
		}

	case "create":
		createCmd.Parse(os.Args[2:])
		if *createName == "" {
			createCmd.PrintDefaults()
			os.Exit(1)
		}
		if err := createMigration(cfg.Database.MigrationDir, *createName); err != nil {
			fmt.Printf("Error creating migration: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Migration created successfully\n")

	default:
		printUsage()
		os.Exit(1)
	}
}

// createMigration creates a new migration file
func createMigration(dir, name string) error {
	// This is a simple implementation, in a real app you would use goose.CreateWithTemplate
	filename := fmt.Sprintf("%s/%s.sql", dir, name)
	
	// Check if file already exists
	if _, err := os.Stat(filename); err == nil {
		return fmt.Errorf("migration file already exists: %s", filename)
	}
	
	// Create the file
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	
	// Write template content
	_, err = file.WriteString(`-- +goose Up
-- SQL in this section is executed when the migration is applied.

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
`)
	
	return err
}

// printUsage prints usage information
func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  migrate up     - Run migrations up")
	fmt.Println("  migrate down   - Roll back the most recent migration")
	fmt.Println("  migrate status - Show migration status")
	fmt.Println("  migrate create -name NAME - Create a new migration")
}