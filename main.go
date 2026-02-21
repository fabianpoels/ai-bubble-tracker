package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/fabianpoels/ai-bubble-tracker/server"
	"github.com/joho/godotenv"
)

func main() {
	environment := flag.String("e", "development", "environment")
	os.Setenv("environment", *environment)

	task := flag.String("task", "server", "Task to run (server)")

	flag.Usage = func() {
		fmt.Println("Usage: go run main.go -e {mode} -task {server|add-user|add-location}")
	}
	flag.Parse()

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error reading env file. Err: %s", err)
	}

	switch *task {
	case "server":
		server.Init()
	}
}
