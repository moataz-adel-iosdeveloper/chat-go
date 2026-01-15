package main

import (
	"chat-go/cmd/db"
	"chat-go/cmd/server"
	"chat-go/internal/config"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	// Load .env only in local development
	if os.Getenv("RAILWAY_ENVIRONMENT") == "" {
		err := godotenv.Load()
		if err != nil {
			log.Println("No .env file found")
			log.Fatal(err)
		}
	}
	config.InitJWTService()
	client, err := db.InitDB()
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	// initialize services with mongo client before starting the HTTP server
	config.InitDB(client)

	// pass mongo client to server and start it
	newServer := server.NewWebServer(client)
	if err := newServer.Connect(); err != nil {
		log.Fatal(err)
	}
}
