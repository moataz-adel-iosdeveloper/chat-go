package MessagesRepository

import (
	"chat-go/internal/models"
	"chat-go/internal/services"
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// helper to get the messages collection at runtime. Returns an error if DB not initialized.
func getMessagesCollection() (*mongo.Collection, error) {
	col := services.GetMessagesCollection()
	if col == nil {
		return nil, errors.New("database not initialized: call services.InitDB before using repositories")
	}
	return col, nil
}

func CreateMessage(message *models.Message) (*models.Message, error) {
	col, err := getMessagesCollection()
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	result, err := col.InsertOne(ctx, message)
	if err != nil {
		return nil, err
	}
	var insertedMessage models.Message
	err = col.FindOne(ctx, bson.M{"_id": result.InsertedID}).Decode(&insertedMessage)
	if err != nil {
		return nil, err
	}
	return &insertedMessage, nil
}

func FindMessageByID(id string) (*models.Message, error) {
	col, err := getMessagesCollection()
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var message models.Message
	err = col.FindOne(ctx, bson.M{"_id": objID}).Decode(&message)
	if err != nil {
		return nil, err
	}
	return &message, nil
}

func FindAllMessagesByConversationID(conversationID string) ([]*models.Message, error) {
	col, err := getMessagesCollection()
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	filter := bson.M{"conversation_id": conversationID}
	cursor, err := col.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var messages []*models.Message
	for cursor.Next(ctx) {
		var message models.Message
		if err := cursor.Decode(&message); err != nil {
			return nil, err
		}
		messages = append(messages, &message)
	}
	return messages, nil
}
