package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Conversation struct {
	ID             primitive.ObjectID   `json:"id" bson:"_id,omitempty"`
	ParticipantIDs []primitive.ObjectID `json:"participant_ids" bson:"participant_ids"`
	LastMessage    *Message             `json:"last_message" bson:"last_message"`
	CreatedAt      primitive.DateTime   `json:"created_at" bson:"created_at"`
	UpdatedAt      primitive.DateTime   `json:"updated_at" bson:"updated_at"`
}

type ConversationResponse struct {
	ID             string    `json:"id"`
	ParticipantIDs []string  `json:"participant_ids"`
	LastMessage    *Message  `json:"last_message,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
