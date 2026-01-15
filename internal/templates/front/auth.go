package front

import (
	"time"
)

// UserResponse represents a user in the UI
type UserResponse struct {
	ID       string
	Username string
	Email    string
}

// ConversationResponse represents a conversation in the UI
type ConversationResponse struct {
	ID             string
	ParticipantIDs []string
	LastMessage    *MessageResponse
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// MessageResponse represents a message in the UI
type MessageResponse struct {
	ID          string
	SenderID    string
	Content     string
	ContentType string
	CreatedAt   time.Time
}

// ChatMessage represents a chat message
type ChatMessage struct {
	ID                string
	SenderID          string
	Content           string
	ContentType       string
	CreatedAt         time.Time
	IsFromCurrentUser bool
}

// ChatPageData represents data for the chat page
type ChatPageData struct {
	ConversationID string
	Messages       []ChatMessage
	OtherUserName  string
}

// // Re-export auth functions
// func Login(isError bool, errorMessage string) interface {
// 	Render(ctx context.Context, w io.Writer) error
// } {
// 	return authLogin(isError, errorMessage)
// }

// func Register(isError bool, errorMessage string) interface {
// 	Render(ctx context.Context, w io.Writer) error
// } {
// 	return authRegister(isError, errorMessage)
// }
