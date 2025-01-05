package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"earthquake/internal/database"
	"earthquake/internal/handlers"

	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Get environment variables
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default to 8080 if PORT is not set
	}

	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		log.Fatal("MONGO_URI not set in .env file")
	}

	// Connect to MongoDB
	database.Connect(mongoURI)

	// Routes
	http.HandleFunc("/test", handlers.HandleTestRequest) 
	http.HandleFunc("/result", handlers.GetTestResult)

	serverAddress := fmt.Sprintf("http://localhost:%s", port)
	fmt.Printf("Server is running on %s\n", serverAddress)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
