package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/fabianpoels/ai-bubble-tracker/db"
	"github.com/fabianpoels/ai-bubble-tracker/server"
	"github.com/joho/godotenv"
)

func main() {
	environment := flag.String("e", "development", "environment")
	os.Setenv("environment", *environment)

	task := flag.String("task", "server", "Task to run")

	flag.Usage = func() {
		fmt.Println("Usage: go run main.go -e {mode} -task {task}")
		fmt.Println("\nAvailable tasks:")
		fmt.Println("  server      - Start the API server")
		fmt.Println("  db-create   - Create database tables")
		fmt.Println("  db-drop     - Drop database tables (careful!)")
	}
	flag.Parse()

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error reading env file. Err: %s", err)
	}

	switch *task {
	case "server":
		server.Init()

	case "db-create":
		if err := db.CreateDatabase(); err != nil {
			log.Fatalf("Failed to create database: %v", err)
		}

	case "db-drop":
		fmt.Print("Are you sure you want to drop all tables? (yes/no): ")
		var confirm string
		fmt.Scanln(&confirm)
		if confirm == "yes" {
			if err := db.DropDatabase(); err != nil {
				log.Fatalf("Failed to drop database: %v", err)
			}
		} else {
			log.Println("Operation cancelled")
		}

	default:
		flag.Usage()
		os.Exit(1)
	}

	// Clean up database connections
	defer db.Close()
}
