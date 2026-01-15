package config

import (
	"chat-go/internal/models"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
)

var db *mongo.Client
var dbName string

func InitDB(client *mongo.Client) {
	db = client
	dbName = os.Getenv("MONGO_DB_NAME")
	if dbName == "" {
		dbName = "app"
	}
}

// GetUserCollection returns the users collection from the provided mongo client.
func GetUserCollection() *mongo.Collection {
	if db == nil || dbName == "" {
		return nil
	}
	return db.Database(dbName).Collection(string(models.UserCollection))
}

func GetMessagesCollection() *mongo.Collection {
	if db == nil || dbName == "" {
		return nil
	}
	return db.Database(dbName).Collection(string(models.MessagesCollection))
}

func GetConversationsCollection() *mongo.Collection {
	if db == nil || dbName == "" {
		return nil
	}
	return db.Database(dbName).Collection(string(models.ConversationsCollection))
}
