package db

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func InitDB() (*mongo.Client, error) {
	// Initialize and return the MongoDB client
	mongoURL := os.Getenv("MONGO_URI")
	if mongoURL == "" {
		log.Println("No MONGO_URL found in .env, using default mongodb://localhost:27017")
		mongoURL = "mongodb://localhost:27017"
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	clientOptions := options.Client().ApplyURI(mongoURL)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal("MongoDB connection error:", err)
		return nil, err
	}
	// Test the connection
	if pinfError := client.Ping(ctx, nil); pinfError != nil {
		log.Fatal("MongoDB ping failed:", pinfError)
		return nil, pinfError
	}
	return client, nil
}
