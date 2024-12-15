package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"zocket/api/routes"
	"zocket/cache"
	"zocket/database"
	"zocket/messaging"

)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: No .env file found, proceeding with system environment variables.")
	}

	// Connect to PostgreSQL
	db, err := database.ConnectPostgres()
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}
	defer db.Close()

	// Connect to Redis
	redisClient, err := cache.ConnectRedis()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer redisClient.Close()

	// Connect to RabbitMQ
	rabbitMQConn, err := messaging.ConnectRabbitMQ()
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer rabbitMQConn.Close()

	// Initialize Router
	r := mux.NewRouter()

	// Register API routes
	// routes.RegisterRoutes(r, db)

	routes.RegisterRoutes(r, db, redisClient, rabbitMQConn)

	// Start HTTP server
	serverAddress := ":8080"
	log.Printf("Starting server on %s", serverAddress)
	if err := http.ListenAndServe(serverAddress, r); err != nil {
		log.Fatalf("Failed to start server on %s: %v", serverAddress, err)
	}
}
