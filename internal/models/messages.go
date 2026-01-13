package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Message struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	ConversationID primitive.ObjectID `bson:"conversation_id" json:"conversation_id"`
	SenderID       primitive.ObjectID `bson:"sender_id" json:"sender_id"`
	IsRead         bool               `bson:"is_read" json:"is_read"`
	Content        string             `bson:"content" json:"content"`
	ContentType    string             `bson:"content_type" json:"content_type"`
	CreatedAt      primitive.DateTime `bson:"createdAt" json:"created_at"`
}

type MessagePayload struct {
	ConversationID string `json:"conversation_id"`
	Content        string `json:"content"`
	ContentType    string `json:"content_type"`
}
