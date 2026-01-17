package ChatControllers

import (
	"chat-go/internal/config"
	"chat-go/internal/models"
	conversationsRepository "chat-go/internal/repositories/conversations"
	"log"
	"net/http"
)

var chatManager *ChatManager

func init() {
	// Initialize any required resources or configurations here
	chatManager = NewChatManager()
}

func ChatHandler(w http.ResponseWriter, req *http.Request) {

	log.Println("Server started from ChatHandler")
	// get another user id from query params
	conversationID := req.URL.Query().Get("conversation_id")
	if conversationID == "" {
		log.Println("user_id query parameter is required")
		http.Error(w, "conversationID query parameter is required", http.StatusBadRequest)
		return
	}
	// Get the authenticated user ID from context
	authUserID := req.Context().Value("userID")
	if authUserID == nil {
		log.Println("Unauthorized")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Find a conversation between the two users by ID
	conversation, err := conversationsRepository.FindConversationByID(conversationID)
	if err != nil {
		log.Println("Error finding conversation:", err)
		http.Error(w, "Error finding conversation", http.StatusInternalServerError)
		return
	}

	connectionHandler(conversation, authUserID.(string), w, req)
}

func connectionHandler(conversation *models.Conversation, userID string, w http.ResponseWriter, req *http.Request) {
	connection, err := config.SocketServe(w, req)
	if err != nil {
		log.Println("WebSocket connection error:", err)
		return
	}
	client := NewClient(connection, chatManager, conversation.ID.Hex(), userID)
	chatManager.RegisterClient(client, conversation)
	go client.ReadMessages()
	go client.WriteMessage()
}
