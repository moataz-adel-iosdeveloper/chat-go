package conversationsRepository

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

// helper to get the conversations collection at runtime. Returns an error if DB not initialized.
func getConversationsCollection() (*mongo.Collection, error) {
	col := services.GetConversationsCollection()
	if col == nil {
		return nil, errors.New("database not initialized: call services.InitDB before using repositories")
	}
	return col, nil
}

func CreateConversation(conversation *models.Conversation) (*models.Conversation, error) {
	col, err := getConversationsCollection()
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	result, err := col.InsertOne(ctx, conversation)
	if err != nil {
		return nil, err
	}
	var insertedConversation models.Conversation
	err = col.FindOne(ctx, bson.M{"_id": result.InsertedID}).Decode(&insertedConversation)
	if err != nil {
		return nil, err
	}
	return &insertedConversation, nil
}

func FindConversationByID(id string) (*models.Conversation, error) {
	col, err := getConversationsCollection()
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var conversation models.Conversation
	err = col.FindOne(ctx, bson.M{"_id": objID}).Decode(&conversation)
	if err != nil {
		return nil, err
	}
	return &conversation, nil
}

func FindConversationByTwoUserID(senderID string, reseverID string) (*models.Conversation, error) {
	col, err := getConversationsCollection()
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	senderObjID, err := primitive.ObjectIDFromHex(senderID)
	if err != nil {
		return nil, err
	}
	reseverObjID, err := primitive.ObjectIDFromHex(reseverID)
	if err != nil {
		return nil, err
	}
	filter := bson.M{
		"participant_ids": bson.M{
			"$all": []primitive.ObjectID{senderObjID, reseverObjID},
		},
	}
	cursor, err := col.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var conversations []*models.Conversation
	for cursor.Next(ctx) {
		var conversation models.Conversation
		if err := cursor.Decode(&conversation); err != nil {
			return nil, err
		}
		conversations = append(conversations, &conversation)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}
	if len(conversations) == 0 {
		return nil, nil
	}
	return conversations[0], nil
}

func UpdateConversation(conversation *models.Conversation) (*models.Conversation, error) {
	col, err := getConversationsCollection()
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	update := bson.M{
		"$set": bson.M{
			"last_message": conversation.LastMessage,
			"updated_at":   conversation.UpdatedAt,
		},
	}

	_, err = col.UpdateOne(
		ctx,
		bson.M{"_id": conversation.ID},
		update,
	)
	if err != nil {
		return nil, err
	}

	var updatedConversation models.Conversation
	err = col.FindOne(ctx, bson.M{"_id": conversation.ID}).Decode(&updatedConversation)
	if err != nil {
		return nil, err
	}

	return &updatedConversation, nil
}
